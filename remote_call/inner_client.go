package remote_call

import (
	"fmt"
	"net"
	netrpc "net/rpc"
	"net/rpc/jsonrpc"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/8treenet/gotree/lib"
)

type InnerClient struct {
	lib.Object
	ping          bool
	BusinesssAddr []*struct {
		Addr      string
		Port      int
		StartPort int
	}
	daos         []DaoNodeInfo
	StartTime    string
	LastInfoTime int64
	dbCountFunc  func() int
	closeMsg     chan bool
}

func (self *InnerClient) Gotree() *InnerClient {
	self.Object.Gotree(self)
	self.ping = false
	self.BusinesssAddr = make([]*struct {
		Addr      string
		Port      int
		StartPort int
	}, 0)
	self.daos = make([]DaoNodeInfo, 0, 32)
	self.AddSubscribe("HookRpcBind", self.hookRpcBind)
	self.StartTime = time.Now().Format("2006-01-02 15:04:05")
	self.LastInfoTime = 0
	self.closeMsg = make(chan bool, 1)
	return self
}

//AddRemoteSerAddr 加入远程地址
func (self *InnerClient) AddRemoteAddr(RemoteAddr string) {
	var innerMaster *InnerMaster
	self.GetComponent(&innerMaster)
	innerMaster.ping = true
	self.ping = true

	innerMaster.addAddr(RemoteAddr)
	return
}

//AddRemoteAddrByNode 加入远程地址节点方式
func (self *InnerClient) AddBusiness(ip string, port int) {
	var item struct {
		Addr      string
		Port      int
		StartPort int
	}
	item.Addr = ip
	item.Port = port
	item.StartPort = port
	self.BusinesssAddr = append(self.BusinesssAddr, &item)
}

//AddRemoteSerAddrByNode 加入远程地址节点方式 如果使用主从方式分布dao, 固定 id 1 为主
func (self *InnerClient) AddDaoByNode(name string, id int, args ...interface{}) {
	dao := DaoNodeInfo{
		Name:  name,
		Port:  "",
		ID:    id,
		Extra: args,
	}
	self.daos = append(self.daos, dao)
	return
}

//hookRpcBind hook rpc server启动
func (self *InnerClient) hookRpcBind(port ...interface{}) {
	localPort := fmt.Sprint(port[0])
	tempdaos := []DaoNodeInfo{}
	for _, v := range self.daos {
		dao := DaoNodeInfo{
			Name:  v.Name,
			Port:  localPort,
			ID:    v.ID,
			Extra: v.Extra,
		}
		tempdaos = append(tempdaos, dao)
	}
	self.daos = tempdaos
	go self.Run()
	return
}

//Run
func (self *InnerClient) Run() {
	var innerMaster *InnerMaster
	var reply int
	self.GetComponent(&innerMaster)
	clientMap := make(map[string]*netrpc.Client)
	timeIndex := 2
	if !self.ping {
		timeIndex = 4
	} else {
		list := innerMaster.addrList()
		for _, addr := range list {
			addrSplit := strings.Split(addr, ":")
			i, err := strconv.Atoi(addrSplit[1])
			if err != nil {
				panic(err.Error())
			}
			self.AddBusiness(addrSplit[0], i)
		}
	}
	for {
		if self.ping {
			for _, v := range self.BusinesssAddr {
				rpcaddr := fmt.Sprintf("%s:%d", v.Addr, v.Port)
				client, ok := clientMap[rpcaddr]
				if !ok {
					var cerr error
					client, cerr = jsonRpc(rpcaddr, 600)
					if cerr == nil {
						clientMap[rpcaddr] = client
					}
				}
				if client != nil {
					if client.Call("InnerServer.Ping", 100, &reply) == nil {
						innerMaster.addAddr(fmt.Sprintf("%s:%d", v.Addr, v.Port))
						//连接成功并且呼叫成功
						continue
					}
					client.Close()
					delete(clientMap, rpcaddr)
				}
				innerMaster.removeAddr(fmt.Sprintf("%s:%d", v.Addr, v.Port))

				for index := 0; index < 9; index++ {
					rpcaddr = fmt.Sprintf("%s:%d", v.Addr, v.StartPort+index)
					client, err := jsonRpc(fmt.Sprintf("%s:%d", v.Addr, v.StartPort+index), 600)
					if err == nil {
						if client.Call("InnerServer.Ping", 100, &reply) == nil {
							v.Port = v.StartPort + index
							innerMaster.addAddr(fmt.Sprintf("%s:%d", v.Addr, v.StartPort+index))
							clientMap[rpcaddr] = client
							break
						}
						client.Close()
					}
				}
			}
		} else {

			var cmd struct {
				DaoList []DaoNodeInfo
			}
			cmd.DaoList = self.daos

			for _, v := range self.BusinesssAddr {
				rpcaddr := fmt.Sprintf("%s:%d", v.Addr, v.Port)
				client, ok := clientMap[rpcaddr]
				if !ok {
					var cerr error
					client, cerr := jsonRpc(rpcaddr, 600)
					if cerr == nil {
						clientMap[rpcaddr] = client
					}
				}

				if client != nil {
					if client.Call("InnerServer.HandShake", cmd, &reply) == nil {
						//连接成功并且呼叫成功
						continue
					}
					client.Close()
					delete(clientMap, rpcaddr)
				}

				for index := 0; index < 9; index++ {
					client, err := jsonRpc(fmt.Sprintf("%s:%d", v.Addr, v.StartPort+index), 600)
					if err == nil {
						if client.Call("InnerServer.HandShake", cmd, &reply) == nil {
							v.Port = v.StartPort + index
							clientMap[rpcaddr] = client
							break
						}
						client.Close()
					}
				}
			}
		}

		for index := 0; index < timeIndex; index++ {
			select {
			case _ = <-self.closeMsg:
				runtime.Goexit()
			default:
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (self *InnerClient) Close() {
	self.closeMsg <- true
}

func jsonRpc(addr string, timeout ...int) (client *netrpc.Client, e error) {
	t := 1200
	if len(timeout) > 0 && timeout[0] > 500 {
		t = timeout[0]
	}
	tcpclient, err := net.DialTimeout("tcp", addr, time.Duration(t)*time.Millisecond)
	if err != nil {
		e = err
		return
	}
	client = jsonrpc.NewClient(tcpclient)
	return
}
