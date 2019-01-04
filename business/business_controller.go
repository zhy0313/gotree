package business

import (
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"jryghq.cn/lib/rpc"
	"jryghq.cn/remote_call"
	"jryghq.cn/utils"
)

var bseq int64
var identification string

func init() {
	bseq = 1
	rand.Seed(time.Now().Unix())
	x := int64(rand.Intn(10000))
	identification = strconv.FormatInt(x, 36)
}

//BusinessController
type BusinessController struct {
	remote_call.RpcController
}

//BusinessController
func (self *BusinessController) BusinessController(child interface{}) *BusinessController {
	self.RpcController.RpcController(self)
	self.AddChild(self, child)
	rpc.GoDict().Set("bseq", getBseq())
	return self
}

//服务定位器获取服务
func (self *BusinessController) Service(child interface{}) {
	err := _scl.Service(child)
	if err != nil {
		utils.Log().WriteError("飞哥:不要乱调用:" + err.Error())
		panic("飞哥:不要乱调用:" + err.Error())
	}
	return
}

// ServiceBroadcast 调用所有注册service的method方法. 潜龙勿用,会使项目非常难以维护
func (self *BusinessController) ServiceBroadcast(method string, arg interface{}) {
	if e := _scl.Broadcast(method, arg); e != nil {
		utils.Log().WriteError("Buesiness ServiceBroadcast errror:" + e.Error())
	}
}

//Async 异步执行
func (self *BusinessController) Async(run func(ac AsyncController), completeds ...func()) {
	var completed func()
	if len(completeds) > 0 {
		completed = completeds[0]
	}
	ac := new(async).async(run, completed)
	go ac.execute()
}

var bseqMutex sync.Mutex

func getBseq() (result string) {
	defer bseqMutex.Unlock()
	bseqMutex.Lock()
	result = identification
	result += strconv.FormatInt(bseq, 36)
	if bseq == math.MaxInt64 {
		bseq = 1
		return
	}
	bseq += 1
	return
}
