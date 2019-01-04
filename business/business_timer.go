package business

import (
	"jryghq.cn/lib"
	"jryghq.cn/utils"
)

type BusinessTimer struct {
	lib.Object
	timerCallBack []trigger
	timerStopTick []lib.StopTick
	_openService  bool
}

func (self *BusinessTimer) BusinessTimer(child interface{}) *BusinessTimer {
	self.Object.Object(self)
	self.AddChild(self, child)
	self.timerCallBack = make([]trigger, 0)
	self.AddSubscribe("TimerOn", self.timer)
	self.AddSubscribe("shutdown", self.stopTick)
	self._openService = false
	self.timerStopTick = make([]lib.StopTick, 0)
	return self
}

//RegisterTickTrigger 注册毫秒tick 触发器 delay=延迟毫秒
func (self *BusinessTimer) RegisterTickTrigger(ms int, fun func(), delay ...int) {
	delayms := 0
	if len(delay) > 0 {
		delayms = delay[0]
	}
	self.timerCallBack = append(self.timerCallBack, trigger{t: ms, tickFun: fun, delay: delayms})
	return
}

//RegisterDayTrigger 注册每天小时 触发器 hour:0 - 23, minute :0 - 59
func (self *BusinessTimer) RegisterDayTrigger(hour, minute int, fun func()) {
	if hour < 0 && hour > 23 {
		utils.Log().WriteError("hour:0 - 23")
	}

	if minute < 0 && minute > 59 {
		utils.Log().WriteError("minute :0 - 59")
	}
	self.timerCallBack = append(self.timerCallBack, trigger{t: hour, t2: minute, dayFun: fun})
	return
}

func (self *BusinessTimer) stopTick(args ...interface{}) {
	for _, stopFun := range self.timerStopTick {
		stopFun.Stop()
	}
}

//服务定位器获取服务
func (self *BusinessTimer) Service(child interface{}) {
	err := _scl.Service(child)
	if err != nil {
		utils.Log().WriteError("飞哥:不要乱调用:" + err.Error())
		panic("飞哥:不要乱调用:" + err.Error())
	}
	return
}

//timer 启动观察者
func (self *BusinessTimer) timer(args ...interface{}) {
	selfName := self.ClassName(self.TopChild())
	//获取配置文件开启的定时器服务
	open := false
	for _, namei := range args {
		name := namei.(string)
		if name == selfName {
			open = true
			break
		}
	}

	if !open {
		return
	}
	for _, fun := range self.timerCallBack {
		if fun.tickFun != nil {
			tick := lib.RunTick(int64(fun.t), fun.tickFun, selfName, fun.delay)
			self.timerStopTick = append(self.timerStopTick, tick)
		}

		if fun.dayFun != nil {
			lib.RunDay(fun.t, fun.t2, fun.dayFun, selfName)
		}
	}
	return
}

func (self *BusinessTimer) OpenTimer() {
	exist := _tsl.CheckService(self.TopChild())
	if exist {
		utils.Log().WriteError("禁止重复实例化")
		panic("禁止重复实例化")
	}
	self._openService = true
	return
}

//Async 异步执行
func (self *BusinessTimer) Async(run func(ac AsyncController), completeds ...func()) {
	var completed func()
	if len(completeds) > 0 {
		completed = completeds[0]
	}
	ac := new(async).async(run, completed)
	go ac.execute()
}

// Broadcast 调用所有注册service的method方法. 潜龙勿用,会使项目非常难以维护
func (self *BusinessTimer) Broadcast(method string, arg interface{}) {
	if e := _scl.Broadcast(method, arg); e != nil {
		utils.Log().WriteError("Buesiness ServiceBroadcast errror:" + e.Error())
	}
}

type trigger struct {
	t       int
	t2      int
	delay   int //延迟毫秒
	tickFun func()
	dayFun  func()
}
