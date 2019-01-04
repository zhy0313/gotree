package dao

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"jryghq.cn/lib"
	"jryghq.cn/lib/rpc"
	"jryghq.cn/remote_call"
	"jryghq.cn/utils"
)

type daoName interface {
	Dao() string
}

var _msl *lib.ServiceLocator //模型服务定位器
var _csl *lib.ServiceLocator //缓存服务定位器
var _esl *lib.ServiceLocator //内存服务定位器
var _api *lib.ServiceLocator //api服务定位器
var _daoOnList []daoNode

type daoNode struct {
	Name  string
	Id    int
	Extra []interface{}
}

func init() {
	utils.LoadConfig("dao")
	if utils.Config().String("sys::mode") == "dev" || utils.Config().String("sys::mode") == "pre" {
		modelProfiler = true
	}
	logOn()
	appStart()

	_msl = new(lib.ServiceLocator).ServiceLocator() //model数据源
	_csl = new(lib.ServiceLocator).ServiceLocator() //cache数据源
	_api = new(lib.ServiceLocator).ServiceLocator() //api数据源
	_esl = new(lib.ServiceLocator).ServiceLocator() //内存数据源
	lib.SetGoDict(rpc.GoDict())
	tp = lib.NewTaskPool(utils.Config().DefaultInt("sys::ApiConcurrency", 768)) //同时最多768 网络请求并发
	tp.Start()

	innerRpcServer := new(remote_call.InnerServerController).InnerServerController()
	remote_call.RpcServerRegister(innerRpcServer)
}

//RpcServerRegister 注册rpc服务
func RegisterController(controller interface{}) {
	if utils.Testing() {
		return
	}
	remote_call.RpcServerRegister(controller)
}

//RegisterModel 注册model
func RegisterModel(service interface{}) {
	if utils.Testing() {
		return
	}
	type init interface {
		DaoInit()
	}
	service.(init).DaoInit()
	if _msl.CheckService(service) {
		utils.Log().WriteError("RegisterModel 重复注册")
		panic("RegisterModel 重复注册")
	}
	_msl.AddService(service)
}

//RegisterModel 注册cache
func RegisterCache(service interface{}) {
	if utils.Testing() {
		return
	}
	type init interface {
		DaoInit()
	}
	service.(init).DaoInit()
	if _csl.CheckService(service) {
		utils.Log().WriteError("RegisterCache 重复注册")
		panic("RegisterCache 重复注册")
	}
	_csl.AddService(service)
}

//RegisterMemory 注册内存
func RegisterMemory(service interface{}) {
	if utils.Testing() {
		return
	}
	type init interface {
		DaoInit()
	}
	service.(init).DaoInit()
	if _esl.CheckService(service) {
		utils.Log().WriteError("RegisterMemory 重复注册")
		panic("RegisterMemory 重复注册")
	}
	_esl.AddService(service)
}

//RegisterApi 注册api
func RegisterApi(service interface{}) {
	if utils.Testing() {
		return
	}
	type init interface {
		DaoInit()
	}
	service.(init).DaoInit()
	if _api.CheckService(service) {
		utils.Log().WriteError("RegisterApi 重复注册")
		panic("RegisterApi 重复注册")
	}
	_api.AddService(service)
}

//daoOn 开启dao
func daoOn() {
	openDao, err := utils.Config().GetSection("dao_on")
	if err != nil {
		utils.Log().WriteError("未找到 dao dao_on:", err)
		os.Exit(-1)
	}
	controllers := remote_call.RpcControllerNames()
	for k, v := range openDao {
		id, e := strconv.Atoi(v)
		if e != nil {
			utils.Log().WriteError("dao id 错误:", k, v)
			continue
		}
		var daoName string
		for index := 0; index < len(controllers); index++ {
			if strings.ToLower(controllers[index]) == k {
				daoName = controllers[index]
			}
		}
		if daoName == "" {
			utils.Log().WriteError("未找到 dao:", k)
			continue
		}

		extra := []interface{}{}
		extraList := strings.Split(utils.Config().String("dao_extra::"+daoName), ",")
		for _, item := range extraList {
			extra = append(extra, item)
		}
		_daoOnList = append(_daoOnList, daoNode{Name: daoName, Id: id, Extra: extra})
	}
}

