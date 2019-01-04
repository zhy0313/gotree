package dao

import (
	"os"
	"sync"
	"time"

	"jryghq.cn/lib"
	"jryghq.cn/utils"
)

//DaoMemory 内存数据源
type DaoMemory struct {
	lib.Object
	daoName   string
	open      bool
	dict      *lib.Dict
	dictMutex sync.Mutex
	dictEep   map[interface{}]int64
}

func (self *DaoMemory) DaoMemory(child interface{}) *DaoMemory {
	self.Object.Object(self)
	self.AddChild(self, child)
	self.daoName = ""
	self.AddSubscribe("MemoryOn", self.memoryOn)

	self.dict = new(lib.Dict).Dict()
	self.dictEep = make(map[interface{}]int64)
	return self
}

//TestOn 单元测试 开启
func (self *DaoMemory) TestOn() {
	mode := utils.Config().String("sys::mode")
	if mode == "prod" {
		utils.Log().WriteError("生产环境不可以使用单元测试model")
		os.Exit(1)
	}
	self.DaoInit()
	self.open = true
}

//daoOn 开启回调
func (self *DaoMemory) memoryOn(arg ...interface{}) {
	daoName := arg[0].(string)
	if daoName == self.daoName {
		self.open = true
	}
}

// Set
func (self *DaoMemory) Set(key, value interface{}) {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	self.dict.Set(key, value)
	return
}

// Get 不存在 返回false
func (self *DaoMemory) Get(key, value interface{}) bool {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	e := self.dict.Get(key, value)
	if e != nil {
		return false
	}
	return true
}

// SetTnx 当key不存在设置成功返回true 否则返回false
func (self *DaoMemory) SetTnx(key, value interface{}) bool {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	if self.dict.Check(key) {
		return false
	}
	self.dict.Set(key, value)
	return true
}

// MultiSet 多条
func (self *DaoMemory) MultiSet(args ...interface{}) {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	if len(args) <= 0 {
		panic("MultiSet len(args) <= 0")
	}
	//多参必须是偶数
	if (len(args) & 1) == 1 {
		panic("MultiSet len(args)&1 == 1")
	}

	for index := 0; index < len(args); index += 2 {
		self.dict.Set(args[index], args[index+1])
	}
	return
}

// MultiSet 多条
func (self *DaoMemory) MultiSetTnx(args ...interface{}) bool {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	if len(args) <= 0 {
		panic("MultiSet len(args) <= 0")
	}
	//多参必须是偶数
	if (len(args) & 1) == 1 {
		panic("MultiSet len(args)&1 == 1")
	}

	for index := 0; index < len(args); index += 2 {
		if self.dict.Check(args[index]) {
			return false
		}
	}
	for index := 0; index < len(args); index += 2 {
		self.dict.Set(args[index], args[index+1])
	}
	return true
}

// Eexpire 设置 key 的生命周期, sec:秒
func (self *DaoMemory) Eexpire(key interface{}, sec int) {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	if !self.dict.Check(key) {
		//不存在,直接返回
		return
	}
	self.dictEep[key] = time.Now().Unix() + int64(sec)
	return
}

// Delete 删除 key
func (self *DaoMemory) Delete(keys ...interface{}) {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	for _, key := range keys {
		self.dict.Del(key)
		delete(self.dictEep, key)
	}
}

// DeleteAll 删除全部数据
func (self *DaoMemory) DeleteAll(key interface{}) {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	self.dict.DelAll()
	self.dictEep = make(map[interface{}]int64)
}

// Incr add 加数据, key必须存在否则errror
func (self *DaoMemory) Incr(key interface{}, addValue int64) (result int64, e error) {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	e = self.dict.Get(key, &result)
	if e != nil {
		return
	}
	result += addValue
	self.dict.Set(key, result)
	return
}

// AllKey 获取全部key
func (self *DaoMemory) AllKey() (result []interface{}) {
	defer self.dictMutex.Unlock()
	self.dictMutex.Lock()
	result = self.dict.Keys()
	return
}

// MemoryTimeout 超时处理
func (self *DaoMemory) MemoryTimeout(now int64) {
	keys := []interface{}{}
	self.dictMutex.Lock()
	for k, v := range self.dictEep {
		if now < v {
			continue
		}
		keys = append(keys, k)
	}
	self.dictMutex.Unlock()

	for index := 0; index < len(keys); index++ {
		self.Delete(keys[index])
	}
	return
}

func (self *DaoMemory) DaoInit() {
	if self.daoName == "" {
		self.daoName = self.TopChild().(daoName).Dao()
	}
}
