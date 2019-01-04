package lib

import (
	"sync"
)

type handlerFunc func(args ...interface{})

//观察者基类
type ObServer struct {
	handleMap map[interface{}]handlerFunc
	mutext    sync.RWMutex
}

//初始化
func (self *ObServer) Gotree() *ObServer {
	self.handleMap = make(map[interface{}]handlerFunc)
	return self
}

//加入订阅者
func (self *ObServer) AddSubscribe(o interface{}, handle handlerFunc) {
	defer self.mutext.Unlock()
	self.mutext.Lock()
	self.handleMap[o] = handle
}

//删除订阅者
func (self *ObServer) DelSubscribe(o interface{}) {
	defer self.mutext.Unlock()
	self.mutext.Lock()
	delete(self.handleMap, o)
}

//发布
func (self *ObServer) NotitySubscribe(args ...interface{}) {
	list := make([]handlerFunc, 0, len(self.handleMap))
	self.mutext.RLock()
	for _, handle := range self.handleMap {
		list = append(list, handle)
	}
	self.mutext.RUnlock()

	for index := 0; index < len(list); index++ {
		list[index](args...)
	}
}

//订阅者长度
func (self *ObServer) SubscribeLen() int {
	defer self.mutext.RUnlock()
	self.mutext.RLock()
	return len(self.handleMap)
}
