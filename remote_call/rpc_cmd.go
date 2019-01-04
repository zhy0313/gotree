package remote_call

import (
	"net/http"
	"strings"

	"jryghq.cn/lib"
	"jryghq.cn/lib/rpc"
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

type RpcHeadCmd struct {
	RpcCmd
	Head string `opt:"empty"`
}

func (self *RpcHeadCmd) RpcHeadCmd(child interface{}) *RpcHeadCmd {
	self.Object.Object(self)
	self.AddChild(self, child)
	return self
}

func (self *RpcHeadCmd) Header(k string) string {
	return _header.Get(self.Head, k)
}

func (self *RpcHeadCmd) SetHeader(k string, v string) {
	_, e := self.GetChild(self)
	if e != nil {
		panic("RpcHeadCmd :header是只读数据")
	}
	self.Head = _header.Set(self.Head, k, v)
}

func (self *RpcHeadCmd) SetHttpHeader(head http.Header) {
	_, e := self.GetChild(self)
	if e != nil {
		panic("RpcHeadCmd :header是只读数据")
	}
	for item := range head {
		self.SetHeader(item, head.Get(item))
	}
}

type RpcCmd struct {
	lib.Object
	Bseq          string `opt:"empty"`
	cacheIdentity interface{}
}

type RpcNode interface {
	RandomRpcAddr() string                    //随机地址
	BalanceRpcAddr() string                   //负载均衡地址
	HostHashRpcAddr(value interface{}) string //热一致性哈希地址
	HashRpcAddr(value interface{}) string     //一致性哈希地址
	SlaveRpcAddr() string                     //返回随机从节点  主节点:节点id=1,当只有主节点返回主节点
	MasterRpcAddr() string                    //返回主节点
	AllNode() (list []*NodeInfo)              //获取全部节点,自定义分发
}

type cmdChild interface {
	Control() string        //远程服务名称
	Action() string         //远程服务方法
	DaoAddr(RpcNode) string //相关服务方法远程地址回调
}

type cmdSerChild interface {
	Control() string //远程服务名称
	Action() string  //远程服务方法
}

type nodeAddr interface {
	GetAddrList(string) *nodeManage //传入远程服务名称, 获取服务Addr列表
}

func (self *RpcCmd) RpcCmd(child interface{}) *RpcCmd {
	self.Object.Object(self)
	self.AddChild(self, child)
	bseq := rpc.GoDict().Get("bseq")
	if bseq == nil {
		return self
	}
	str, ok := bseq.(string)
	if ok {
		self.Bseq = str
	}
	self.cacheIdentity = nil
	return self
}

//CacheOn 开启命令缓存 identification标识 通常填写 id
func (self *RpcCmd) CacheOn(identity interface{}) {
	if self.Bseq == "" {
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
		panic("顶级子类未实现 interface :" + className)
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
	return childObj.DaoAddr(nm), nil
}

//获取ServiceMethod
func (self *RpcCmd) ServiceMethod() string {
	child := self.TopChild()
	childObj, ok := child.(cmdSerChild)
	if !ok {
		className := self.ClassName(child)
		panic("顶级子类未实现 interface :" + className)
	}

	return childObj.Control() + "." + childObj.Action()
}

var _header RpcHeader
