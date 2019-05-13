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

package dao

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"time"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
)

var tp *lib.LimiteGo

type Request func(request *helper.HTTPRequest)

type DaoApi struct {
	lib.Object
	open    bool
	apiName string
	host    string

	dayMax     int
	hourMax    int
	minMax     int
	dayCount   int
	hourCount  int
	minCount   int
	countMutex sync.Mutex
}

func (self *DaoApi) Gotree(child interface{}) *DaoApi {
	self.Object.Gotree(self)
	self.AddChild(self, child)
	self.apiName = ""

	self.dayMax = 0
	self.hourMax = 0
	self.minMax = 0
	self.dayCount = 0
	self.hourCount = 0
	self.minCount = 0
	self.apiOn()
	return self
}

//TestOn 单元测试 开启
func (self *DaoApi) TestOn() {
	mode := helper.Config().String("sys::Mode")
	if mode == "prod" {
		helper.Log().WriteError("生产环境不可以使用单元测试api")
		panic("生产环境不可以使用单元测试api")
	}
	self.apiOn()
}

type apiName interface {
	Api() string
}

//apiOn
func (self *DaoApi) apiOn() {
	self.open = true
	self.apiName = self.TopChild().(apiName).Api()
	self.host = helper.Config().String("api::" + self.apiName)
}

//HttpGet
func (self *DaoApi) HttpGet(apiAddr string, args map[string]interface{}, callback ...Request) (result []byte, e error) {
	req := helper.HttpGet(self.host + apiAddr + "?" + httpBuildQuery(args))
	req = req.SetTimeout(3*time.Second, 10*time.Second)
	fun := func() error {
		var reqerr error
		result, reqerr = req.Bytes()
		return reqerr
	}

	if len(callback) > 0 {
		callback[0](req)
	}

	if self.limited() {
		e = self.count()
		if e != nil {
			return
		}
	}
	mode := helper.Config().String("sys::Mode")
	//2次重试
	for index := 0; index < 2; index++ {
		e = tp.Go(fun)
		if mode == "dev" {
			helper.Log().WriteInfo("ApiGet:", self.host+apiAddr, "ReqData:", args, "ResData:", string(result), "error:", e)
		}
		if e != nil {
			continue
		}

		break
	}

	return
}

//HttpPost
func (self *DaoApi) HttpPost(apiAddr string, args map[string]interface{}, callback ...Request) (result []byte, e error) {
	req := helper.HttpPost(self.host + apiAddr)
	req = req.SetTimeout(3*time.Second, 10*time.Second)
	for k, v := range args {
		req.Param(k, fmt.Sprint(v))
	}

	if len(callback) > 0 {
		callback[0](req)
	}

	if self.limited() {
		e = self.count()
		if e != nil {
			return
		}
	}

	fun := func() error {
		var reqerr error
		result, reqerr = req.Bytes()
		return reqerr
	}

	mode := helper.Config().String("sys::Mode")
	//2次重试
	for index := 0; index < 2; index++ {
		e = tp.Go(fun)
		if mode == "dev" {
			helper.Log().WriteInfo("ApiPost:", self.host+apiAddr, "ReqData:", args, "ResData:", string(result), "error:", e)
		}
		if e != nil {
			continue
		}

		break
	}

	return
}

//HttpPostJson
func (self *DaoApi) HttpPostJson(apiAddr string, raw interface{}, callback ...Request) (result []byte, e error) {
	req := helper.HttpPost(self.host + apiAddr)
	req = req.SetTimeout(3*time.Second, 10*time.Second)
	req.JSONBody(raw)

	if len(callback) > 0 {
		callback[0](req)
	}

	if self.limited() {
		e = self.count()
		if e != nil {
			return
		}
	}

	fun := func() error {
		var reqerr error
		result, reqerr = req.Bytes()
		return reqerr
	}

	mode := helper.Config().String("sys::Mode")
	//2次重试
	for index := 0; index < 2; index++ {
		e = tp.Go(fun)
		if mode == "dev" {
			helper.Log().WriteInfo("ApiPostJson:", self.host+apiAddr, "ReqData:", raw, "ResData:", string(result), "error:", e)
		}
		if e != nil {
			continue
		}

		break
	}

	return
}

//HostAddr 获取本api dao的host地址
func (self *DaoApi) HostAddr() string {
	return self.host
}

//limited api限制
func (self *DaoApi) count() error {
	defer self.countMutex.Unlock()
	self.countMutex.Lock()
	if self.minMax > 0 && self.minCount >= self.minMax {
		return errors.New("每分钟调用频次超过限制, host:" + self.host)
	}
	if self.hourMax > 0 && self.hourCount >= self.hourMax {
		return errors.New("每小时调用频次超过限制, host:" + self.host)
	}
	if self.dayMax > 0 && self.dayCount >= self.dayMax {
		return errors.New("每天调用频次超过限制, host:" + self.host)
	}

	if self.minMax > 0 {
		self.minCount += 1
	}
	if self.hourMax > 0 {
		self.hourCount += 1
	}
	if self.dayMax > 0 {
		self.dayCount += 1
	}
	return nil
}

//limited api限制是否开启
func (self *DaoApi) limited() bool {
	if self.dayMax > 0 || self.minMax > 0 || self.hourMax > 0 {
		return true
	}
	return false
}

//DayCountLimit 每日限制
func (self *DaoApi) DayCountLimit(count int) {
	if self.dayMax > 0 {
		return
	}
	self.dayMax = count
	lib.RunDefaultTimer(24*60*60*1000, func() {
		self.countMutex.Lock()
		self.dayCount = 0
		self.countMutex.Unlock()
	})
}

//HourCountLimit 每小时限制
func (self *DaoApi) HourCountLimit(count int) {
	if self.hourMax > 0 {
		return
	}
	self.hourMax = count
	lib.RunDefaultTimer(60*60*1000, func() {
		self.countMutex.Lock()
		self.hourCount = 0
		self.countMutex.Unlock()
	})
}

//MinCountLimit 每分钟限制
func (self *DaoApi) MinCountLimit(count int) {
	if self.minMax > 0 {
		return
	}
	self.minMax = count
	lib.RunDefaultTimer(60*1000, func() {
		self.countMutex.Lock()
		self.minCount = 0
		self.countMutex.Unlock()
	})
}

//HttpBuildQuery转换get参数
func httpBuildQuery(args map[string]interface{}) string {
	result := ""
	for k, v := range args {
		result += k + "=" + fmt.Sprint(v) + "&"
	}

	return url.PathEscape(strings.TrimSuffix(result, "&"))
}
