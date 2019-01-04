package business

import (
	"strconv"
	"strings"

	"jryghq.cn/lib"
	"jryghq.cn/lib/rpc"
	"jryghq.cn/remote_call"
	"jryghq.cn/utils"
)

type BusinessService struct {
	lib.Object
	_openService bool
	_head        remote_call.RpcHeader
}

func (self *BusinessService) BusinessService(child interface{}) *BusinessService {
	self.Object.Object(self)
	self.AddChild(self, child)
	self._openService = false
	return self
}

func (self *BusinessService) CallDao(obj interface{}, reply interface{}) error {
	if !self._openService {
		utils.Log().WriteError("禁止重复实例化后调用")
		panic("禁止重复实例化后调用")
	}
	var client *remote_call.RpcClient
	if e := _ssl.GetComponent(&client); e != nil {
		return e
	}
	return client.Call(obj, reply)
}

//Header 读取go栈中的kv数据
func (self *BusinessService) ReqHeader(k string) string {
	value := rpc.GoDict().Get("head")
	if value == nil {
		return ""
	}
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return self._head.Get(str, k)
}

func (self *BusinessService) TestOn(testDaos ...string) {
	mode := utils.Config().String("sys::mode")
	if mode == "prod" {
		utils.Log().WriteError("生产环境不可以使用单元测试service")
		panic("生产环境不可以使用单元测试service")
	}
	rpc.GoDict().Set("bseq", "ServiceUnit")
	self._openService = true

	var im *remote_call.InnerMaster
	_ssl.GetComponent(&im)

	for _, dao := range testDaos {
		daoNameId := strings.Split(dao, ":")
		id, _ := strconv.Atoi(daoNameId[1])
		im.LocalAddNode(daoNameId[0], "127.0.0.1", "6666", id)
	}
	return
}

func (self *BusinessService) OpenService() {
	exist := _scl.CheckService(self.TopChild())
	if exist {
		utils.Log().WriteError("禁止重复实例化")
		panic("禁止重复实例化")
	}
	self._openService = true
	return
}

//Async 异步执行
func (self *BusinessService) Async(run func(ac AsyncController), completeds ...func()) {
	var completed func()
	if len(completeds) > 0 {
		completed = completeds[0]
	}
	ac := new(async).async(run, completed)
	go ac.execute()
}

// Broadcast 调用所有注册service的method方法. 潜龙勿用,会使项目非常难以维护
func (self *BusinessService) Broadcast(method string, arg interface{}) {
	if e := _scl.Broadcast(method, arg); e != nil {
		utils.Log().WriteError("Buesiness ServiceBroadcast errror:" + e.Error())
	}
}
