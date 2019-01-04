package remote_call

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
)

type RpcController struct {
	lib.Object
	conn net.Conn
}

//RpcServer 构造
func (self *RpcController) Gotree(child interface{}) *RpcController {
	self.Object.Gotree(self)
	self.Object.AddChild(self, child)
	return self
}

//Prepare 底层rpc系统反射实例化后(构造函数后)调用,init() new实例化里,不调用. method是方法, argv是远程调用参数,
//可重写记录日志等
func (self *RpcController) Prepare(method string, argv interface{}) {

}

//Finish 结束后调用 reply是要返回的参数
func (self *RpcController) Finish(method string, reply interface{}, e error) {

}

//RemoteAddr 获取调用者ip
func (self *RpcController) RemoteAddr() string {
	list := strings.Split(self.conn.RemoteAddr().String(), ":")
	if len(list) > 0 {
		return list[0]
	}
	return ""
}

//RemotePort 获取调用者port
func (self *RpcController) RemotePort() int {
	list := strings.Split(self.conn.RemoteAddr().String(), ":")
	if len(list) > 1 {
		port, _ := strconv.Atoi(list[1])
		return port
	}
	return 0
}

//以下内部处理
func (self *RpcController) TestRpcControllerTest__(cmd int, result *int) error {
	return nil
}

//RpcInvoke rpc service 反射 远程连接ip
func (self *RpcController) RpcInvoke(conn net.Conn) bool {
	close := atomic.LoadInt32(&rpccloseChan)
	if close == 1 {
		return false
	}
	if close == 2 {
		conn.Close()
		return false
	}
	self.conn = conn
	return true
}

//用于rpcserver 注册,不可重写
func (self *RpcController) RpcName() string {
	//获取顶级子类
	child := self.TopChild()
	controller := reflect.TypeOf(child).Elem().String()
	name := self.controllerName(controller)
	if name == "" {
		helper.Log().WriteError("rpc controller 不符合规定:" + controller)
	}
	return name
}

//controllerName 切割名字
func (self *RpcController) controllerName(name string) string {
	list := strings.Split(name, ".")
	if len(list) < 2 {
		return ""
	}

	list = strings.Split(list[len(list)-1], "Controller")
	if len(list) < 1 {
		return ""
	}
	return list[0]
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	}
	return fmt.Sprintf("%v", src)
}
