package business

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"jryghq.cn/lib"
	"jryghq.cn/lib/rpc"
	"jryghq.cn/utils"
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
	bseq      string
	exit      chan bool
}

//Async 异步执行
func (self *async) async(run func(ac AsyncController), completef func()) *async {
	self.Object.Object(self)
	self.run = run
	self.completef = completef
	self.AddSubscribe("shutdown", self.shutdown)
	bseqi := rpc.GoDict().Get("bseq")
	if bseqi != nil {
		str, ok := bseqi.(string)
		if ok {
			self.bseq = str
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
			atomic.AddInt32(&asyncNum, -1)
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
	callCompletef := false
	defer func() {
		if perr := recover(); perr != nil {
			if !callCompletef && self.completef != nil {
				//未执行Complete 并且Complete不为空
				self.completef()
			}
			utils.Log().WriteError(perr)
		}
		atomic.AddInt32(&asyncNum, -1)
		if self.bseq != "" {
			rpc.GoDict().Remove()
		}
		self.DelSubscribe("shutdown")
	}()

	if self.bseq != "" {
		rpc.GoDict().Set("bseq", self.bseq)
	}

	self.run(self)
	if self.completef != nil {
		callCompletef = true
		self.completef()
	}
}

func AsyncNum() int {
	num := atomic.LoadInt32(&asyncNum)
	return int(num)
}
