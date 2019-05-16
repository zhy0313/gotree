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
	"math"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/8treenet/gotree/helper"
)

//系统启动时间
var sysStartTime time.Time
var currentTimeNum int32 = 0
var _gGoDict *GoDict

func SetGoDict(godict *GoDict) {
	_gGoDict = godict
}
func init() {
	sysStartTime = time.Now()
	x := int64(rand.Intn(10000))
	identification = strconv.FormatInt(x, 36)
	tseq = 1
}

func CurrentTimeNum() int {
	var result int32
	result = atomic.LoadInt32(&currentTimeNum)
	return int(result)
}

func RunTickStopTimer(tick int64, child TickStopRun) {
	t := timer{
		tickStopfunc: child,
		n:            tick,
	}
	go t.tickLoop()
}

type StopTick interface {
	Stop()
}

func RunTick(tick int64, child func(), name string, delay int) StopTick {
	t := timer{
		tickfunc: child,
		n:        tick,
	}
	t.stop = make(chan bool, 1)
	if delay > 0 {
		go func() {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			t.triggerLoop(name)
		}()
	} else {
		go t.triggerLoop(name)
	}
	return &t
}

func RunDefaultTimer(tick int64, child func()) {
	t := timer{
		timefunc: child,
		n:        tick,
	}
	go t.defaultLoop()
}

func RunDay(hour, minute int, child func(), name string) {
	t := timer{
		timeDay: child,
		n:       0,
	}
	go t.dayLoop(hour, minute, name)
}

type TickStopRun func(stop *bool)

//定时器基类
type timer struct {
	timefunc     func()
	tickStopfunc TickStopRun
	tickfunc     func()
	timeDay      func()
	n            int64 //每秒或者指定时间和日期
	stop         chan bool
}

func (self *timer) Stop() {
	self.stop <- true
}

func (self *timer) triggerLoop(name string) {
	atomic.AddInt32(&currentTimeNum, 1)
	defer func() {
		atomic.AddInt32(&currentTimeNum, -1)
		if perr := recover(); perr != nil {
			helper.Log().WriteError(perr)
		}
		_gGoDict.Remove()
	}()

	n := self.n
	if n > 1000 {
		n = 1000
	}
	var deltaTime int64 = int64(0)

	for {
		isBreak := false
		select {
		case stop := <-self.stop:
			isBreak = stop
		default:
			isBreak = false
		}
		if isBreak {
			break
		}

		if deltaTime >= self.n {
			_gGoDict.Set("gseq", "t-"+getTseq())
			self.tickfunc()
			deltaTime = 0
		}

		time.Sleep(time.Duration(n) * time.Millisecond)
		deltaTime += n
	}
}

func (self *timer) tickLoop() {
	defer func() {
		if perr := recover(); perr != nil {
			helper.Log().WriteError(perr)
		}
	}()
	for {
		stop := false
		self.tickStopfunc(&stop)
		if stop {
			break
		}
		time.Sleep(time.Duration(self.n) * time.Millisecond)
	}
}

func (self *timer) defaultLoop() {
	defer func() {
		if perr := recover(); perr != nil {
			helper.Log().WriteError(perr)
		}
	}()
	for {
		time.Sleep(time.Duration(self.n) * time.Millisecond)
		self.timefunc()
	}
}

func (self *timer) dayLoop(hour, minute int, name string) {
	defer func() {
		if perr := recover(); perr != nil {
			helper.Log().WriteError(perr)
		}
		_gGoDict.Remove()
	}()

	currentHour := time.Now().Hour()
	currentMinute := time.Now().Minute()
	var n int64
	if currentHour > hour || (currentHour == hour && currentMinute >= minute) {
		td := time.Now().AddDate(0, 0, 1)
		local, _ := time.LoadLocation("Local")
		nt, _ := time.ParseInLocation("2006-01-02", td.Format("2006-01-02"), local)
		n = nt.Unix() + int64(hour*60*60+(minute*60))
	} else {
		td := time.Now()
		local, _ := time.LoadLocation("Local")
		nt, _ := time.ParseInLocation("2006-01-02", td.Format("2006-01-02"), local)
		n = nt.Unix() + int64(hour*60*60+(minute*60))
	}

	for {
		unix := time.Now().Unix()
		diff := n - unix
		diffTick := int(float32(diff) * float32(0.01))
		if unix > n {
			_gGoDict.Set("gseq", "t-"+getTseq())
			self.timeDay()
			break
		}
		time.Sleep(time.Duration(diffTick) * time.Second)
	}

	time.Sleep(5 * time.Second)
	self.dayLoop(hour, minute, name)
}

//GetSysClock 获取程序执行时间
func GetSysClock() int64 {
	return time.Now().Unix() - sysStartTime.Unix()
}

var tseq int64
var identification string
var tseqMutex sync.Mutex

func getTseq() (result string) {
	defer tseqMutex.Unlock()
	tseqMutex.Lock()
	result = identification
	result += strconv.FormatInt(tseq, 36)
	if tseq == math.MaxInt64 {
		tseq = 1
		return
	}
	tseq += 1
	return
}
