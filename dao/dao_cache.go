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
	"fmt"
	"strconv"
	"strings"

	"github.com/8treenet/gotree/dao/redis"
	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
)

type DaoCache struct {
	lib.Object
	open    bool
	daoName string
}

func (self *DaoCache) Gotree(child interface{}) *DaoCache {
	self.Object.Gotree(self)
	self.AddChild(self, child)
	self.AddSubscribe("DaoTelnet", self.daoTelnet)
	self.AddSubscribe("CacheOn", self.cacheOn)
	self.daoName = ""
	return self
}

//TestOn 单元测试 开启
func (self *DaoCache) TestOn() {
	mode := helper.Config().String("sys::Mode")
	if mode == "prod" {
		helper.Log().WriteError("生产环境不可以使用单元测试cache")
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
	redisinfo := helper.Config().String("redis::" + self.daoName)
	if redisinfo == "" {
		helper.Log().WriteError("配置文件 dao:" + self.daoName + "redis地址错误或未找到")
	}
	list := strings.Split(redisinfo, ";")
	m := map[string]string{}
	for _, item := range list {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			helper.Log().WriteError("配置文件 dao:" + self.daoName + "redis地址错误")
			continue
		}
		m[kv[0]] = kv[1]
	}

	client := redis.GetClient(self.daoName)
	if client != nil {
		//已注册
		return
	}

	maxIdleConns := helper.Config().String("redis::" + self.daoName + "MaxIdleConns")
	maxOpenConns := helper.Config().String("redis::" + self.daoName + "MaxOpenConns")
	if maxIdleConns == "" {
		maxIdleConns = helper.Config().DefaultString("sys::RedisMaxIdleConns", "1")
	}
	if maxOpenConns == "" {
		maxOpenConns = helper.Config().DefaultString("sys::RedisMaxOpenConns", "2")
	}
	imaxIdleConns, ei := strconv.Atoi(maxIdleConns)
	imaxOpenConns, eo := strconv.Atoi(maxOpenConns)
	if ei != nil || eo != nil || imaxIdleConns == 0 || imaxOpenConns == 0 || imaxIdleConns > imaxOpenConns {
		helper.Log().WriteError("连接dao redis:"+self.daoName+"失败, 错误原因: MaxIdleConns或MaxOpenConns 参数错误,", imaxIdleConns, imaxOpenConns)
		panic("严重错误")
	}

	db, _ := strconv.Atoi(m["database"])
	helper.Log().WriteInfo("connect redis: MaxIdleConns:" + maxIdleConns + "MaxOpenConns:" + maxOpenConns + " config:" + fmt.Sprint(m))
	client, e := redis.NewCache(m["server"], m["password"], db, imaxIdleConns, imaxOpenConns)
	if e != nil {
		helper.Log().WriteError(e)
		panic(e.Error())
	}
	redis.AddDatabase(self.daoName, client)
}

//Do
func (self *DaoCache) Do(cmd string, args ...interface{}) (reply interface{}, e error) {
	if !self.open || self.daoName == "" {
		helper.Log().WriteError("cache error: 未开启dao:" + self.daoName)
		panic("cache error: 未开启dao:" + self.daoName)
	}
	reply, e = redis.Do(self.daoName, cmd, args...)
	return
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
