package remote_call

import (
	"math/rand"
	"net/rpc"
	"sync"
	"time"

	"github.com/8treenet/gotree/lib"
)

var connPool *rpcConnPool

func init() {
	connPool = new(rpcConnPool).Gotree()
}

type rpcConn struct {
	lib.Object
	client    *rpc.Client
	takeCount int  //当前使用数量
	status    int8 //0创建, 1连接成功, 2网络错误
	addr      string
	id        int
	mutex     sync.Mutex
	exit      chan bool
}

func (self *rpcConn) Gotree(addr string, id int) *rpcConn {
	self.Object.Gotree(self)
	self.status = 0
	self.addr = addr
	self.id = id
	self.exit = make(chan bool, 1)
	return self
}

func (self *rpcConn) connect() (e error) {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	self.client, e = jsonRpc(self.addr)
	if e != nil {
		self.status = 2
		return
	}
	self.status = 1
	return
}

func (self *rpcConn) timeout() {
	//一定连接成功过才会进入
	for index := 0; index < 240; index++ {
		time.Sleep(500 * time.Millisecond)
		if self.statu() != 1 {
			self.client.Close()
			connPool.delConn(self.addr, self.id)
			return
		}
		select {
		case _ = <-self.exit:
			self.client.Close()
			connPool.delConn(self.addr, self.id)
			return
		default:
		}
	}

	//结束判断是否还有连接 等待30秒
	connPool.delConn(self.addr, self.id)
	for index := 0; index < 60; index++ {
		if self.getTakeCount() <= 0 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	self.client.Close()
	return
}

func (self *rpcConn) take() {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	self.takeCount += 1
	return
}

func (self *rpcConn) release(netok bool) {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	self.takeCount -= 1
	if !netok && self.client != nil && self.status == 1 {
		self.status = 2
	}
	return
}

func (self *rpcConn) statu() int8 {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	return self.status
}

func (self *rpcConn) getTakeCount() int {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	return self.takeCount
}

type rpcConnPool struct {
	lib.Object
	pool  map[string]map[int]*rpcConn
	mutex sync.Mutex
}

func (self *rpcConnPool) Gotree() *rpcConnPool {
	self.Object.Gotree(self)
	self.pool = make(map[string]map[int]*rpcConn)
	return self
}

func (self *rpcConnPool) takeConn(addr string) (*rpcConn, error) {
	var conn *rpcConn
	m := self.addrMap(addr)
	for {
		i := self.random()
		self.mutex.Lock()
		var ok bool
		conn, ok = m[i]
		if ok {
			if conn.statu() != 1 {
				conn = nil
			}
		} else {
			conn = new(rpcConn).Gotree(addr, i)
			m[i] = conn
		}
		self.mutex.Unlock()
		if conn != nil {
			break
		}
	}

	if conn.statu() == 0 {
		err := conn.connect()
		if err != nil {
			self.delConn(addr, conn.id)
			return conn, err
		}
		go conn.timeout()
	}
	conn.take()
	return conn, nil
}

var maxConnCount int = 8

func (self *rpcConnPool) random() int {
	return 1 + rand.Intn(maxConnCount)
}

func (self *rpcConnPool) addrMap(addr string) map[int]*rpcConn {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	result, ok := self.pool[addr]
	if ok {
		return result
	}
	result = make(map[int]*rpcConn)
	self.pool[addr] = result
	return result
}

func (self *rpcConnPool) delConn(addr string, id int) {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	result, ok := self.pool[addr]
	if !ok {
		return
	}
	delete(result, id)
	return
}

func (self *rpcConnPool) Close() {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	for _, item := range self.pool {
		for _, conn := range item {
			conn.exit <- true
		}
	}
	return
}
