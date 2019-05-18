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

package business

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
)

func viewDaos() {
	client, _ := utilsClient()
	var replys string
	client.Call("InnerServer.DaoServerInfo", 100, &replys)
	fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 34, "viewDaos-client Dao container service info:"))
	for _, str := range strings.Split(replys, ";") {
		fmt.Println(str)
	}

	replys = ""
	client.Call("InnerServer.DaoStatus", 100, &replys)
	client.Close()
	fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 34, "viewDaos-client Dao component connection info:"))
	for _, str := range strings.Split(replys, ";") {
		fmt.Println(str)
	}
}

func viewRuntime() {
	client, addr := utilsClient()
	if client == nil {
		fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 31, "viewRuntime-client business service not startup"))
		os.Exit(0)
	}
	var replys string
	client.Call("InnerServer.BusinessInfo", 100, &replys)
	var pid int
	client.Call("InnerServer.ProcessId", 100, &pid)
	fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 34, "service info Listening address: "+addr+", pid:"+fmt.Sprint(pid)))
	fmt.Println(replys)
	client.Close()
}

func viewQps(top int) {
	//如果top 大于1 只刷新前30
	client, _ := utilsClient()
	if client == nil {
		fmt.Println(fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 31, "viewQps-client business service not startup"))
		os.Exit(0)
	}
	var replys string
	var restart string
	client.Call("InnerServer.ComQps", 100, &replys)
	client.Call("InnerServer.ComQpsBeginTime", 100, &restart)
	client.Close()
	if top == 1 {
		fmt.Printf("%46s\n", "Com qps clear: "+restart)
		for _, str := range strings.Split(replys, "??") {
			fmt.Println(str)
		}
		return
	}
	fmt.Printf("%46s\n", "Com top 30 qps clear: "+restart)
	listqps := strings.Split(replys, "??")
	if len(listqps) > 30 {
		listqps = listqps[0:30]
	}
	for _, str := range listqps {
		fmt.Println(str)
	}
}

func printSystem() {
	if helper.InSlice(os.Args, "status") {
		viewRuntime()
		viewDaos()
		os.Exit(0)
		return
	}
	if helper.InSlice(os.Args, "qps") {
		count := 1
		if helper.InSlice(os.Args, "-t") {
			count = 180
		}
		for index := 0; index < count; index++ {
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			viewQps(count)
			time.Sleep(1 * time.Second)
		}
		os.Exit(0)
		return
	}
}

func utilsClient() (*rpc.Client, string) {
	addr := helper.Config().String("dispersed::BindAddr")
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
	addr := helper.Config().String("dispersed::BindAddr")
	if addr == "" {
		panic("undefined BindAddr")
	}
	list := strings.Split(addr, ":")
	port, _ := strconv.Atoi(list[1])
	if helper.InSlice(os.Args, "daemon") {
		lib.AppDaemon()
		os.Exit(0)
		return
	}

	if helper.InSlice(os.Args, "restart") || helper.InSlice(os.Args, "start") {
		lib.AppRestart("business", list[0], port)
		os.Exit(0)
		return
	}
	if helper.InSlice(os.Args, "stop") {
		lib.AppStop("business", list[0], port)
		os.Exit(0)
		return
	}
}
