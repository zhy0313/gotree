package remote_call

import (
	"encoding/json"
	"sync"
	"time"

	"jryghq.cn/utils"

	"jryghq.cn/lib"
)

const (
	_CHECK_QPS_RESET = 10800000 //3小时重置
)

type qps struct {
	Count int64 //调用次数
	AvgMs int64 //平均用时
	MaxMs int64 //最高用时
	MinMs int64 //最低用时
}

type RpcQps struct {
	lib.Object
	mutex     sync.Mutex
	dict      map[string]*qps
	beginTime int64
}

func (self *RpcQps) RpcQps() *RpcQps {
	self.Object.Object(self)
	self.dict = make(map[string]*qps)
	self.beginTime = time.Now().Unix()
	lib.RunTickStopTimer(_CHECK_QPS_RESET, self.tick) //定时器检测超时节点
	self.AddSubscribe("DaoQps", self.list)
	return self
}

func (self *RpcQps) Qps(serviceMethod string, ms int64) {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	if ms < 0 {
		utils.Log().WriteError("RpcQps ms < 0 ServiceMethod:", serviceMethod)
		return
	}

	var _qps *qps
	dqps, ok := self.dict[serviceMethod]
	if ok {
		_qps = dqps
	} else {
		_qps = new(qps)
		self.dict[serviceMethod] = _qps
	}

	if ms == 0 {
		ms = 1
	}

	if _qps.Count == 0 {
		_qps.Count = 1
		_qps.AvgMs = ms
		_qps.MaxMs = ms
		_qps.MinMs = ms
		return
	}

	_qps.Count += 1
	_qps.AvgMs = (_qps.AvgMs + ms) / 2
	if ms > _qps.MaxMs {
		_qps.MaxMs = ms
	}
	if ms < _qps.MinMs {
		_qps.MinMs = ms
	}
}

func (self *RpcQps) tick(stop *bool) {
	var list []struct {
		ServiceMethod string
		Count         int64 //调用次数
		AvgMs         int64 //平均用时
		MaxMs         int64 //最高用时
		MinMs         int64 //最低用时
	}

	self.list(&list)
	if len(list) > 0 {
		data, e := json.Marshal(list)
		if e == nil {
			utils.Log().WriteInfo("business qps", string(data))
		}
	}

	self.mutex.Lock()
	self.beginTime = time.Now().Unix()
	self.dict = make(map[string]*qps)
	self.mutex.Unlock()
}

func (self *RpcQps) list(args ...interface{}) {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	ret := args[0].(*[]struct {
		ServiceMethod string
		Count         int64 //调用次数
		AvgMs         int64 //平均用时
		MaxMs         int64 //最高用时
		MinMs         int64 //最低用时
	})

	list := make([]struct {
		ServiceMethod string
		Count         int64 //调用次数
		AvgMs         int64 //平均用时
		MaxMs         int64 //最高用时
		MinMs         int64 //最低用时
	}, 0)

	for key, item := range self.dict {
		var additem struct {
			ServiceMethod string
			Count         int64 //调用次数
			AvgMs         int64 //平均用时
			MaxMs         int64 //最高用时
			MinMs         int64 //最低用时
		}

		additem.ServiceMethod = key
		additem.Count = item.Count
		additem.AvgMs = item.AvgMs
		additem.MaxMs = item.MaxMs
		additem.MinMs = item.MinMs
		list = append(list, additem)
	}
	*ret = list
	return
}
