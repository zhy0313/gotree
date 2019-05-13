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
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
	rpc "github.com/8treenet/gotree/lib/rpc"
)

func init() {
	startTime = time.Now().Format("2006-01-02 15:04:05")
}

type InnerServerController struct {
	RpcController
}

func (self *InnerServerController) Gotree() *InnerServerController {
	self.RpcController.Gotree(self)
	return self
}

//HandShake dao容器服务器发来的握手
func (self *InnerServerController) HandShake(cmd struct {
	DaoList []DaoNodeInfo
}, ret *int) error {
	//获取远程dao容器服务器ip
	ip := self.RemoteAddr()
	unix := time.Now().Unix()

	//添加dao节点
	for _, dni := range cmd.DaoList {
		node := NodeInfo{
			name:     dni.Name,
			lastUnxi: unix,
			ip:       ip,
			id:       dni.ID,
			port:     dni.Port,
			Extra:    dni.Extra,
		}
		//通知节点接入
		self.NotitySubscribe("HandShakeAddNode", node)
	}

	*ret = 666
	return nil
}

//api 服务器发来的握手
func (self *InnerServerController) Ping(arg interface{}, ret *int) error {
	*ret = 666
	return nil
}

//ProcessId 获取进程id
func (self *InnerServerController) ProcessId(arg interface{}, ret *int) error {
	*ret = os.Getpid()
	return nil
}

//api 服务器发来的握手
func (self *InnerServerController) DaoStatus(arg interface{}, ret *string) error {
	var list []*NodeInfo
	self.NotitySubscribe("InnerDaoInfo", &list)
	for _, item := range list {
		if *ret != "" {
			*ret += ";"
		}
		*ret += "dao:" + item.name + " id:" + fmt.Sprint(item.id) + " ip:" + item.ip + " port:" + item.port + " extra:" + fmt.Sprint(item.Extra)
	}
	return nil
}

//DaoServerInfo dao服务器信息
func (self *InnerServerController) DaoServerInfo(arg interface{}, ret *string) error {
	if *ret != "" {
		*ret += ";"
	}
	var list []*NodeInfo
	self.NotitySubscribe("InnerDaoInfo", &list)
	addrs := make(map[string]bool)
	for _, item := range list {
		addrs[item.ip+":"+item.port] = true
	}

	for addr, _ := range addrs {
		client, err := jsonRpc(addr)
		if err != nil {
			continue
		}
		var result string
		if client.Call("InnerServer.DaoInfo", "666", &result) != nil {
			client.Close()
			continue
		}
		client.Close()
		if *ret != "" {
			*ret += ";"
		}

		list := strings.Split(result, ",")
		if len(list) < 8 {
			continue
		}

		*ret += fmt.Sprintf("dao服务地址:%s, 使用内存:%smb, GCCPU占用:%s, GC次数:%s, 请求:%s, 数据库连接:%s, 日志队列:%s, 异步队列:%s, 启动时间:%s", addr, list[0], list[1], list[2], list[3], list[5], list[6], list[7], list[4])
	}
	return nil
}

//DaoInfo dao服务器信息
func (self *InnerServerController) DaoInfo(arg interface{}, ret *string) error {
	*ret = strings.Join(sysInfo(), ",")
	return nil
}

//BusinessInfo Business服务信息
func (self *InnerServerController) BusinessInfo(arg interface{}, ret *string) error {
	*ret = ""
	list := sysInfo()
	if len(list) < 7 {
		return nil
	}

	*ret = fmt.Sprintf("使用内存:%smb, GCCPU占用:%s, GC次数:%s, 请求:%d, 定时:%d, 异步:%d, 日志队列:%s, 启动时间:%s", list[0], list[1], list[2], rpc.CurrentCallNum(), lib.CurrentTimeNum(), AsynNumFunc(), list[6], startTime)
	return nil
}

//BusinessInfo Business服务信息
func (self *InnerServerController) DaoQps(arg interface{}, ret *string) error {
	var list []struct {
		ServiceMethod string
		Count         int64 //调用次数
		AvgMs         int64 //平均用时
		MaxMs         int64 //最高用时
		MinMs         int64 //最低用时
	}
	self.NotitySubscribe("DaoQps", &list)
	helper.SliceSort(&list, "AvgMs")
	*ret += fmt.Sprintf("%46s %12s %10s %10s %10s", "Call", "Count", "MaxMs", "MinMs", "AvgMs")
	for _, item := range list {
		if *ret != "" {
			*ret += "??"
		}

		callcount := fmt.Sprint(item.Count)
		if item.Count > 1000 {
			callcount = fmt.Sprintf("%.2fk", float32(item.Count)/1000.0)
		}

		*ret += fmt.Sprintf("%46s \x1b[0;31m%12s\x1b[0m \x1b[0;31m%10d\x1b[0m \x1b[0;31m%10d\x1b[0m \x1b[0;31m%10d\x1b[0m", item.ServiceMethod, callcount, item.MaxMs, item.MinMs, item.AvgMs)
	}
	return nil
}

var startTime string
var AsynNumFunc func() int

func sysInfo() (result []string) {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	result = append(result, fmt.Sprint(m.Alloc/1024/1024))        //内存占用mb
	result = append(result, fmt.Sprintf("%.3f", m.GCCPUFraction)) //gc cpu占用
	result = append(result, fmt.Sprint(m.NumGC))                  //gc 次数
	result = append(result, fmt.Sprint(rpc.CurrentCallNum()))     //当前 go 数量
	result = append(result, startTime)                            //系统启动时间
	if dbCountFunc != nil {
		result = append(result, fmt.Sprint(dbCountFunc())) //数据库总连接数
	} else {
		result = append(result, "") //数据库总连接数
	}
	result = append(result, fmt.Sprint(helper.Log().QueueLen())) //日志待处理
	if queueCountFunc != nil {
		result = append(result, fmt.Sprint(queueCountFunc())) //队列待处理
	} else {
		result = append(result, "0") //队列待处理
	}
	return
}

func SetDbCountFunc(fun func() int) {
	dbCountFunc = fun
	return
}

func SetQueueCountFunc(fun func() int) {
	queueCountFunc = fun
	return
}

var dbCountFunc func() int
var queueCountFunc func() int
