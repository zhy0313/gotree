package remote_call

import (
	"sync"
	"time"

	"github.com/8treenet/gotree/lib"
)

func init() {
	_breakerManage.dict = new(lib.Dict).Gotree()
}

type breakerCmd interface {
	Control() string
	Action() string
}

//RegisterBreaker 注册熔断
func RegisterBreaker(cmd breakerCmd, rangeSec int, timeoutRatio float32, resumeSec int) {
	serm := cmd.Control() + "." + cmd.Action()
	_breakerManage.AddBreaker(serm, rangeSec, timeoutRatio, resumeSec)
}

type RpcBreak struct {
	lib.Object
}

func (self *RpcBreak) Gotree() *RpcBreak {
	return self
}

func (self *RpcBreak) Breaking(cmd cmdCall) bool {
	serm := cmd.ServiceMethod()
	return _breakerManage.Breaking(serm)
}

func (self *RpcBreak) Call(cmd cmdCall, timeout bool) {
	serm := cmd.ServiceMethod()
	_breakerManage.Call(serm, timeout)
	return
}

func (self *RpcBreak) RunTick() {
	_breakerManage.Run()
}

func (self *RpcBreak) StopTick() {
	_breakerManage.Stop()
}

type breaker struct {
	Conf struct {
		RangeSec     int     //时间范围
		TimeoutRatio float32 //超时比例
		ResumeSec    int     //恢复时间
	}

	CallSum     int  //总调用次数
	CallTimeout int  //超时次数
	Tick        int  //为0重置
	Breaking    bool //是否熔断
	ResumeTick  int  //恢复计数
}

type breakerManage struct {
	dict      *lib.Dict
	dictMutex sync.Mutex
	keys      []interface{}
	stop      chan bool
}

func (self *breakerManage) AddBreaker(serviceMethod string, rangeSec int, timeoutRatio float32, resumeSec int) {
	if self.dict == nil {
		self.dict = new(lib.Dict).Gotree()
	}
	if self.dict.Check(serviceMethod) {
		return
	}
	b := &breaker{}
	b.Conf.RangeSec = rangeSec
	b.Conf.TimeoutRatio = timeoutRatio
	b.Conf.ResumeSec = resumeSec
	b.Tick = b.Conf.RangeSec
	self.dict.Set(serviceMethod, b)
}

func (self *breakerManage) Breaking(serviceMethod string) bool {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	var b *breaker
	if err := self.dict.Get(serviceMethod, &b); err != nil {
		return false
	}
	return b.Breaking
}

func (self *breakerManage) Call(serviceMethod string, timeout bool) {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	var b *breaker
	if err := self.dict.Get(serviceMethod, &b); err != nil {
		return
	}
	b.CallSum += 1
	if timeout {
		b.CallTimeout += 1
	}
}

func (self *breakerManage) Run() {
	self.keys = self.dict.Keys()
	self.stop = make(chan bool)
	go func() {
		for {
			over := false
			select {
			case stop := <-self.stop:
				over = stop
				break
			default:
				self.Tick()
			}
			time.Sleep(1 * time.Second)
			if over {
				break
			}
		}
	}()
}

func (self *breakerManage) Stop() {
	self.stop <- true
}

func (self *breakerManage) Tick() {
	for index := 0; index < len(self.keys); index++ {
		self.breaker(self.keys[index])
	}
}

func (self *breakerManage) breaker(cmd interface{}) {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	var b *breaker
	if err := self.dict.Get(cmd, &b); err != nil {
		return
	}
	//熔断中
	if b.Breaking {
		b.ResumeTick -= 1
		if b.ResumeTick > 0 {
			return
		}
		b.Breaking = false
		b.CallSum = 0
		b.CallTimeout = 0
		b.Tick = b.Conf.RangeSec
		b.ResumeTick = 0
		return
	}

	//是否要熔断
	if b.CallSum > 0 && b.CallTimeout > 0 {
		r := float32(b.CallTimeout) / float32(b.CallSum)
		if r >= b.Conf.TimeoutRatio {
			//超过超时比例,触发熔断
			b.Breaking = true
			b.CallSum = 0
			b.CallTimeout = 0
			b.ResumeTick = b.Conf.ResumeSec
			b.Tick = 0
			return
		}
	}

	//正常中
	b.Tick -= 1
	if b.Tick > 0 {
		return
	}
	//重置
	b.CallSum = 0
	b.CallTimeout = 0
	b.Tick = b.Conf.RangeSec
	b.ResumeTick = 0
}

var _breakerManage breakerManage
