package lib

import (
	"jryghq.cn/utils"
	"reflect"
)

//组件基类
type Components struct {
	dict *Dict
}

//初始化
func (self *Components) Components() *Components {
	self.dict = new(Dict).Dict()
	return self
}

type enterComponent interface {
	EnterComponent()
}

type updateComponent interface {
	UpdateComponent(c *Components)
}

//加入组件
func (self *Components) AddComponent(obj interface{}) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		utils.Log().WriteError("AddComponent != reflect.Ptr")
	}
	self.dict.Set(t.Elem().Name(), obj)
	if app, ok := obj.(updateComponent); ok {
		app.UpdateComponent(self)
	}
	if app, ok := obj.(enterComponent); ok {
		app.EnterComponent()
	}
}

//移除组件
func (self *Components) RemoveComponent(obj interface{}) {
	t := reflect.TypeOf(obj)
	self.dict.Del(t.Name())
}

//获取组件
func (self *Components) GetComponent(obj interface{}) error {
	t := reflect.TypeOf(obj)
	return self.dict.Get(t.Elem().Elem().Name(), obj)
}

//广播组件内所有实现method的方法
func (self *Components) Broadcast(method string, arg interface{}) {
	list := self.dict.Keys()
	for _, v := range list {
		com := self.dict.GetInterface(v)
		if com == nil {
			continue
		}
		value := reflect.ValueOf(com).MethodByName(method)
		if value.Kind() != reflect.Invalid {
			value.Call([]reflect.Value{reflect.ValueOf(arg)})
		}
	}
}
