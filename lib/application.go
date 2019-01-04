package lib

import (
	"bytes"
	"fmt"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/8treenet/gotree/helper"
)

func AppDaemon() {
	for {
		var stderr bytes.Buffer
		cmd := exec.Command(os.Args[0])
		cmd.Stderr = &stderr
		cmd.Start()
		err := cmd.Wait()
		if err == nil {
			os.Exit(0)
		}

		helper.Log().WriteDaemonError("异常错误,10秒后开始重启 :", stderr.String())
		time.Sleep(10 * time.Second)
	}
}

func AppStart(name, addr string, port int) {
	over := make(chan bool, 1)
	go func() {
		for i := 0; i != 10; i = i + 1 {
			fmt.Fprintf(os.Stdout, "启动进度 : %%%d\r", i*10)
			time.Sleep(time.Millisecond * 300)
		}
		over <- true
	}()

	for index := 0; index < 10; index++ {
		_, err := jsonrpc.Dial("tcp", fmt.Sprintf("%s:%d", addr, port+index))
		if err == nil {
			fmt.Println(name+" 正在运行中 :", fmt.Sprintf("%s:%d", addr, port+index))
			os.Exit(0)
		}
	}

	cmdStart := exec.Command("nohup", os.Args[0], "daemon", "&")
	err := cmdStart.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	<-over
	if err == nil {
		fmt.Printf("启动进度 :%%%d\n", 100)
		fmt.Println("启动"+name+" daemon pid:", cmdStart.Process.Pid)
	}
}

func AppStop(name, addr string, port int) {
	for index := 0; index < 10; index++ {
		client, err := jsonrpc.Dial("tcp", fmt.Sprintf("%s:%d", addr, port+index))
		if err != nil {
			continue
		}
		var pid int
		if client.Call("InnerServer.ProcessId", 100, &pid) != nil {
			continue
		}
		process, err := os.FindProcess(pid)
		if err != nil {
			continue
		}
		fmt.Println("关闭"+name+" pid:", pid)
		process.Signal(syscall.SIGINT)
	}
}

func AppRestart(name, addr string, port int) {
	over := make(chan bool, 1)
	go func() {
		for i := 0; i != 10; i = i + 1 {
			fmt.Fprintf(os.Stdout, "重启进度 : %%%d\r", i*10)
			time.Sleep(time.Millisecond * 300)
		}
		over <- true
	}()

	AppStop(name, addr, port)
	cmdStart := exec.Command("nohup", os.Args[0], "daemon", "&")
	err := cmdStart.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	<-over
	if err == nil {
		fmt.Printf("重启进度 : %%%d\n", 100)
		fmt.Println("启动"+name+" daemon pid:", cmdStart.Process.Pid)
	}
}
