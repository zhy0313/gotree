package business

import (
	"strings"
	"time"

	"jryghq.cn/lib"
	"jryghq.cn/lib/rpc"
	"jryghq.cn/remote_call"
	"jryghq.cn/utils"
)

var _scl *lib.ServiceLocator //业务服务定位器
var _ssl *lib.ServiceLocator //系统服务定位器
var _tsl *lib.ServiceLocator //定时业务
var _timer []string

func init() {
	utils.LoadConfig("business")
	logOn()
	appStart()
	_scl = new(lib.ServiceLocator).ServiceLocator()
	_ssl = new(lib.ServiceLocator).ServiceLocator()
	_tsl = new(lib.ServiceLocator).ServiceLocator()

	concurrency := concurrency()
	//注册相关组件
	_ssl.AddComponent(new(remote_call.InnerMaster).InnerMaster())
	innerRpcServer := new(remote_call.InnerServerController).InnerServerController()
	client := new(remote_call.RpcClient).RpcClient(concurrency, retry())
	client.RecoverLog()
	_ssl.AddComponent(client)
	_ssl.AddComponent(new(remote_call.RpcQps).RpcQps())
	remote_call.RpcServerRegister(innerRpcServer)
	_timer = []string{}
	lib.SetGoDict(rpc.GoDict())
	remote_call.AsynNumFunc = AsyncNum
}

//RpcServerRegister 注册rpc服务
func RegisterController(controller interface{}) {
	remote_call.RpcServerRegister(controller)
}

//RegisterService注册service
func RegisterService(service interface{}) {
	type openService interface {
		OpenService()
	}

	os := service.(openService)
	os.OpenService()
	_scl.AddService(service)
}

//RegisterTimer注册定时器
func RegisterTimer(service interface{}) {
	type openService interface {
		OpenTimer()
	}

	os := service.(openService)
	os.OpenTimer()
	_tsl.AddService(service)
}

func Run(args ...interface{}) {
	rpc.GoDict().Remove()
	printSystem()
	timerOn()
	var bindAddr string
	if len(args) == 0 {
		bindAddr = utils.Config().String("BindAddr")
	}
	remote_call.StartServerInfoCheck()
	//通知所有定时服务
	_tsl.NotitySubscribe("TimerOn", timerServices()...)
	utils.Log().WriteInfo("jryg business run ....")
	var client *remote_call.RpcClient
	_ssl.GetComponent(&client)
	_scl.NotitySubscribe("startup")
	rpcSer := remote_call.RpcServerRun(bindAddr)
	_scl.NotitySubscribe("shutdown")
	utils.Log().WriteInfo("jryg business close...")
	//优雅关闭
	for index := 0; index < 30; index++ {
		//当前rpc数量, 天定时器, 异步数量
		num := rpc.CurrentCallNum()
		timenum := lib.CurrentTimeNum()
		anum := AsyncNum()
		utils.Log().WriteInfo("jryg business close: 定时服务剩余:", timenum, "请求服务剩余:", num, "异步服务剩余:", anum)
		if num <= 0 && anum <= 0 && timenum <= 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
	rpcSer.Close()
	client.Close()
	utils.Log().WriteInfo("jryg business close: success")
	utils.Log().Close()
}

func timerServices() (list []interface{}) {
	for _, serName := range _timer {
		list = append(list, serName)
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
	//生产环境定位 ~/log/xx目录
	utils.Log().Init(dir, rpc.GoDict())
}

//读取配置文件启动定时 service服务
func timerOn() {
	serList := strings.Split(utils.Config().String("TimerOn"), ",")
	for _, ser := range serList {
		_timer = append(_timer, ser)
	}
}

func concurrency() (concurrency int) {
	//读取最大调用dao并发数量
	concurrency = utils.Config().DefaultInt("sys::CallDaoConcurrency", 128)
	return
}

func retry() (retry int) {
	retry = utils.Config().DefaultInt("sys::CallDaoRetry", 5)
	return
}
