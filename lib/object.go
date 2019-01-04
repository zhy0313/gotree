package lib

import (
	"errors"
	"reflect"
	"sync"
)

func init() {
	//全局观察者 继承Object可使用
	globalObserver = make(map[string]*ObServer)
}

//基类
type Object struct {
	component *Components //组件
	child     map[interface{}]interface{}
}

//初始化
func (self *Object) Gotree(child interface{}) *Object {
	self.child = make(map[interface{}]interface{})
	self.AddChild(self, child)
	self.component = new(Components).Gotree()
	return self
}

//修改组件
func (self *Object) UpdateComponent(component *Components) {
	self.component = component
}

//GetComponent
func (self *Object) GetComponent(obj interface{}) error {
	return self.component.GetComponent(obj)
}

//AddComponent
func (self *Object) AddComponent(obj interface{}) {
	self.component.AddComponent(obj)
}

//Broadcast
func (self *Object) Broadcast(method string, arg interface{}) {
	self.component.Broadcast(method, arg)
}

var globalObserver map[string]*ObServer
var globalObserverMutex sync.Mutex

//加入全局订阅者
func (self *Object) AddSubscribe(event string, handle handlerFunc) {
	defer globalObserverMutex.Unlock()
	globalObserverMutex.Lock()
	ob, ok := globalObserver[event]
	if ok {
		ob.AddSubscribe(self, handle)
		return
	}
	ob = new(ObServer).Gotree()
	ob.AddSubscribe(self, handle)
	globalObserver[event] = ob
}

//删除全局订阅者
func (self *Object) DelSubscribe(event string) {
	defer globalObserverMutex.Unlock()
	globalObserverMutex.Lock()
	ob, ok := globalObserver[event]
	if !ok {
		return
	}
	ob.DelSubscribe(self)
	if ob.SubscribeLen() == 0 {
		delete(globalObserver, event)
	}
}

//通知全局订阅
func (self *Object) NotitySubscribe(event string, args ...interface{}) {
	globalObserverMutex.Lock()
	ob, ok := globalObserver[event]
	globalObserverMutex.Unlock()
	if !ok {
		return
	}
	ob.NotitySubscribe(args...)
}

//通知全局订阅
func (self *Object) Observer(event string) *ObServer {
	globalObserverMutex.Lock()
	ob, ok := globalObserver[event]
	globalObserverMutex.Unlock()
	if ok {
		return ob
	}
	return nil
}

//AddChild 添加子类
func (self *Object) AddChild(parnet interface{}, child ...interface{}) {
	if len(child) == 0 {
		return
	}
	c := child[0]
	if c == nil {
		return
	}
	self.child[parnet] = c
}

//GetChild 获取子类
func (self *Object) GetChild(parnet interface{}) (child interface{}, err error) {
	err = nil
	child, ok := self.child[parnet]
	if !ok {
		err = errors.New("undefined")
	}
	return
}

//TopChild 获取顶级子类
func (self *Object) TopChild() (result interface{}) {
	result = self
	for {
		c, err := self.GetChild(result)
		if err != nil {
			return
		}
		result = c
	}
}

//ClassName 获取实例名字
func (self *Object) ClassName(class interface{}) (name string) {
	return reflect.TypeOf(class).Elem().Name()
}
