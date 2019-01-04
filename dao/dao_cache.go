package dao

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"jryghq.cn/dao/redis"
	"jryghq.cn/lib"
	"jryghq.cn/utils"
)

type DaoCache struct {
	lib.Object
	open    bool
	daoName string
}

func (self *DaoCache) DaoCache(child interface{}) *DaoCache {
	self.Object.Object(self)
	self.AddChild(self, child)
	self.AddSubscribe("DaoTelnet", self.daoTelnet)
	self.AddSubscribe("CacheOn", self.cacheOn)
	self.daoName = ""
	return self
}

//TestOn 单元测试 开启
func (self *DaoCache) TestOn() {
	mode := utils.Config().String("sys::mode")
	if mode == "prod" {
		utils.Log().WriteError("生产环境不可以使用单元测试cache")
		panic("生产环境不可以使用单元测试cache")
	}
	self.DaoInit()
	self.redisOn()
}

//daoOn 开启回调
func (self *DaoCache) daoTelnet(args ...interface{}) {
	dao := self.TopChild().(daoName)
	daoName := dao.Dao()

	for _, arg := range args {
		dao := arg.(daoNode)
		if dao.Name == daoName {
			self.redisOn()
			return
		}
	}
}

//cacheOn 开启回调
func (self *DaoCache) cacheOn(arg ...interface{}) {
	daoName := arg[0].(string)
	if daoName == self.daoName {
		self.redisOn()
	}
}

//daoOn 开启回调
func (self *DaoCache) redisOn() {
	self.open = true
	if !connectDao(self.daoName + "cache") {
		return
	}
	redisinfo := utils.Config().String("redis::" + self.daoName)
	if redisinfo == "" {
		utils.Log().WriteError("配置文件 dao:" + self.daoName + "redis地址错误或未找到")
	}
	list := strings.Split(redisinfo, ";")
	m := map[string]string{}
	for _, item := range list {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			utils.Log().WriteError("配置文件 dao:" + self.daoName + "redis地址错误")
			continue
		}
		m[kv[0]] = kv[1]
	}

	client := redis.GetClient(self.daoName)
	if client != nil {
		//已注册
		return
	}

	maxIdleConns := utils.Config().String("redis::" + self.daoName + "MaxIdleConns")
	maxOpenConns := utils.Config().String("redis::" + self.daoName + "MaxOpenConns")
	if maxIdleConns == "" {
		maxIdleConns = utils.Config().DefaultString("sys::RedisMaxIdleConns", "1")
	}
	if maxOpenConns == "" {
		maxOpenConns = utils.Config().DefaultString("sys::RedisMaxOpenConns", "2")
	}
	imaxIdleConns, ei := strconv.Atoi(maxIdleConns)
	imaxOpenConns, eo := strconv.Atoi(maxOpenConns)
	if ei != nil || eo != nil || imaxIdleConns == 0 || imaxOpenConns == 0 || imaxIdleConns > imaxOpenConns {
		utils.Log().WriteError("连接dao redis:"+self.daoName+"失败, 错误原因: MaxIdleConns或MaxOpenConns 参数错误,", imaxIdleConns, imaxOpenConns)
		panic("严重错误")
	}

	db, _ := strconv.Atoi(m["database"])
	utils.Log().WriteInfo("jryg connect redis: MaxIdleConns:" + maxIdleConns + "MaxOpenConns:" + maxOpenConns + " config:" + fmt.Sprint(m))
	client, e := redis.NewCache(m["server"], m["password"], db, imaxIdleConns, imaxOpenConns)
	if e != nil {
		utils.Log().WriteError(e)
		panic(e.Error())
	}
	redis.AddDatabase(self.daoName, client)
}

//Do
func (self *DaoCache) Do(cmd string, args ...interface{}) (reply interface{}, e error) {
	if !self.open {
		utils.Log().WriteError("cache error: 未开启dao:" + self.daoName)
		panic("cache error: 未开启dao:" + self.daoName)
	}
	if self.daoName == "" {
		utils.Log().WriteError("这是一个未注册的dao")
		panic("这是一个未注册的dao")
	}

	reply, e = redis.Do(self.daoName, cmd, args...)
	return
}

//Set
func (self *DaoCache) Set(cmd string, args ...interface{}) error {
	if !self.open {
		utils.Log().WriteError("cache error: 未开启dao:" + self.daoName)
		panic("cache error: 未开启dao:" + self.daoName)
	}
	if self.daoName == "" {
		utils.Log().WriteError("这是一个未注册的dao")
		panic("这是一个未注册的dao")
	}

	size := len(args)
	newargs, _ := json.Marshal(args[size-1])
	args[size-1] = string(newargs)
	_, e := redis.Do(self.daoName, cmd, args...)
	return e
}

//Get
func (self *DaoCache) Get(out interface{}, cmd string, args ...interface{}) (bool, error) {
	if !self.open {
		utils.Log().WriteError("cache error: 未开启dao:" + self.daoName)
		panic("cache error: 未开启dao:" + self.daoName)
	}
	if self.daoName == "" {
		utils.Log().WriteError("这是一个未注册的dao")
		panic("这是一个未注册的dao")
	}
	value, e := redis.Do(self.daoName, cmd, args...)
	if e != nil {
		return false, e
	}

	if value == nil {
		return false, nil
	}
	return true, json.Unmarshal(value.([]byte), out)
}

func (self *DaoCache) Connections(m map[string]int) {
	if !self.open {
		return
	}
	n := redis.GetConnects(self.daoName)
	m[self.daoName] = n
}

func (self *DaoCache) DaoInit() {
	if self.daoName == "" {
		self.daoName = self.TopChild().(daoName).Dao()
	}
}
