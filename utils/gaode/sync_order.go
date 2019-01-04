package gaode

import (
	"time"
)

// SyncOrder
type SyncOrder struct {
	Key   string
	Type  string
	Time  int64
	Otype int // 订单类型	1:实时订单	2:预约订单	3:接机单 4:送机单 5:日租 6:半日租
}

// SyncOrder 初始化
func (s *SyncOrder) SyncOrder() *SyncOrder {
	s.Key = "79d13a44d4a4c738e1403de1d3e6a4b7"
	s.Type = "orders"
	s.Otype = 1
	return s
}

//GetTs 计算 Ts
func (s *SyncOrder) GetTs() *SyncOrder {
	timestamp := time.Now().Local().UnixNano() / 1e6
	s.Time = timestamp
	return s
}
