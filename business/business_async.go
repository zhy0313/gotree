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
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
	"github.com/8treenet/gotree/lib/rpc"
)

type AsyncController interface {
	Sleep(millisecond int64) //休眠
	CancelCompleted()        //取消完成回调
}

type async struct {
	lib.Object
	run       func(ac AsyncController)
	completef func()
	mutex     sync.Mutex
	gseq      string
	exit      chan bool
}

//Async 异步执行
func (self *async) Gotree(run func(ac AsyncController), completef func()) *async {
	self.Object.Gotree(self)
	self.run = run
	self.completef = completef
	self.AddSubscribe("shutdown", self.shutdown)
	bseqi := rpc.GoDict().Get("gseq")
	if bseqi != nil {
		str, ok := bseqi.(string)
		if ok {
			self.gseq = str
		}
	}
	self.exit = make(chan bool, 1)
	return self
}

// Sleep 休眠 毫秒
func (self *async) Sleep(millisecond int64) {
	for {
		var sleepMillisecond int64
		if millisecond < 1 {
			break
		}
		if millisecond < 500 {
			sleepMillisecond = millisecond
			millisecond = 0
		} else {
			sleepMillisecond = 500
			millisecond -= 500
		}
		time.Sleep(time.Duration(sleepMillisecond) * time.Millisecond)

		select {
		case _ = <-self.exit:
			if self.completef != nil {
				self.completef()
			}
			runtime.Goexit()
		default:
			continue
		}
	}
}

//CancelCompleted 取消完成回调
func (self *async) CancelCompleted() {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	self.completef = nil
}

//shutdown 停服终端处理
func (self *async) shutdown(args ...interface{}) {
	self.exit <- true
}

var asyncNum int32 = 0

//Async 异步执行
func (self *async) execute() {
	atomic.AddInt32(&asyncNum, 1)
	defer func() {
		if perr := recover(); perr != nil {
			helper.Log().WriteError(perr)
		}
		atomic.AddInt32(&asyncNum, -1)
		if self.gseq != "" {
			rpc.GoDict().Remove()
		}
		self.DelSubscribe("shutdown")
	}()

	if self.gseq != "" {
		rpc.GoDict().Set("gseq", self.gseq)
	}

	self.run(self)
	if self.completef != nil {
		self.completef()
		self.completef = nil
	}
}

func AsyncNum() int {
	num := atomic.LoadInt32(&asyncNum)
	return int(num)
}
