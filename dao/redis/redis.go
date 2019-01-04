package redis

import (
	"errors"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/8treenet/gotree/helper"
	"github.com/garyburd/redigo/redis"
)

var dbs map[string]*Cache

func init() {
	dbs = make(map[string]*Cache)
}

func Do(name, cmd string, to ...interface{}) (reply interface{}, err error) { //, to interface{}
	rc := GetClient(name)
	if rc == nil {
		return nil, errors.New("Redis 连接无效")
	}

	reply, err = rc.do(cmd, to...)
	return
}

// Cache is Redis cache adapter.
type Cache struct {
	p     *redis.Pool // redis connection pool
	mutex sync.Mutex
}

// AddDatabase 加入 枚举库 enum/db 和缓存
func AddDatabase(name string, c *Cache) {
	dbs[name] = c
}

// GetClient 通过枚举读取client
func GetClient(name string) *Cache {
	return dbs[name]
}

// GetConnects 获取连接数
func GetConnects(name string) int {
	c, ok := dbs[name]
	if !ok {
		return 0
	}
	return c.p.IdleCount() + c.p.ActiveCount()
}

// NewCache 创建新的 Redis Cache 对象。 connection - Redis服务器和端口号，例如: "127.0.0.1:6379", dbNum - Redis 数据库序号
func NewCache(connection string, password string, dbNum int, idle int, maxOpen int) (cache *Cache, err error) {
	cache = &Cache{}

	// initialize a new pool
	cache.p = &redis.Pool{
		MaxIdle:     idle,
		IdleTimeout: 600 * time.Second,
		MaxActive:   maxOpen,
		Dial: func() (c redis.Conn, err error) {
			c, err = redis.DialTimeout("tcp", connection, 3*time.Second, 15*time.Second, 15*time.Second)
			//c, err = redis.Dial("tcp", connection)
			if err != nil {
				return
			}

			defer func() {
				if err != nil {
					c.Close()
				}
			}()

			if password != "" {
				_, err = c.Do("AUTH", password)
				if err != nil {
					return
				}
			}

			if dbNum != 0 {
				_, err = c.Do("SELECT", dbNum)
				if err != nil {
					return
				}
			}

			return
		},
	}

	c := cache.p.Get()
	defer c.Close()
	err = c.Err()
	if err != nil {
		cache = nil
	}
	return
}

//GetRc 获取连接池中的链接
func (rc *Cache) getRc() (redis.Conn, error) {
	defer rc.mutex.Unlock()
	rc.mutex.Lock()
	u := time.Now().Unix()
	for {
		//当前正在活动中的链接 小于最大链接 = 有空闲链接
		if rc.p.ActiveCount() < rc.p.MaxActive {
			return rc.p.Get(), nil
		}

		//没有空闲链接且等待15秒 返回超时
		if time.Now().Unix()-u > 15 {
			break
		}
		runtime.Gosched()
	}
	return nil, errors.New("timeout")
}

// actually do the redis cmds
func (rc *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if strings.ToLower(commandName) == "select" {
		helper.Log().WriteError("运行中不可以切库")
	}

	c, err := rc.getRc()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return c.Do(commandName, args...)
}
