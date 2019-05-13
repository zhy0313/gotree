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

package chart

import (
	"fmt"
	"hash/crc32"
	"sort"
)

//一致性哈希分发
type DummyNode struct {
	dummyList []dummy //虚拟列表
}

//创建虚拟节点
func (self *DummyNode) DummyNode() *DummyNode {
	self.dummyList = make([]dummy, 0, 2048)
	return self
}

//Add加入地址
func (self *DummyNode) Add(addr interface{}) {
	self.Remove(addr)
	pa := new(physicalAddr)
	pa.addr = addr

	dlist := pa.createDummy()
	self.dummyList = append(self.dummyList, dlist...)
	sort.Sort(dummys(self.dummyList))
}

//Remove移除地址
func (self *DummyNode) Remove(addr interface{}) {
	newList := make([]dummy, 0, 2048) //全部节点
	for index := 0; index < len(self.dummyList); index++ {
		if self.dummyList[index].pAddr.addr == addr {
			self.dummyList[index].pAddr = nil
			continue
		}

		newList = append(newList, self.dummyList[index])
	}
	self.dummyList = newList
}

//CheckAddr检查该物理地址是否存在
func (self *DummyNode) Check(addr interface{}) bool {
	for index := 0; index < len(self.dummyList); index++ {
		if self.dummyList[index].pAddr.addr == addr {
			return true
		}
	}
	return false
}

//Search查找节点位置
func (self *DummyNode) Search(value uint32) interface{} {
	if len(self.dummyList) == 0 {
		return nil
	}
	head := self.dummyList[0]
	tail := self.dummyList[len(self.dummyList)-1]
	if value < head.value {
		return head.pAddr.addr
	}

	if value > tail.value {
		return head.pAddr.addr
	}

	return self.half(value, 0, len(self.dummyList)-1)
}

//half 2分查找节点位置
func (self *DummyNode) half(value uint32, begin int, end int) interface{} {

	if begin == end || (end-begin) == 1 {
		if self.dummyList[end].value == value {
			return self.dummyList[end].pAddr.addr
		}
		return self.dummyList[begin].pAddr.addr
	}

	halfIndex := begin + (end-begin)/2
	halfDummy := self.dummyList[halfIndex]
	if value > halfDummy.value {
		return self.half(value, halfIndex, end)
	}
	return self.half(value, begin, halfIndex)
}

//物理地址
type physicalAddr struct {
	addr interface{} //地址
}

//虚拟节点
type dummy struct {
	pAddr *physicalAddr //对应的物理地址
	value uint32        //虚拟节点值
}

type dummys []dummy

func (self dummys) Len() int {
	return len(self)
}

func (self dummys) Less(i, j int) bool {
	return self[i].value < self[j].value
}

func (self dummys) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self *physicalAddr) createDummy() []dummy {
	result := make([]dummy, 512)
	//为物理地址创建512虚拟节点
	for index := 0; index < 512; index++ {
		result[index].pAddr = self
		result[index].value = crc32.ChecksumIEEE([]byte(fmt.Sprint(self.addr) + fmt.Sprint(index+1)))
	}
	return result
}
