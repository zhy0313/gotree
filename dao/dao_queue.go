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
	"fmt"
	"sync/atomic"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
	"github.com/8treenet/gotree/lib/rpc"
)

type daoQueue struct {
	lib.Object
	maxGo     int32          //最大go程
	queue     chan queueCast //队列管道
	currentGo int32          //当前执行go程
	com       string
	name      string
}

func (self *daoQueue) Gotree(queueLen, max int, com, name string) *daoQueue {
	self.Object.Gotree(self)
	self.queue = make(chan queueCast, queueLen)
	self.maxGo = int32(max)
	self.currentGo = 0
	self.com = com
	self.name = name
	return self
}

//Cast 异步调用
func (self *daoQueue) cast(fun func() error) {
	self.openAssist()
	gdict := rpc.GoDict()
	var bseqstr string
	if gdict != nil {
		seq := gdict.Get("gseq")
		if seq != nil {
			bseqstr = seq.(string)
		}
	}
	q := queueCast{f: fun, gseq: bseqstr}
	self.queue <- q
	return
}

//openAssist 开启辅助go程
func (self *daoQueue) openAssist() {
	current := atomic.LoadInt32(&self.currentGo)
	if current >= self.maxGo {
		//已达到最大go程
		return
	}

	//如果当前管道有数据 开启辅助go处理
	if len(self.queue) > 0 {
		go self.assistRun()
	}
}

//mainRun 主go程
func (self *daoQueue) mainRun() {
	atomic.AddInt32(&self.currentGo, 1)
	for {
		fun := <-self.queue
		self.execute(fun)
	}
}

//assistRun 辅go程
func (self *daoQueue) assistRun() {
	atomic.AddInt32(&self.currentGo, 1)
	iterator := new(lib.Iterator).Gotree()
	var sleepLen int64 = 0 //总共休眠时长
	for {
		select {
		case fun := <-self.queue:
			self.execute(fun)
			iterator.ResetTimer()
			sleepLen = 0
		default:
			sleepLen += iterator.Sleep()
		}

		//30万毫秒 辅助线程最多存活5分钟
		if sleepLen > 300000 {
			break
		}
	}
	atomic.AddInt32(&self.currentGo, -1)
}

func (self *daoQueue) execute(f queueCast) {
	defer func() {
		if perr := recover(); perr != nil {
			helper.Log().WriteWarn(self.com+"."+self.name+" queue: ", fmt.Sprint(perr))
		}
	}()

	rpc.GoDict().Set("gseq", f.gseq)
	err := f.f()
	if err != nil {
		helper.Log().WriteWarn(self.com+"."+self.name+" queue: ", err)
	}
}

type queueCast struct {
	f    func() error
	gseq string
}
