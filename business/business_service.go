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
package business

import (
	"strconv"
	"strings"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
	"github.com/8treenet/gotree/lib/rpc"
	"github.com/8treenet/gotree/remote_call"
)

type BusinessService struct {
	lib.Object
	_openService bool
	_head        remote_call.RpcHeader
}

func (self *BusinessService) Gotree(child interface{}) *BusinessService {
	self.Object.Gotree(self)
	self.AddChild(self, child)
	self._openService = false
	return self
}

func (self *BusinessService) CallDao(obj interface{}, reply interface{}) error {
	if !self._openService {
		helper.Log().WriteError("禁止重复实例化后调用")
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
	mode := helper.Config().String("sys::Mode")
	if mode == "prod" {
		helper.Log().WriteError("生产环境不可以使用单元测试service")
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
		helper.Log().WriteError("禁止重复实例化")
		panic("禁止重复实例化")
	}
	self._openService = true
	return
}
