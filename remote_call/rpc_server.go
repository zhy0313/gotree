package remote_call

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"

	"jryghq.cn/lib"
	rpc "jryghq.cn/lib/rpc"
	jsonrpc "jryghq.cn/lib/rpc/jsonrpc"
	"jryghq.cn/utils"
)

var rs *rpcServer

func init() {
	rs = new(rpcServer).rpcServer()
}

//启动rpc服务 args[0] = string(x.x.x.x:2321)
func RpcServerRun(args ...interface{}) *rpcServer {
	bindAddr := "0.0.0.0:8888" //默认
	if len(args) > 0 {
		bindAddr = args[0].(string)
	}

	if len(args) > 1 {
		sercall := args[1].(func(svrName string))
		list := rs.SvrNames()
		for _, name := range list {
			sercall(name)
		}
	}

	ip := strings.Split(bindAddr, ":")[0]
	port, _ := strconv.Atoi(strings.Split(bindAddr, ":")[1])
	for index := 0; index < 10; index++ {
		var err error
		rs.socket, err = net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port+index))
		if err != nil && index == 9 {
			utils.Log().WriteError(err.Error())
			panic(err.Error())
		}
		if err == nil {
			port += index
			utils.Log().WriteInfo("rpc server bind addr:" + fmt.Sprintf("%s:%d", ip, port))
			break
		}
	}

	go rs.run()
	rs.NotitySubscribe("HookRpcBind", strconv.Itoa(port))

	//监听kill pid 和 control+c
	channel := make(chan os.Signal)
	signal.Notify(channel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV, syscall.SIGFPE, syscall.SIGABRT, syscall.SIGILL)
	<-channel
	atomic.AddInt32(&rpccloseChan, 1)
	return rs
}

//RpcServerRegister 注册rpc服务
func RpcServerRegister(controller interface{}) {
	type rpcname interface {
		RpcName() string
	}

	rc, ok := controller.(rpcname)
	if !ok {
		utils.Log().WriteError("RpcServerRegister 注册失败:" + reflect.TypeOf(controller).String())
	}
	rs.register(rc.RpcName(), controller)
}

//RpcUnServerRegister 取消注册rpc服务
func RpcUnServerRegister(name string) {
	rs.unregister(name)
}

//RpcControllerNames 返回controller名字
func RpcControllerNames() []string {
	return rs.SvrNames()
}

var rpccloseChan int32

type rpcServer struct {
	socket net.Listener
	srv    *rpc.Server
	lib.Object
}

func (self *rpcServer) rpcServer() *rpcServer {
	self.Object.Object(self)
	self.srv = rpc.NewServer()
	rpccloseChan = 0
	return self
}

func (self *rpcServer) unregister(name string) {
	if name == "InnerServer" {
		return //内部服务
	}
	self.srv.UnRegister(name)
}

//register 注册相关rpc服务
func (self *rpcServer) register(name string, controller interface{}) {
	err := self.srv.RegisterName(name, controller)
	if err != nil {
		utils.Log().WriteError(err.Error())
	}
}

//svrNames 已注册服务名字
func (self *rpcServer) SvrNames() []string {
	return self.srv.ServiceNameList()
}

//run 启动
func (self *rpcServer) run() {
	for {
		conn, err := self.socket.Accept()
		if err != nil {
			continue
		}
		if atomic.LoadInt32(&rpccloseChan) == 2 {
			conn.Close()
			continue
		}
		go self.srv.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

//Close 关闭
func (self *rpcServer) Close() {
	atomic.AddInt32(&rpccloseChan, 1)
	if self.socket == nil {
		return
	}
	self.socket.Close()
	self.srv.Close()
}
