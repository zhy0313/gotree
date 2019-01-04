package lib

import (
	"runtime"
	"time"
)

type Iterator struct {
	count               int //当前梯度
	maxCount            int //最大梯度
	tick                int //梯度内执行的次数
	increaseMillisecond int //毫秒基数
	maxMillisecond      int //时间曲线长度,
}

//Iterator 初始化曲线sleep数据
func (self *Iterator) Iterator() *Iterator {
	self.count = 0
	self.maxCount = 32
	self.tick = 0
	self.increaseMillisecond = 8
	self.maxMillisecond = 95000 //95秒后 每次512毫秒沉睡
	return self
}

//Sleep 睡眠
func (self *Iterator) Sleep() {
	defer self.next()
	if self.count == 0 {
		runtime.Gosched()
		return
	}
	time.Sleep(time.Duration(self.increaseMillisecond*self.count) * time.Millisecond)
}

//ResetTimer 重置
func (self *Iterator) ResetTimer() {
	self.count = 0
	self.tick = 0
}

func (self *Iterator) next() {
	if self.count >= self.maxCount {
		return
	}

	//0按1计算
	count := self.count
	if count == 0 {
		count = 1
	}

	//计算出本梯度内的最大次数
	maxTick := int(self.maxMillisecond / self.maxCount / (self.increaseMillisecond * count))
	if self.tick >= maxTick {
		self.count = self.count + 1
		self.tick = 0
		return
	}
	self.tick = self.tick + 1
}
