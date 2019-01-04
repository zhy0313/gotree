package business

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"strconv"
	"strings"

	"jryghq.cn/lib"

	"jryghq.cn/utils"
)

func viewDaos() {
	client, _ := utilsClient()
	var replys string
	client.Call("InnerServer.DaoServerInfo", 100, &replys)
	fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 34, "dao容器服务信息(刷新频次 生产1分钟, 测试同心跳时间):"))
	for _, str := range strings.Split(replys, ";") {
		fmt.Println(str)
	}

	replys = ""
	client.Call("InnerServer.DaoStatus", 100, &replys)
	client.Close()
	fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 34, "dao组件连接信息:"))
	for _, str := range strings.Split(replys, ";") {
		fmt.Println(str)
	}
}

func viewRuntime() {
	client, addr := utilsClient()
	if client == nil {
		fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 31, "本机未启动business"))
		os.Exit(0)
	}
	var replys string
	client.Call("InnerServer.BusinessInfo", 100, &replys)
	var pid int
	client.Call("InnerServer.ProcessId", 100, &pid)
	fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 34, "本服务信息 监听地址: "+addr+", pid:"+fmt.Sprint(pid)))
	fmt.Println(replys)
	client.Close()
}

func viewQps() {
	client, _ := utilsClient()
	if client == nil {
		fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 31, "本机未启动business"))
		os.Exit(0)
	}
	var replys string
	client.Call("InnerServer.DaoQps", 100, &replys)
	client.Close()
	fmt.Println("dao qps:(远程方法, 调用次数, 最高用时, 最低用时, 平均用时)")

	for _, str := range strings.Split(replys, "??") {
		fmt.Println(str)
	}
}

func printSystem() {
	if utils.InArray(os.Args, "status") {
		viewRuntime()
		viewDaos()
		os.Exit(0)
		return
	}
	if utils.InArray(os.Args, "qps") {
		viewQps()
		os.Exit(0)
		return
	}
}

func utilsClient() (*rpc.Client, string) {
	addr := utils.Config().String("BindAddr")
	if addr == "" {
		return nil, addr
	}
	list := strings.Split(addr, ":")
	port, _ := strconv.Atoi(list[1])
	for index := 0; index < 10; index++ {
		client, err := jsonrpc.Dial("tcp", fmt.Sprintf("%s:%d", list[0], port+index))
		if err != nil {
			continue
		}
		var pid int
		if client.Call("InnerServer.ProcessId", 100, &pid) == nil {
			return client, fmt.Sprintf("%s:%d", list[0], port+index)
		}
		client.Close()
	}
	return nil, ""
}

func appStart() {
	addr := utils.Config().String("BindAddr")
	if addr == "" {
		panic("undefined BindAddr")
	}
	list := strings.Split(addr, ":")
	port, _ := strconv.Atoi(list[1])
	if utils.InArray(os.Args, "daemon") {
		lib.AppDaemon()
		os.Exit(0)
		return
	}
	if utils.InArray(os.Args, "start") {
		lib.AppStart("businesss", list[0], port)
		os.Exit(0)
		return
	}
	if utils.InArray(os.Args, "restart") {
		lib.AppRestart("businesss", list[0], port)
		os.Exit(0)
		return
	}
	if utils.InArray(os.Args, "stop") {
		lib.AppStop("businesss", list[0], port)
		os.Exit(0)
		return
	}
}
