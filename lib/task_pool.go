package lib

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"jryghq.cn/utils"
)

var _gGoDict *GoDict

func NewTaskSimple() *TaskSimple {
	return new(TaskSimple).TaskSimple()
}

func NewTaskPool(length int) *TaskPool {
	return new(TaskPool).TaskPool(length)
}

func NewTaskGroup() *TaskGroup {
	return new(TaskGroup).TaskGroup()
}

func SetGoDict(godict *GoDict) {
	_gGoDict = godict
}

//单一task
type TaskSimple struct {
	TaskPool
}

//任务池
type TaskPool struct {
	queue     chan *task //队列管道
	close     bool       //是否关闭
	length    int32      //池大小
	group     bool       //是否是 group
	bseq      string     //bseq
	closeLock *sync.RWMutex
}

//组合task
type TaskGroup struct {
	TaskPool
	list []*task
}

//四种回调方式, 同步接口,异步接口,同步匿名函数,异步匿名函数
//同步接口
type CallBack interface {
	CallBack(result interface{}) error //同步调用等待回执result
}

//异步接口
type CastBack interface {
	CastBack() //异步调用
}

//任务
type task struct {
	call         CallBack     //同步方式
	cast         CastBack     //异步方式
	castFunction func()       //异步匿名函数方式
	callFunction func() error //同步匿名函数方式
	result       interface{}  //返回值
	done         chan error   //错误返回
}

//TaskPopl 构造
func (self *TaskPool) TaskPool(length int) *TaskPool {
	self.length = int32(length)
	self.queue = make(chan *task, length*16)
	self.close = false
	self.closeLock = new(sync.RWMutex)
	return self
}

//Start 启动 go程池
func (self *TaskPool) Start() {
	for index := int32(0); index < self.length; index++ {
		go self.run()
	}
}

//Run 处理调用
func (self *TaskPool) run() {
	iterator := new(Iterator).Iterator()
	if self.group && self.bseq != "" {
		_gGoDict.Set("bseq", self.bseq)
	}
	for {
		isRun := self.runTask()
		if self.IsClose() {
			break
		}

		//如果是 group 直接continue,无需sleep
		if self.group {
			continue
		}

		if !isRun {
			iterator.Sleep()
			continue
		}
		iterator.ResetTimer()
	}
	atomic.AddInt32(&self.length, -1)
	if self.group {
		_gGoDict.Remove()
	}
}

//runTask 处理task
func (self *TaskPool) runTask() bool {
	var t *task
	if self.group {
		//如果是 group 加入recover
		defer func() {
			if perr := recover(); perr != nil {
				utils.Log().WriteError(perr)
				if t != nil {
					//回执错误信息
					t.done <- errors.New(fmt.Sprint(perr))
				}
			}
		}()
	}
	select {
	case t = <-self.queue:
		if t.call != nil {
			t.done <- t.call.CallBack(t.result)
		} else if t.cast != nil {
			t.cast.CastBack()
		} else if t.callFunction != nil {
			t.done <- t.callFunction()
		} else if t.castFunction != nil {
			t.castFunction()
		}
		return true
	default:
		return false
	}
}

//Close 关闭
func (self *TaskPool) IsClose() bool {
	defer self.closeLock.RUnlock()
	self.closeLock.RLock()
	return self.close
}

//Close 关闭
func (self *TaskPool) Close() {
	//优雅关闭
	self.closeLock.Lock()
	self.close = true
	self.closeLock.Unlock()
}

//Call 同步调用
func (self *TaskPool) Call(call CallBack, result interface{}) error {
	//如果关闭 不接受调用
	if self.IsClose() {
		return errors.New("任务池:关闭")
	}

	ts := new(task)
	ts.call = call
	ts.result = result
	ts.done = make(chan error)
	self.queue <- ts

	//等待返回
	err := <-ts.done
	close(ts.done)
	return err
}

//Cast 异步调用
func (self *TaskPool) Cast(cast CastBack) error {
	//如果关闭 不接受调用
	if self.IsClose() {
		return errors.New("任务池:关闭")
	}

	ts := new(task)
	ts.cast = cast
	self.queue <- ts
	return nil
}

//CallFunc 同步匿名回调
func (self *TaskPool) CallFunc(fun func() error) error {
	if self.IsClose() {
		return errors.New("任务池:关闭")
	}

	ts := new(task)
	ts.callFunction = fun
	ts.done = make(chan error)
	self.queue <- ts

	//等待返回
	err := <-ts.done
	close(ts.done)
	return err
}

//CastFunc 异步匿名回调
func (self *TaskPool) CastFunc(fun func()) error {
	if self.IsClose() {
		return errors.New("任务池:关闭")
	}
	ts := new(task)
	ts.castFunction = fun
	self.queue <- ts
	return nil
}

//TaskPopl 构造
func (self *TaskSimple) TaskSimple() *TaskSimple {
	self.TaskPool.TaskPool(1)
	self.Start()
	return self
}

//TaskGroup 构造
func (self *TaskGroup) TaskGroup() *TaskGroup {
	self.list = make([]*task, 0, 5)
	self.TaskPool.TaskPool(5)
	self.group = true
	return self
}

//CallFAddCallFuncunc 同步匿名回调
func (self *TaskGroup) AddFuncByGroup(fun func() error) error {
	if self.IsClose() {
		return errors.New("任务池:关闭")
	}

	ts := new(task)
	ts.callFunction = fun
	ts.done = make(chan error)

	self.list = append(self.list, ts)
	return nil
}

//Wait 等待所有任务执行完成
func (self *TaskGroup) WaitByGroup() error {
	//如果关闭 不接受调用
	if self.IsClose() {
		return nil
	}
	defer func() {
		go self.Close()
	}()

	if _gGoDict != nil {
		//如果有 _taskGoDict 并且有bseq 读取并设置
		bseq := _gGoDict.Get("bseq")
		if bseq != nil {
			str, ok := bseq.(string)
			if ok {
				self.bseq = str
			}
		}
	}

	//限制并发长度最大64
	if len(self.list) < 64 {
		self.length = int32(len(self.list))
	}
	self.queue = make(chan *task, len(self.list))

	self.Start()
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
