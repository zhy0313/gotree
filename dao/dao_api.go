package dao

import (
	"errors"
	"fmt"
	"sync"

	"time"

	"jryghq.cn/lib"
	"jryghq.cn/utils"
)

var tp *lib.TaskPool

type Request func(request *utils.HTTPRequest)

type DaoApi struct {
	lib.Object
	open    bool
	daoName string
	host    string

	dayMax     int
	hourMax    int
	minMax     int
	dayCount   int
	hourCount  int
	minCount   int
	countMutex sync.Mutex
}

func (self *DaoApi) DaoApi(child interface{}) *DaoApi {
	self.Object.Object(self)
	self.AddChild(self, child)
	self.AddSubscribe("DaoTelnet", self.daoTelnet)
	self.AddSubscribe("ApiOn", self.apiOn)
	self.daoName = ""

	self.dayMax = 0
	self.hourMax = 0
	self.minMax = 0
	self.dayCount = 0
	self.hourCount = 0
	self.minCount = 0
	return self
}

//TestOn 单元测试 开启
func (self *DaoApi) TestOn() {
	mode := utils.Config().String("sys::mode")
	if mode == "prod" {
		utils.Log().WriteError("生产环境不可以使用单元测试api")
		panic("生产环境不可以使用单元测试api")
	}
	self.DaoInit()
	self.apiOn()
}

//daoOn 开启回调
func (self *DaoApi) daoTelnet(args ...interface{}) {
	dao := self.TopChild().(daoName)
	daoName := dao.Dao()

	for _, arg := range args {
		dao := arg.(daoNode)
		if dao.Name == daoName {
			self.apiOn()
			b := utils.HttpGet(self.host)
			req := b.SetTimeout(1*time.Second, 1*time.Second)
			_, reqerr := req.Bytes()
			if reqerr != nil {
				utils.Log().WriteWarn("连接dao api:" + self.daoName + "失败, 错误原因:" + reqerr.Error())
			}
			return
		}
	}
}

//daoOn 开启回调
func (self *DaoApi) apiOn(arg ...interface{}) {
	if len(arg) > 0 {
		daoName := arg[0].(string)
		if daoName != self.daoName {
			return
		}
	}

	self.open = true
	self.host = utils.Config().String("api::" + self.daoName)
}

//HttpGet
func (self *DaoApi) HttpGet(apiAddr string, args map[string]interface{}, callback ...Request) (result []byte, e error) {
	req := utils.HttpGet(self.host + apiAddr + "?" + utils.HttpBuildQuery(args))
	req = req.SetTimeout(3*time.Second, 15*time.Second)
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
	mode := utils.Config().String("sys::mode")
	//3次重试
	for index := 0; index < 3; index++ {
		e = tp.CallFunc(fun)
		if mode == "dev" {
			utils.Log().WriteInfo("ApiGet:", self.host+apiAddr, "ReqData:", args, "ResData:", string(result), "error:", e)
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
	req := utils.HttpPost(self.host + apiAddr)
	req = req.SetTimeout(3*time.Second, 15*time.Second)
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

	mode := utils.Config().String("sys::mode")
	//3次重试
	for index := 0; index < 3; index++ {
		e = tp.CallFunc(fun)
		if mode == "dev" {
			utils.Log().WriteInfo("ApiPost:", self.host+apiAddr, "ReqData:", args, "ResData:", string(result), "error:", e)
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
	req := utils.HttpPost(self.host + apiAddr)
	req = req.SetTimeout(3*time.Second, 15*time.Second)
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

	mode := utils.Config().String("sys::mode")
	//3次重试
	for index := 0; index < 3; index++ {
		e = tp.CallFunc(fun)
		if mode == "dev" {
			utils.Log().WriteInfo("ApiPostJson:", self.host+apiAddr, "ReqData:", raw, "ResData:", string(result), "error:", e)
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

func (self *DaoApi) DaoInit() {
	if self.daoName == "" {
		self.daoName = self.TopChild().(daoName).Dao()
	}
}
