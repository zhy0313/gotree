// Copyright gotree Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dao

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
	"github.com/8treenet/gotree/lib/rpc"
	"github.com/8treenet/gotree/remote_call"
)

type comName interface {
	Com() string
}

var _msl *lib.ServiceLocator //模型服务定位器
var _csl *lib.ServiceLocator //缓存服务定位器
var _esl *lib.ServiceLocator //内存服务定位器
var _api *lib.ServiceLocator //api服务定位器
var _daoOnList []comNode
var queueMap map[string]*daoQueue

type comNode struct {
	Name  string
	Id    int
	Extra []interface{}
}

func init() {
	helper.LoadConfig("dao")
	if helper.Config().String("sys::Mode") == "dev" {
		modelProfiler = true
	}
	logOn()
	appStart()
	queueMap = make(map[string]*daoQueue)

	_msl = new(lib.ServiceLocator).Gotree() //model数据源
	_csl = new(lib.ServiceLocator).Gotree() //cache数据源
	_api = new(lib.ServiceLocator).Gotree() //api数据源
	_esl = new(lib.ServiceLocator).Gotree() //内存数据源

	helper.SetGoDict(rpc.GoDict())
	lib.SetGoDict(rpc.GoDict())
	tp = new(lib.LimiteGo).Gotree(helper.Config().DefaultInt("sys::ApiConcurrency", 1024))

	innerRpcServer := new(remote_call.InnerServerController).Gotree()
	remote_call.RpcServerRegister(innerRpcServer)
}

//RpcServerRegister 注册rpc服务
func RegisterController(controller interface{}) {
	if helper.Testing() {
		return
	}
	remote_call.RpcServerRegister(controller)
}

//RegisterModel 注册model
func RegisterModel(service interface{}) {
	if helper.Testing() {
		return
	}
	type init interface {
		DaoInit()
	}
	service.(init).DaoInit()
	if _msl.CheckService(service) {
		helper.Exit("RegisterModel duplicate registration")
	}
	_msl.AddService(service)
}

//RegisterModel 注册cache
func RegisterCache(service interface{}) {
	if helper.Testing() {
		return
	}
	type init interface {
		DaoInit()
	}
	service.(init).DaoInit()
	if _csl.CheckService(service) {
		helper.Exit("RegisterCache duplicate registration")
	}
	_csl.AddService(service)
}

//RegisterMemory 注册内存
func RegisterMemory(service interface{}) {
	if helper.Testing() {
		return
	}
	type init interface {
		DaoInit()
	}
	service.(init).DaoInit()
	if _esl.CheckService(service) {
		helper.Exit("RegisterMemory duplicate registration")
	}
	_esl.AddService(service)
}

//RegisterApi 注册api
func RegisterApi(service interface{}) {
	if helper.Testing() {
		return
	}
	if _api.CheckService(service) {
		helper.Exit("RegisterApi duplicate registration")
	}
	_api.AddService(service)
}

//RegisterQueue 队列
//queueName:队列名字
//queueLen:队列长度
//goroutine:队列消费的go程数量,默认是1
func RegisterQueue(controller interface{}, queueName string, queueLen int, goroutine ...int) {
	type rpcname interface {
		RpcName() string
	}

	rc, ok := controller.(rpcname)
	if !ok {
		helper.Exit("RegisterQueue error")
	}
	dao := rc.RpcName()
	mgo := 1
	if len(goroutine) > 0 && goroutine[0] > 0 {
		mgo = goroutine[0]
	}
	q := new(daoQueue).Gotree(queueLen, mgo, dao, queueName)
	go q.mainRun()
	queueMap[dao+"_"+queueName] = q
}

//daoOn 开启dao
func daoOn() {
	openDao, err := helper.Config().GetSection("com_on")
	if err != nil {
		helper.Exit("daoOn-openDao Not found com.conf com_on:" + err.Error())
	}
	controllers := remote_call.RpcControllerNames()
	for k, v := range openDao {
		id, e := strconv.Atoi(v)
		if e != nil {
			helper.Log().Error("daoOn-openDao dao id error:", k, v)
			continue
		}
		var comName string
		for index := 0; index < len(controllers); index++ {
			if strings.ToLower(controllers[index]) == k {
				comName = controllers[index]
			}
		}
		if comName == "" {
			helper.Log().Error("daoOn-openDao error:Not found com:", k)
			continue
		}

		extra := []interface{}{}
		extraList := strings.Split(helper.Config().String("com_extra::"+comName), ",")
		for _, item := range extraList {
			extra = append(extra, item)
		}
		_daoOnList = append(_daoOnList, comNode{Name: comName, Id: id, Extra: extra})
	}
}

func Run(args ...interface{}) {
	var bindAddr string
	if len(args) == 0 {
		bindAddr = helper.Config().String("dispersed::BindAddr")
	}

	tick := lib.RunTick(1000, memoryTimeout, "memoryTimeout", 3000)
	daoOn()
	telnet()

	ic := new(remote_call.InnerClient).Gotree()
	initInnerClient(ic)

	//通知所有model dao关联
	task := helper.NewGroup()
	for _, daoItem := range daos() {
		daonode := daoItem.(comNode)
		task.Add(func() error {
			_msl.NotitySubscribe("ModelOn", daonode.Name)
			_msl.NotitySubscribe("CacheOn", daonode.Name)
			_msl.NotitySubscribe("MemoryOn", daonode.Name)
			_msl.NotitySubscribe("ApiOn", daonode.Name)
			return nil
		})
	}
	if e := task.Wait(); e != nil {
		helper.Exit(e.Error())
	}
	_esl.NotitySubscribe("startup")

	rpcser := remote_call.RpcServerRun(bindAddr, func(svrName string) {
		for _, item := range _daoOnList {
			if svrName == item.Name {
				helper.Log().Notice("startup:", svrName, "id:", item.Id)
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
		qlen := allQueueLen()
		helper.Log().Notice("Run dao close: Request service surplus:", num, "Queue surplus:", qlen)
		if num <= 0 && qlen <= 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
	_esl.NotitySubscribe("shutdown")
	rpcser.Close()
	tick.Stop()
	helper.Log().Notice("dao close")
	helper.Log().Close()
}

func initInnerClient(ic *remote_call.InnerClient) {
	baddrs := helper.Config().String("dispersed::BusinessAddrs")
	if baddrs == "" {
		helper.Log().Error("initInnerClient-BusinessAddrs baddrs address is empty.")
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
	remote_call.SetDbCountFunc(dbConnectNum)
	remote_call.SetQueueCountFunc(allQueueLen)
}

func daos() (list []interface{}) {
	for _, dao := range _daoOnList {
		list = append(list, dao)
	}
	return
}

func logOn() {
	dir := helper.Config().DefaultString("sys::LogDir", "log")
	helper.Log().Init(dir, rpc.GoDict())
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

func connectDao(comName string) bool {
	defer ormConMutex.Unlock()
	ormConMutex.Lock()
	if helper.InSlice(ormConnect, comName) {
		return false
	}
	ormConnect = append(ormConnect, comName)
	return true
}

// allQueueLen 当前全部队列长度
func allQueueLen() (result int) {
	for _, q := range queueMap {
		result += len(q.queue)
	}
	return
}
