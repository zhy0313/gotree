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
	"errors"
	"fmt"
	"sync"

	"github.com/8treenet/gotree/lib/g"
)

type GoDict struct {
	m     map[string]*Dict
	mutex *sync.Mutex
}

func (self *GoDict) Gotree() *GoDict {
	self.m = make(map[string]*Dict)
	self.mutex = new(sync.Mutex)
	return self
}

func (self *GoDict) Set(key string, value interface{}) {
	defer self.mutex.Unlock()
	id := self.goPoint()
	self.mutex.Lock()

	dict, ok := self.m[id]
	if !ok {
		dict = new(Dict).Gotree()
		self.m[id] = dict
	}
	dict.Set(key, value)
	return
}

func (self *GoDict) Get(key string) interface{} {
	defer self.mutex.Unlock()
	id := self.goPoint()
	self.mutex.Lock()
	v, ok := self.m[id]
	if !ok {
		return nil
	}

	return v.GetInterface(key)
}

func (self *GoDict) Eval(key string, value interface{}) error {
	defer self.mutex.Unlock()
	id := self.goPoint()
	self.mutex.Lock()
	v, ok := self.m[id]
	if !ok {
		return errors.New("undefined key:" + key)
	}
	return v.Get(key, value)
}

func (self *GoDict) Del(key string) {
	defer self.mutex.Unlock()
	id := self.goPoint()
	self.mutex.Lock()

	v, ok := self.m[id]
	if !ok {
		return
	}
	v.Del(key)
}

func (self *GoDict) Remove() {
	defer self.mutex.Unlock()
	id := self.goPoint()
	self.mutex.Lock()
	v, ok := self.m[id]
	if ok {
		v.DelAll()
	}
	delete(self.m, id)
}

// func (self *GoDict) goId() string {
// 	var buf [64]byte
// 	n := runtime.Stack(buf[:], false)
// 	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
// 	return idField
// }

func (self *GoDict) goPoint() string {
	return fmt.Sprint(g.RuntimePointer())
}
