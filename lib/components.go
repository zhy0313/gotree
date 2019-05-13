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

package lib

import (
	"reflect"

	"github.com/8treenet/gotree/helper"
)

//组件基类
type Components struct {
	dict *Dict
}

//初始化
func (self *Components) Gotree() *Components {
	self.dict = new(Dict).Gotree()
	return self
}

type enterComponent interface {
	EnterComponent()
}

type updateComponent interface {
	UpdateComponent(c *Components)
}

//加入组件
func (self *Components) AddComponent(obj interface{}) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		helper.Log().WriteError("AddComponent != reflect.Ptr")
	}
	self.dict.Set(t.Elem().Name(), obj)
	if app, ok := obj.(updateComponent); ok {
		app.UpdateComponent(self)
	}
	if app, ok := obj.(enterComponent); ok {
		app.EnterComponent()
	}
}

//移除组件
func (self *Components) RemoveComponent(obj interface{}) {
	t := reflect.TypeOf(obj)
	self.dict.Del(t.Name())
}

//获取组件
func (self *Components) GetComponent(obj interface{}) error {
	t := reflect.TypeOf(obj)
	return self.dict.Get(t.Elem().Elem().Name(), obj)
}

//广播组件内所有实现method的方法
func (self *Components) Broadcast(method string, arg interface{}) {
	list := self.dict.Keys()
	for _, v := range list {
		com := self.dict.GetInterface(v)
		if com == nil {
			continue
		}
		value := reflect.ValueOf(com).MethodByName(method)
		if value.Kind() != reflect.Invalid {
			value.Call([]reflect.Value{reflect.ValueOf(arg)})
		}
	}
}
