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
	"reflect"

	"github.com/8treenet/gotree/helper"

	"github.com/8treenet/gotree/dao/orm"
	"github.com/8treenet/gotree/remote_call"
)

//ComController
type ComController struct {
	remote_call.RpcController
	selfName string
}

//Gotree
func (self *ComController) Gotree(child interface{}) *ComController {
	self.RpcController.Gotree(self)
	self.AddChild(self, child)

	type fun interface {
		RpcName() string
	}
	self.selfName = self.TopChild().(fun).RpcName()
	return self
}

//Model 服务定位器获取model
func (self *ComController) Model(child interface{}) {
	modelDao := reflect.ValueOf(child).Elem().Interface().(comName).Com()
	if self.selfName != modelDao {
		helper.Exit("model 不在一个 com 下,不要乱调用")
	}

	err := _msl.Service(child)
	if err != nil {
		helper.Exit("禁止调用 " + err.Error())
	}
	return
}

//Cache 服务定位器获取Cache
func (self *ComController) Cache(child interface{}) {
	cacheDao := reflect.ValueOf(child).Elem().Interface().(comName).Com()
	if self.selfName != cacheDao {
		helper.Exit("cachae 不在一个 com 下,不要乱调用")
	}

	err := _csl.Service(child)
	if err != nil {
		helper.Exit("禁止调用 " + err.Error())
	}
	return
}

//Api 服务定位器获取Api
func (self *ComController) Api(child interface{}) {
	err := _api.Service(child)
	if err != nil {
		helper.Exit("禁止调用 " + err.Error())
	}
	return
}

//Memory 服务定位器获取Memory
func (self *ComController) Memory(child interface{}) {
	apiDao := reflect.ValueOf(child).Elem().Interface().(comName).Com()
	if self.selfName != apiDao {
		helper.Exit("memory 不在一个 com 下,不要乱调用")
	}

	err := _esl.Service(child)
	if err != nil {
		helper.Exit("禁止调用 " + err.Error())
	}
	return
}

//Transaction 事务
func (self *ComController) Transaction(fun func() error) error {
	return orm.Transaction(self.selfName, fun)
}

//Queue 队列处理
func (self *ComController) Queue(name string, fun func() error) {
	q, ok := queueMap[self.selfName+"_"+name]
	if !ok {
		helper.Exit("未注册队列:" + self.selfName + "." + name)
	}
	q.cast(fun)
}

//TotalPage 总页数
func (self *ComController) TotalPage(size, pageSize int) int {
	if size%pageSize == 0 {
		return size / pageSize
	}
	return size/pageSize + 1
}
