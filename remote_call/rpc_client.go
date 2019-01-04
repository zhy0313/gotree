package remote_call

import (
	"errors"
	"fmt"
	"math/rand"
	"net/rpc"
	"reflect"
	"strings"
	"time"

	"jryghq.cn/lib"
	jrygrpc "jryghq.cn/lib/rpc"
	"jryghq.cn/utils"
)

type RpcClient struct {
	lib.Object
	task       *lib.TaskPool
	retry      int
	recoverLog bool
	connReuse  bool //连接池,连接复用 默认开启
}

type cmdCall interface {
	RemoteAddr(interface{}) (string, error) //远程地址 0.0.0.1:9123
	ServiceMethod() string                  //远程方法
	Cache() interface{}                     //如果开启缓存返回标识否则nil
}

type businessCmd interface {
	BusinessAddr() string //bussiness地址
}

//并发数量 和rpc失败重试次数
func (self *RpcClient) RpcClient(concurrency int, retry int) *RpcClient {
	self.Object.Object(self)
	self.task = lib.NewTaskPool(concurrency)
	self.task.Start()
	maxConnCount += concurrency / 64

	self.retry = 5
	if retry < 9 || retry > 0 {
		self.retry = retry
	}
	self.recoverLog = false
	self.connReuse = true
	return self
}

//Call cmd obj参数, reply &daostruct 或 &businessstruct
func (self *RpcClient) Call(obj interface{}, reply interface{}) (err error) {
	var innerMaster *InnerMaster
	var identity string
	var identityUse bool = false
	var beginMs int64

	self.GetComponent(&innerMaster)
	if innerMaster == nil {
		panic("未获取到InnerMaster")
	}
	cmd := obj.(cmdCall)
	serviceMethod := cmd.ServiceMethod()
	if !innerMaster.ping {
		beginMs = time.Now().UnixNano() / 1e6
	}

	defer func() {
		if !innerMaster.ping && !identityUse && (err == nil || err != unknownNetwork) {
			//无错误或非网络错误 并且 非cmd缓存 加入统计
			self.Qps(serviceMethod, time.Now().UnixNano()/1e6-beginMs)
		}
		if err != nil {
			if !innerMaster.ping {
				childObj, _ := obj.(cmdSerChild)
				err = errors.New("call dao " + childObj.Control() + "." + childObj.Action() + ", error:" + err.Error())
			}
			return
		}
		//如果未开启缓存 并且使用缓存获取的数据不处理
		if identity == "" && !identityUse {
			return
		}
		cacheValue := reflect.Indirect(reflect.ValueOf(reply)).Interface()
		jrygrpc.GoDict().Set(identity, cacheValue)
	}()

	cacheIdentity := cmd.Cache()
	if cacheIdentity != nil {
		identity = self.ClassName(cmd) + "_" + fmt.Sprint(cacheIdentity)
		if jrygrpc.GoDict().Eval(identity, reply) == nil {
			identityUse = true
			return
		}
	}

	fun := func() error {
		if self.recoverLog {
			defer func() {
				if perr := recover(); perr != nil {
					utils.Log().WriteError(perr)
				}
			}()
		}
		var addr string
		var resultErr error
		if innerMaster.ping {
			//api层调用Business
			bc, ok := obj.(businessCmd)
			if ok {
				addr = bc.BusinessAddr()
			} else {
				addr = innerMaster.randomRpcAddr()
			}
		} else {
			//Business层调用dao
			addr, resultErr = cmd.RemoteAddr(innerMaster)
			if resultErr != nil {
				return resultErr
			}
		}
		if addr == "" {
			return ErrNetwork
		}

		jrc, e := self.createJsonCall(addr)
		if e != nil {
			return ErrNetwork
		}

		callDone := jrc.client.Go(serviceMethod, cmd, reply, make(chan *rpc.Call, 1)).Done
		e = errors.New("超时请求:" + serviceMethod)
		for index := 1; index < 100; index++ {
			forbreak := false
			select {
			case call := <-callDone:
				forbreak = true
				e = call.Error
				break
			default:
				time.Sleep(time.Duration(index*5) * time.Millisecond)
			}
			if forbreak {
				break
			}
		}
		self.releaseJsonCall(jrc, e) //释放
		resultErr = e

		//如果是网络错误,sleep后重试
		if resultErr == nil {
			return resultErr
		}
		emsg := e.Error()
		if emsg == ErrShutdown.Error() || emsg == Unexpected.Error() || emsg == ErrConnect.Error() || strings.Contains(emsg, "closed network connection") || strings.Contains(emsg, "read: connection reset by peer") || strings.Contains(emsg, "write: broken pipe") || strings.Contains(emsg, "ServerShutDown") {
			return ErrNetwork
		}
		return resultErr
	}

	for index := 1; index <= self.retry; index++ {
		err = self.task.CallFunc(fun)
		if err == nil {
			return
		}

		if err != ErrNetwork {
			return err
		}
		time.Sleep(time.Duration(500+rand.Intn(1500)) * time.Millisecond) //休眠后重试, 重试次数越多,sleep时间越久
	}
	err = unknownNetwork
	return
}

func (self *RpcClient) createJsonCall(addr string) (client *rpcConn, err error) {
	if !self.connReuse {
		client, e := jsonRpc(addr)
		if e != nil {
			return nil, e
		}
		r := new(rpcConn)
		r.client = client
		return r, nil
	}
	//复用连接
	return connPool.takeConn(addr)
}

func (self *RpcClient) releaseJsonCall(client *rpcConn, err error) {
	if !self.connReuse {
		client.client.Close()
		return
	}
	//复用连接
	netok := true
	if err == nil {
		client.release(netok)
		return
	}
	emsg := err.Error()
	if emsg == ErrShutdown.Error() || emsg == ErrConnect.Error() || emsg == ErrNetwork.Error() {
		netok = false
	}
	client.release(netok)
	return
}

func (self *RpcClient) Close() {
	utils.Log().WriteInfo("jryg business close: 关闭task...")
	self.task.Close()
	utils.Log().WriteInfo("jryg business close: task已关闭")
	utils.Log().WriteInfo("jryg business close: 关闭dao连接池...")
	connPool.Close()
	utils.Log().WriteInfo("jryg business close: dao连接池已关闭")
	return
}

// RecoverLog 开启
func (self *RpcClient) RecoverLog() {
	self.recoverLog = true
	return
}

func (self *RpcClient) Qps(serviceMethod string, ms int64) {
	go func() {
		var rpcQps *RpcQps
		self.GetComponent(&rpcQps)
		if rpcQps == nil {
			panic("未获取到InnerQps")
		}
		rpcQps.Qps(serviceMethod, ms)
	}()
}

var ErrShutdown = errors.New("connection is shut down")
var ErrConnect = errors.New("dial is fail")
var ErrNetwork = errors.New("connection is shut down")
var Unexpected = errors.New("unexpected EOF")
var unknownNetwork = errors.New("未知网络错误")
