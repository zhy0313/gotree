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

package helper

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
)

type godictInter interface {
	Remove()
	Set(key string, value interface{})
	Get(key string) interface{}
}

var _gGoDict godictInter

func NewGroup() *TaskGroup {
	return new(TaskGroup).taskGroup()
}

func SetGoDict(godict godictInter) {
	_gGoDict = godict
}

//任务池
type TaskPool struct {
	queue     chan *task //队列管道
	close     bool       //是否关闭
	length    int        //池大小
	gseq      string     //gseq
	closeLock *sync.RWMutex
}

//组合task
type TaskGroup struct {
	TaskPool
	list []*task
}

//任务
type task struct {
	callFunction func() error //同步匿名函数方式
	result       interface{}  //返回值
	done         chan error   //错误返回
}

//TaskPopl 构造
func (self *TaskPool) taskPool(length int) *TaskPool {
	self.length = length
	self.queue = make(chan *task, length*16)
	self.close = false
	self.closeLock = new(sync.RWMutex)
	return self
}

//Start 启动 go程池
func (self *TaskPool) start() {
	for index := 0; index < self.length; index++ {
		go self.run()
	}
}

//Run 处理调用
func (self *TaskPool) run() {
	defer func() {
		if _gGoDict != nil {
			_gGoDict.Remove()
		}
	}()
	if self.gseq != "" {
		_gGoDict.Set("gseq", self.gseq)
	}
	for {
		self.runTask()
		if self.isQuit() {
			break
		}
		runtime.Gosched()
	}
}

//runTask 处理task
func (self *TaskPool) runTask() {
	var callTask *task
	defer func() {
		if perr := recover(); perr != nil {
			Log().WriteError(perr)
			if callTask != nil {
				callTask.done <- errors.New(fmt.Sprint(perr))
			}
		}
	}()

	select {
	case callTask = <-self.queue:
		if callTask.callFunction != nil {
			callTask.done <- callTask.callFunction()
		}
		return
	default:
		return
	}
}

//isClose 关闭
func (self *TaskPool) isQuit() bool {
	defer self.closeLock.RUnlock()
	self.closeLock.RLock()
	return self.close
}

//close 关闭
func (self *TaskPool) quit() {
	//优雅关闭
	self.closeLock.Lock()
	self.close = true
	self.closeLock.Unlock()
}

//TaskGroup 构造
func (self *TaskGroup) taskGroup() *TaskGroup {
	self.list = make([]*task, 0, 5)
	self.TaskPool.taskPool(5)
	return self
}

//CallFAddCallFuncunc 同步匿名回调
func (self *TaskGroup) Add(fun func() error) {
	if self.isQuit() {
		return
	}

	ts := new(task)
	ts.callFunction = fun
	ts.done = make(chan error)

	self.list = append(self.list, ts)
	return
}

//Wait 等待所有任务执行完成
func (self *TaskGroup) Wait(numGoroutine ...int) error {
	//如果关闭 不接受调用
	if self.isQuit() {
		return nil
	}
	defer func() {
		self.quit()
	}()

	if _gGoDict != nil {
		//如果有 _taskGoDict 并且有bseq 读取并设置
		gseq := _gGoDict.Get("gseq")
		if gseq != nil {
			str, ok := gseq.(string)
			if ok {
				self.gseq = str
			}
		}
	}
	_gonum := 64
	if len(numGoroutine) > 0 {
		_gonum = numGoroutine[0]
	}

	//限制并发长度最大64
	if len(self.list) < _gonum {
		self.length = len(self.list)
	} else {
		self.length = _gonum
	}
	self.queue = make(chan *task, len(self.list))

	self.start()
	for index := 0; index < len(self.list); index++ {
		self.queue <- self.list[index]
	}

	var result error
	for index := 0; index < len(self.list); index++ {
		err := <-self.list[index].done
		close(self.list[index].done)
		if err != nil {
			result = err
		}
	}
	return result
}
