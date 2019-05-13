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
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib/rpc"
	"github.com/8treenet/gotree/remote_call"
)

var bseq int64
var identification string

func init() {
	bseq = 1
	rand.Seed(time.Now().Unix())
	x := int64(rand.Intn(10000))
	identification = strconv.FormatInt(x, 36)
}

//BusinessController
type BusinessController struct {
	remote_call.RpcController
}

//BusinessController
func (self *BusinessController) Gotree(child interface{}) *BusinessController {
	self.RpcController.Gotree(self)
	self.AddChild(self, child)
	rpc.GoDict().Set("bseq", getBseq())
	return self
}

//服务定位器获取服务
func (self *BusinessController) Service(child interface{}) {
	err := _scl.Service(child)
	if err != nil {
		helper.Log().WriteError("飞哥:不要乱调用:" + err.Error())
		panic("飞哥:不要乱调用:" + err.Error())
	}
	return
}

// ServiceBroadcast 调用所有注册service的method方法. 潜龙勿用,会使项目非常难以维护
func (self *BusinessController) ServiceBroadcast(method string, arg interface{}) {
	if e := _scl.Broadcast(method, arg); e != nil {
		helper.Log().WriteError("Buesiness ServiceBroadcast errror:" + e.Error())
	}
}

//Async 异步执行
func (self *BusinessController) Async(run func(ac AsyncController), completeds ...func()) {
	var completed func()
	if len(completeds) > 0 {
		completed = completeds[0]
	}
	ac := new(async).Gotree(run, completed)
	go ac.execute()
}

var bseqMutex sync.Mutex

func getBseq() (result string) {
	defer bseqMutex.Unlock()
	bseqMutex.Lock()
	result = identification
	result += strconv.FormatInt(bseq, 36)
	if bseq == math.MaxInt64 {
		bseq = 1
		return
	}
	bseq += 1
	return
}
