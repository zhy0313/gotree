package lib

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type LimiteGo struct {
	count    int64 //当前go数量
	maxCount int64
}

//LimiteGo
func (self *LimiteGo) Gotree(maxCount int) *LimiteGo {
	self.count = 0
	self.maxCount = int64(maxCount)
	return self
}

func (self *LimiteGo) Go(fun func() error) (e error) {
	for index := 0; index < 100; index++ {
		if (atomic.LoadInt64(&self.count)) > self.maxCount {
			time.Sleep(30 * time.Millisecond)
			continue
		}
		break
	}

	defer func() {
		if perr := recover(); perr != nil {
			e = errors.New(fmt.Sprint(perr))
		}
		atomic.AddInt64(&self.count, -1)
	}()
	atomic.AddInt64(&self.count, 1)
	e = fun()
	return
}
