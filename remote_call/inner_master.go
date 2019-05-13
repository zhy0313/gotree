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
	"math/rand"
	"sync"
	"time"

	"github.com/8treenet/gotree/lib"
)

const (
	_CHECK_NODE_TICK     = 2000 //检测超时定时器触发 单位毫秒
	_CHECK_NODE_TIMEOUT  = 16   //检测node超时最长时间 单位秒
	_HANDSHAKE_NODE_TICK = 11   //node握手发送时间
)

type InnerMaster struct {
	lib.Object
	mutex   sync.Mutex
	nodeMap map[string]*nodeManage //节点地址[dao_name] = dao ip list
	addrMap map[string]int64       //普通地址
	ping    bool                   //true ping模式, false node模式
}

func (self *InnerMaster) Gotree() *InnerMaster {
	self.Object.Gotree(self)
	self.ping = false
	self.nodeMap = make(map[string]*nodeManage)
	self.addrMap = make(map[string]int64)
	lib.RunTickStopTimer(_CHECK_NODE_TICK, self.tick) //定时器检测超时节点
	rand.Seed(time.Now().Unix())

	self.AddSubscribe("HandShakeAddNode", self.addNode)
	self.AddSubscribe("InnerDaoInfo", self.innerDaoInfo)
	return self
}

//Tick 处理超时dao
func (self *InnerMaster) tick(stop *bool) {
	timeoutUnxi := time.Now().Unix() - _CHECK_NODE_TIMEOUT
	list := self.allNode()

	if self.ping {
		//如果是ping 模式
		*stop = true
		return
	}

	for _, item := range list {
		if item.lastUnxi == -1 {
			// -1 不处理
			continue
		}
		if item.lastUnxi < timeoutUnxi {
			self.removeNode(item)
		}
	}

	return
}

//GetAddrList 通过SerName的 获取该服务的远程地址列表
func (self *InnerMaster) GetAddrList(SerName string) *nodeManage {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	m, ok := self.nodeMap[SerName]
	if !ok {
		return nil
	}
	return m
}

func (self *InnerMaster) allNode() (list []*NodeInfo) {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	for _, m := range self.nodeMap {
		list = append(list, m.AllNode()...)
	}
	return
}

//AddNode 加入节点
func (self *InnerMaster) LocalAddNode(name, ip, port string, id int) {
	node := NodeInfo{
		name:     name,
		lastUnxi: -1,
		ip:       ip,
		port:     port,
		id:       id,
	}
	self.addNode(node)
}

//AddNode 加入节点
func (self *InnerMaster) addNode(args ...interface{}) {
	defer self.mutex.Unlock()
	self.mutex.Lock()

	node := args[0].(NodeInfo)
	_, ok := self.nodeMap[node.name]
	if !ok {
		//不存在该组件
		self.nodeMap[node.name] = new(nodeManage).nodeManage()
	}

	self.nodeMap[node.name].addNode(&node)
}

//RemoveNode 删除节点
func (self *InnerMaster) removeNode(node *NodeInfo) {
	defer self.mutex.Unlock()
	self.mutex.Lock()

	np, ok := self.nodeMap[node.name]
	if !ok {
		return
	}
	np.removeNode(node)
}

//AddAddr 加入地址
func (self *InnerMaster) addAddr(addr string) {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	self.addrMap[addr] = time.Now().Unix()
}

func (self *InnerMaster) removeAddr(addr string) {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	delete(self.addrMap, addr)
}

//RandomRpcAddr 随机节点地址 ping模式
func (self *InnerMaster) randomRpcAddr() string {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	addrList := []string{}
	timeoutUnxi := time.Now().Unix() - _CHECK_NODE_TIMEOUT
	for k, v := range self.addrMap {
		if v < timeoutUnxi {
			continue
		}
		addrList = append(addrList, k)
	}

	if len(addrList) == 0 {
		return ""
	}

	index := rand.Intn(len(addrList))
	return addrList[index]
}

//获取远程地址
func (self *InnerMaster) addrList() []string {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	addrList := []string{}
	for k, _ := range self.addrMap {
		addrList = append(addrList, k)
	}
	return addrList
}

//InnerDaoInfo 查看节点
func (self *InnerMaster) innerDaoInfo(args ...interface{}) {
	ret := args[0].(*[]*NodeInfo)
	list := self.allNode()
	*ret = list
}
