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
	"sync"
)

type handlerFunc func(args ...interface{})

//观察者基类
type ObServer struct {
	handleMap map[interface{}]handlerFunc
	mutext    sync.RWMutex
}

//初始化
func (self *ObServer) Gotree() *ObServer {
	self.handleMap = make(map[interface{}]handlerFunc)
	return self
}

//加入订阅者
func (self *ObServer) AddSubscribe(o interface{}, handle handlerFunc) {
	defer self.mutext.Unlock()
	self.mutext.Lock()
	self.handleMap[o] = handle
}

//删除订阅者
func (self *ObServer) DelSubscribe(o interface{}) {
	defer self.mutext.Unlock()
	self.mutext.Lock()
	delete(self.handleMap, o)
}

//发布
func (self *ObServer) NotitySubscribe(args ...interface{}) {
	list := make([]handlerFunc, 0, len(self.handleMap))
	self.mutext.RLock()
	for _, handle := range self.handleMap {
		list = append(list, handle)
	}
	self.mutext.RUnlock()

	for index := 0; index < len(list); index++ {
		list[index](args...)
	}
}

//订阅者长度
func (self *ObServer) SubscribeLen() int {
	defer self.mutext.RUnlock()
	self.mutext.RLock()
	return len(self.handleMap)
}
