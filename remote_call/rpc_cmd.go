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
	"net/http"
	"strings"

	"github.com/8treenet/gotree/lib"
	"github.com/8treenet/gotree/lib/rpc"
)

type RpcHeader struct {
}

func (self *RpcHeader) Set(src, key, val string) string {
	if src != "" {
		src += "??"
	}
	src += key + "::" + val
	return src
}

func (self *RpcHeader) Get(src, key string) string {
	list := strings.Split(src, "??")
	for _, kvitem := range list {
		kv := strings.Split(kvitem, "::")
		if len(kv) != 2 {
			continue
		}
		if kv[0] == key {
			return kv[1]
		}
	}
	return ""
}

type RpcCmd struct {
	lib.Object
	Gseq          string `opt:"empty"`
	Head          string `opt:"empty"`
	cacheIdentity interface{}
}

func (self *RpcCmd) Header(k string) string {
	return _header.Get(self.Head, k)
}

func (self *RpcCmd) SetHeader(k string, v string) {
	_, e := self.GetChild(self)
	if e != nil {
		panic("RpcCmd-SetHeader :Header is read-only data")
	}
	self.Head = _header.Set(self.Head, k, v)
}

func (self *RpcCmd) SetHttpHeader(head http.Header) {
	_, e := self.GetChild(self)
	if e != nil {
		panic("RpcCmd-SetHttpHeader :Header is read-only data")
	}
	for item := range head {
		self.SetHeader(item, head.Get(item))
	}
}

type ComNode interface {
	RandomAddr() string                    //随机地址
	BalanceAddr() string                   //负载均衡地址
	HostHashRpcAddr(value interface{}) string //热一致性哈希地址
	HashRpcAddr(value interface{}) string     //一致性哈希地址
	SlaveAddr() string                     //返回随机从节点  主节点:节点id=1,当只有主节点返回主节点
	MasterAddr() string                    //返回主节点
	AllNode() (list []*NodeInfo)              //获取全部节点,自定义分发
}

type cmdChild interface {
	Control() string        //远程服务名称
	Action() string         //远程服务方法
	ComAddr(ComNode) string //相关服务方法远程地址回调
}

type cmdSerChild interface {
	Control() string //远程服务名称
	Action() string  //远程服务方法
}

type nodeAddr interface {
	GetAddrList(string) *nodeManage //传入远程服务名称, 获取服务Addr列表
}

func (self *RpcCmd) Gotree(child interface{}) *RpcCmd {
	self.Object.Gotree(self)
	self.AddChild(self, child)
	gseq := rpc.GoDict().Get("gseq")
	if gseq == nil {
		return self
	}
	str, ok := gseq.(string)
	if ok {
		self.Gseq = str
	}
	self.cacheIdentity = nil
	return self
}

//CacheOn 开启命令缓存 identification标识 通常填写 id
func (self *RpcCmd) CacheOn(identity interface{}) {
	if self.Gseq == "" {
		//非请求会话,不可开启缓存
		return
	}
	self.cacheIdentity = identity
	return
}

func (self *RpcCmd) Cache() interface{} {
	return self.cacheIdentity
}

//client调用RemoteAddr 传入NodeMaster, 获取远程地址
func (self *RpcCmd) RemoteAddr(naddr interface{}) (string, error) {
	child := self.TopChild()
	childObj, ok := child.(cmdChild)
	if !ok {
		className := self.ClassName(child)
func (self *RpcCmd) RemoteAddr(naddr interface{}) (string, error) {
		panic("RpcCmd-RemoteAddr:Subclass is not implemented,interface :" + className)
	}

	//获取子类要调用的服务名
	serName := childObj.Control()
	//通过服务名调用NodeMaster 获取该服务远程地址列表
	node := naddr.(nodeAddr)
	nm := node.GetAddrList(serName)
	if nm == nil || nm.Len() == 0 {
		return "", ErrNetwork
	}

	//传入子类该服务远程地址列表,计算要使用的节点
	return childObj.ComAddr(nm), nil
}

//获取ServiceMethod
func (self *RpcCmd) ServiceMethod() string {
	child := self.TopChild()
	childObj, ok := child.(cmdSerChild)
	if !ok {
		className := self.ClassName(child)
		panic("RpcCmd-ServiceMethod:Subclass is not implemented, interface :" + className)
	}

	return childObj.Control() + "." + childObj.Action()
}

var _header RpcHeader
