package dao

import (
	"reflect"

	"jryghq.cn/dao/orm"
	"jryghq.cn/remote_call"
)

//DaoController
type DaoController struct {
	remote_call.RpcController
	selfName string
}

//DaoController
func (self *DaoController) DaoController(child interface{}) *DaoController {
	self.RpcController.RpcController(self)
	self.AddChild(self, child)

	type fun interface {
		RpcName() string
	}
	self.selfName = self.TopChild().(fun).RpcName()
	return self
}

//Model 服务定位器获取model
func (self *DaoController) Model(child interface{}) {
	modelDao := reflect.ValueOf(child).Elem().Interface().(daoName).Dao()
	if self.selfName != modelDao {
		panic("model 不在一个dao下,不要乱调用")
	}

	err := _msl.Service(child)
	if err != nil {
		panic("禁止调用:" + err.Error())
	}
	return
}

//Cache 服务定位器获取Cache
func (self *DaoController) Cache(child interface{}) {
	cacheDao := reflect.ValueOf(child).Elem().Interface().(daoName).Dao()
	if self.selfName != cacheDao {
		panic("cache不在一个dao下,不要乱调用")
	}

	err := _csl.Service(child)
	if err != nil {
		panic("禁止调用:" + err.Error())
	}
	return
}

//Api 服务定位器获取Api
func (self *DaoController) Api(child interface{}) {
	apiDao := reflect.ValueOf(child).Elem().Interface().(daoName).Dao()
	if self.selfName != apiDao {
		panic("api不在一个dao下,不要乱调用")
	}

	err := _api.Service(child)
	if err != nil {
		panic("禁止调用:" + err.Error())
	}
	return
}

//Memory 服务定位器获取Memory
func (self *DaoController) Memory(child interface{}) {
	apiDao := reflect.ValueOf(child).Elem().Interface().(daoName).Dao()
	if self.selfName != apiDao {
		panic("Memory不在一个dao下,不要乱调用")
	}

	err := _esl.Service(child)
	if err != nil {
		panic("禁止调用:" + err.Error())
	}
	return
}

//Transaction 事务
func (self *DaoController) Transaction(fun func() error) error {
	return orm.Transaction(self.selfName, fun)
}

//TotalPage 总页数
func (self *DaoController) TotalPage(size, pageSize int) int {
	if size%pageSize == 0 {
		return size / pageSize
	}
	return size/pageSize + 1
}
