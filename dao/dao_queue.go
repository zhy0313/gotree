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
	dao       string
	name      string
}

func (self *daoQueue) Gotree(queueLen, max int, dao, name string) *daoQueue {
	self.Object.Gotree(self)
	self.queue = make(chan queueCast, queueLen)
	self.maxGo = int32(max)
	self.currentGo = 0
	self.dao = dao
	self.name = name
	return self
}

//Cast 异步调用
func (self *daoQueue) cast(fun func() error) {
	self.openAssist()
	gdict := rpc.GoDict()
	var bseqstr string
	if gdict != nil {
		seq := gdict.Get("bseq")
		if seq != nil {
			bseqstr = seq.(string)
		}
	}
	q := queueCast{f: fun, bseq: bseqstr}
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
			helper.Log().WriteWarn(self.dao+"."+self.name+" queue: ", fmt.Sprint(perr))
		}
	}()

	rpc.GoDict().Set("bseq", f.bseq)
	err := f.f()
	if err != nil {
		helper.Log().WriteWarn(self.dao+"."+self.name+" queue: ", err)
	}
}

type queueCast struct {
	f    func() error
	bseq string
}
