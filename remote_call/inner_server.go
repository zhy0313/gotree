package remote_call

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"jryghq.cn/utils"

	"jryghq.cn/lib"
	rpc "jryghq.cn/lib/rpc"
)

func init() {
	serverInfo = make(map[string][]string)
	serverInfoTime = make(map[string]int64)
	startTime = time.Now().Format("2006-01-02 15:04:05")
}

type InnerServerController struct {
	RpcController
}

func (self *InnerServerController) InnerServerController() *InnerServerController {
	self.RpcController.RpcController(self)
	return self
}

//HandShake dao容器服务器发来的握手
func (self *InnerServerController) HandShake(cmd struct {
	DaoList []DaoNodeInfo
	Info    []string `opt:"empty"` //程序信息
}, ret *int) error {
	//获取远程dao容器服务器ip
	ip := self.RemoteAddr()
	unix := time.Now().Unix()
	var list []*NodeInfo
	self.NotitySubscribe("InnerDaoInfo", &list)

	//添加dao节点
	for _, dni := range cmd.DaoList {
		for _, listItem := range list {
			if listItem.name == dni.Name && listItem.id == dni.ID && (listItem.ip != ip || listItem.port != dni.Port) {
				emsg := fmt.Sprintf("business内部连接警告, 已存在dao:(名字:%s,id:%d,ip:%s,端口:%s)", dni.Name, dni.ID, listItem.ip, listItem.port)
				utils.Log().WriteWarn(emsg)
			}
		}
		node := NodeInfo{
			name:     dni.Name,
			lastUnxi: unix,
			ip:       ip,
			id:       dni.ID,
			port:     dni.Port,
			Extra:    dni.Extra,
		}
		if len(cmd.Info) > 0 {
			addServerInfo(node.ip+":"+node.port, cmd.Info)
		}
		//通知节点接入
		self.NotitySubscribe("HandShakeAddNode", node)
	}

	*ret = 666
	return nil
}

func (self *InnerServerController) DaoOff(cmd struct {
	DaoList []DaoNodeInfo
	Info    []string `opt:"empty"` //程序信息
}, ret *int) error {
	ip := self.RemoteAddr()
	unix := time.Now().Unix() - 15
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
		if len(cmd.Info) > 0 {
			addServerInfo(node.ip+":"+node.port, cmd.Info)
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
	list := getServerInfo()
	for _, item := range list {
		if *ret != "" {
			*ret += ";"
		}
		str := fmt.Sprintf("dao服务地址:%s, 使用内存:%smb, GCCPU占用:%s, GC次数:%s, 并发:%s, 数据库连接:%s, 启动时间:%s", item.Addr, item.List[0], item.List[1], item.List[2], item.List[3], item.List[4], item.List[5])
		*ret += str
	}
	return nil
}

//BusinessInfo Business服务信息
func (self *InnerServerController) BusinessInfo(arg interface{}, ret *string) error {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	num := rpc.CurrentCallNum() + lib.CurrentTimeNum() + AsynNumFunc()
	*ret = fmt.Sprintf("使用内存:%smb, GCCPU占用:%s, GC次数:%s, 并发:%d, 启动时间:%s", fmt.Sprint(m.Alloc/1024/1024), fmt.Sprintf("%.3f", m.GCCPUFraction), fmt.Sprint(m.NumGC), num, startTime)
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
	utils.ArraySort(&list, "AvgMs")
	for _, item := range list {
		if *ret != "" {
			*ret += "??"
		}

		callcount := fmt.Sprint(item.Count)
		if item.Count > 1000 {
			callcount = fmt.Sprintf("%.2fk", float32(item.Count)/1000.0)
		}

		*ret += fmt.Sprintf("%46s \x1b[0;31m%12s\x1b[0m \x1b[0;31m%10d\x1b[0mms \x1b[0;31m%10d\x1b[0mms \x1b[0;31m%10d\x1b[0mms", item.ServiceMethod, callcount, item.MaxMs, item.MinMs, item.AvgMs)
	}
	return nil
}

var serverInfo map[string][]string
var serverInfoTime map[string]int64
var serverInfoMutex sync.Mutex
var startTime string
var AsynNumFunc func() int

func addServerInfo(addr string, list []string) {
	defer serverInfoMutex.Unlock()
	serverInfoMutex.Lock()
	serverInfo[addr] = list
	serverInfoTime[addr] = time.Now().Unix()
}
func delServerInfo(addr string) {
	defer serverInfoMutex.Unlock()
	serverInfoMutex.Lock()
	delete(serverInfo, addr)
	delete(serverInfoTime, addr)
}

func getServerInfo() (result []struct {
	List []string
	Addr string
}) {
	defer serverInfoMutex.Unlock()
	serverInfoMutex.Lock()
	for k, item := range serverInfo {
		var resultItem struct {
			List []string
			Addr string
		}
		resultItem.List = item
		resultItem.Addr = k
		result = append(result, resultItem)
	}
	return
}

func serverInfoTick(stop *bool) {
	list := []string{}
	timeoutUnxi := time.Now().Unix() - 80
	serverInfoMutex.Lock()
	for k, unix := range serverInfoTime {
		if unix < timeoutUnxi {
			list = append(list, k)
		}
	}
	serverInfoMutex.Unlock()
	for _, key := range list {
		delServerInfo(key)
	}
}

func StartServerInfoCheck() {
	lib.RunTickStopTimer(30000, serverInfoTick) //定时器检测超时daoinfo
}
