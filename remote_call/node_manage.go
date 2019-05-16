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
	"hash/crc32"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/8treenet/gotree/lib/chart"
)

//节点管理器
type nodeManage struct {
	nodes    map[string]*NodeInfo //[ip]组件
	seq      int
	mutex    sync.Mutex
	hotNode  *chart.DummyNode //热一致性哈希节点方式
	hashNode *chart.DummyNode //一致性哈希节点方式
}

func (self *nodeManage) nodeManage() *nodeManage {
	self.nodes = make(map[string]*NodeInfo)
	self.seq = 1000
	self.hotNode = new(chart.DummyNode).DummyNode()
	self.hashNode = new(chart.DummyNode).DummyNode()
	rand.Seed(time.Now().Unix())
	return self
}

//RandomAddr 随机节点地址
func (self *nodeManage) RandomAddr() string {
	defer self.mutex.Unlock()
	self.mutex.Lock()

	list := self.allNode()
	if len(list) == 0 {
		return ""
	}

	index := rand.Intn(len(list))
	return list[index].RpcAddr()
}

//BalanceAddr 均衡节点地址
func (self *nodeManage) BalanceAddr() string {
	defer self.mutex.Unlock()
	self.mutex.Lock()

	list := self.allNode()
	if len(list) == 0 {
		return ""
	}

	seq := self.seq
	if seq > 9999999999 {
		self.seq = 1000
	} else {
		self.seq++
	}
	return list[seq%len(list)].RpcAddr()
}

//HostHashRpcAddr  热一致性哈希节点地址
func (self *nodeManage) HostHashRpcAddr(value interface{}) string {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	searchValue := crc32.ChecksumIEEE([]byte(fmt.Sprint(value)))
	v := self.hotNode.Search(searchValue)
	if v == nil {
		return ""
	}
	nodeid := v.(int)

	for _, node := range self.allNode() {
		if node.id == nodeid {
			return node.RpcAddr()
		}
	}
	return ""
}

//HotHashRpcId 热一致性哈希节点id
func (self *nodeManage) HotHashRpcId(value interface{}) int {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	searchValue := crc32.ChecksumIEEE([]byte(fmt.Sprint(value)))
	v := self.hotNode.Search(searchValue)
	if v == nil {
		return 0
	}
	nodeid := v.(int)
	return nodeid
}

//HashRpcAddr 一致性哈希节点地址
func (self *nodeManage) HashRpcAddr(value interface{}) string {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	searchValue := crc32.ChecksumIEEE([]byte(fmt.Sprint(value)))
	v := self.hashNode.Search(searchValue)
	if v == nil {
		return ""
	}
	nodeid := v.(int)

	for _, node := range self.allNode() {
		if node.id == nodeid {
			return node.RpcAddr()
		}
	}
	return ""
}

//HashRpcId 一致性哈希节点id
func (self *nodeManage) HashRpcId(value interface{}) int {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	searchValue := crc32.ChecksumIEEE([]byte(fmt.Sprint(value)))
	v := self.hashNode.Search(searchValue)
	if v == nil {
		return 0
	}
	nodeid := v.(int)
	return nodeid
}

//SlaveAddr 从节点地址
func (self *nodeManage) SlaveAddr() string {
	defer self.mutex.Unlock()
	self.mutex.Lock()

	list := self.allNode()
	if len(list) == 0 {
		return ""
	}

	if len(list) == 1 {
		return list[0].RpcAddr()
	}

	slvae := []*NodeInfo{}
	for _, ni := range list {
		if ni.id != 1 {
			//排除主节点
			slvae = append(slvae, ni)
		}
	}

	index := rand.Intn(len(slvae))
	return slvae[index].RpcAddr()
}

//MasterAddr 返回主节点
func (self *nodeManage) MasterAddr() string {
	defer self.mutex.Unlock()
	self.mutex.Lock()

	list := self.allNode()
	if len(list) == 0 {
		return ""
	}

	for _, node := range list {
		if node.id == 1 {
			return node.RpcAddr()
		}
	}
	return ""
}

//removeNode 删除节点
func (self *nodeManage) removeNode(node *NodeInfo) {
	defer self.mutex.Unlock()
	self.mutex.Lock()

	locNode, ok := self.nodes[node.ip+"_"+node.port]
	if !ok {
		return
	}

	//删除一致性虚拟节点
	self.hotNode.Remove(locNode.id)
	delete(self.nodes, node.ip+"_"+node.port)
}

//addNode 加入节点
func (self *nodeManage) addNode(node *NodeInfo) {
	addr := node.ip + "_" + node.port
	delList := self.allNode()
	for _, item := range delList {
		if item.id == node.id && addr != item.ip+"_"+item.port {
			self.removeNode(item)
		}
	}
	defer self.mutex.Unlock()
	self.mutex.Lock()
	locNode, ok := self.nodes[addr]
	if ok {
		//存在该name ip 的dao 更新访问时间
		locNode.lastUnxi = node.lastUnxi
		return
	}
	self.nodes[addr] = node

	//加入一致性虚拟节点和普通节点
	self.hotNode.Add(node.id)
	maxkv := ""
	for _, extra := range node.Extra {
		str := extra.(string)
		if !strings.Contains(str, "MaxID") {
			continue
		}
		maxkv = str
	}
	if maxkv == "" {
		return
	}

	if kv := strings.Split(maxkv, ":"); len(kv) == 2 {
		maxid, err := strconv.Atoi(kv[1])
		if err != nil {
			return
		}
		for index := 1; index <= maxid; index++ {
			if self.hashNode.Check(index) {
				continue
			}
			self.hashNode.Add(index)
		}
	}
}

//all 全部节点
func (self *nodeManage) allNode() (list []*NodeInfo) {
	for _, m := range self.nodes {
		list = append(list, m)
	}
	return
}

//all 全部节点
func (self *nodeManage) AllNode() (list []*NodeInfo) {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	return self.allNode()
}

//Len
func (self *nodeManage) Len() int {
	defer self.mutex.Unlock()
	self.mutex.Lock()
	return len(self.nodes)
}

//节点信息
type NodeInfo struct {
	lastUnxi int64  //最后请求时间时间戳
	name     string //节点名字
	id       int    //节点id 从1开始
	ip       string
	port     string
	Extra    []interface{}
}

//RpcAddr 节点地址
func (self *NodeInfo) RpcAddr() string {
	return self.ip + ":" + self.port
}

//dao
type DaoNodeInfo struct {
	Name  string        //节点名称
	Port  string        //节点端口
	ID    int           //节点id, 假如该节点存在多个
	Extra []interface{} `opt:"empty"` //扩展信息, 如lbs 等经纬度 区域划分
}