func Run(args ...interface{}) {
	var bindAddr string
	if len(args) == 0 {
		bindAddr = utils.Config().String("BindAddr")
	}

	tick := lib.RunTick(1000, memoryTimeout, "memoryTimeout", 3000)
	daoOn()
	telnet()

	ic := new(remote_call.InnerClient).InnerClient()
	initInnerClient(ic)

	//通知所有model dao关联
	task := lib.NewTaskGroup()
	for _, daoItem := range daos() {
		daonode := daoItem.(daoNode)
		task.AddFuncByGroup(func() error {
			_msl.NotitySubscribe("ModelOn", daonode.Name)
			_msl.NotitySubscribe("CacheOn", daonode.Name)
			_msl.NotitySubscribe("MemoryOn", daonode.Name)
			_msl.NotitySubscribe("ApiOn", daonode.Name)
			return nil
		})
	}
	task.WaitByGroup()
	_esl.NotitySubscribe("startup")

	rpcser := remote_call.RpcServerRun(bindAddr, func(svrName string) {
		for _, item := range _daoOnList {
			if svrName == item.Name {
				utils.Log().WriteInfo("启动:", svrName, "id:", item.Id)
				_daoOnList = append(_daoOnList, item)
				return
			}
		}
		//如果未开启dao 取消注册rpc
		remote_call.RpcUnServerRegister(svrName)
	})
	ic.Close()
	//优雅关闭
	for index := 0; index < 30; index++ {
		num := rpc.CurrentCallNum()
		utils.Log().WriteInfo("jryg dao close: 请求服务剩余:", num)
		if num <= 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
	_esl.NotitySubscribe("shutdown")
	rpcser.Close()
	tick.Stop()
	utils.Log().WriteInfo("jryg dao close")
	utils.Log().Close()
}

func initInnerClient(ic *remote_call.InnerClient) {
	baddrs := utils.Config().String("BusinessAddrs")
	if baddrs == "" {
		utils.Log().WriteError("BusinessAddrs地址为空")
	}
	for index := 0; index < len(_daoOnList); index++ {
		ic.AddDaoByNode(_daoOnList[index].Name, _daoOnList[index].Id, _daoOnList[index].Extra...)
	}

	list := strings.Split(baddrs, ",")
	for _, item := range list {
		addr := strings.Split(item, ":")
		port, _ := strconv.Atoi(addr[1])
		ic.AddBusiness(addr[0], port)
	}
	ic.SetDbCountFunc(dbConnectNum)
}

func daos() (list []interface{}) {
	for _, dao := range _daoOnList {
		list = append(list, dao)
	}
	return
}

func logOn() {
	mode := utils.Config().String("sys::mode")
	if mode != "prod" {
		//如果是测试当前目录创建日志文件 并开启fmt.print
		utils.Log().Debug()
	}
	dir := utils.Config().DefaultString("sys::LogDir", "log")
	utils.Log().Init(dir, rpc.GoDict())
}

// memoryTimeout 内存超时检测
func memoryTimeout() {
	_esl.Broadcast("MemoryTimeout", time.Now().Unix())
}

func dbConnectNum() int {
	mconnects := make(map[string]int)
	cconnects := make(map[string]int)
	var result int
	_msl.Broadcast("Connections", mconnects)
	_csl.Broadcast("Connections", cconnects)
	for _, item := range mconnects {
		result += item
	}

	for _, item := range cconnects {
		result += item
	}
	return result
}

var ormConnect []string = []string{}
var ormConMutex sync.Mutex

func connectDao(daoName string) bool {
	defer ormConMutex.Unlock()
	ormConMutex.Lock()
	if utils.InArray(ormConnect, daoName) {
		return false
	}
	ormConnect = append(ormConnect, daoName)
	return true
}
