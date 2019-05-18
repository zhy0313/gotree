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

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
	rpc "github.com/8treenet/gotree/lib/rpc"
	jsonrpc "github.com/8treenet/gotree/lib/rpc/jsonrpc"
)

var rs *rpcServer

func init() {
	rs = new(rpcServer).Gotree()
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
			helper.Log().WriteError(err.Error())
			panic(err.Error())
		}
		if err == nil {
			port += index
			helper.Log().WriteInfo("rpc server bind addr:" + fmt.Sprintf("%s:%d", ip, port))
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
		helper.Log().WriteError("RpcServerRegister:registration failed:" + reflect.TypeOf(controller).String())
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

func (self *rpcServer) Gotree() *rpcServer {
	self.Object.Gotree(self)
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
		helper.Log().WriteError(err.Error())
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
