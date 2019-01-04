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
