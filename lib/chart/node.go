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
)

//哈希分发
type Node struct {
	mod   uint32 //取模
	value interface{}
}

func (self *Node) Node(mod uint32) *Node {
	self.mod = mod
	return self
}

func (self *Node) HashNodeSum(v interface{}, modarg ...int) int {
	value := fmt.Sprint(v)
	mod := self.mod
	if len(modarg) > 0 {
		mod = uint32(modarg[0])
	}
	result := crc32.ChecksumIEEE([]byte(value)) % mod
	return int(result) + 1
}
