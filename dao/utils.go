package dao

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
)

func telnet() {
	if !helper.InSlice(os.Args, "telnet") {
		return
	}
	helper.Log().Debug()
	_msl.NotitySubscribe("DaoTelnet", daos()...)
	baddrs := helper.Config().String("BusinessAddrs")
	if baddrs == "" {
		helper.Log().WriteWarn("BusinessAddrs地址为空")
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}

	list := strings.Split(baddrs, ",")
	for _, item := range list {
		_, err := net.DialTimeout("tcp", item, time.Duration(2*time.Second))
		if err != nil {
			helper.Log().WriteWarn("BusinessAddrs连接失败", item)
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		}
	}
	os.Exit(0)
}

func appStart() {
	addr := helper.Config().String("BindAddr")
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
	if helper.InSlice(os.Args, "start") {
		lib.AppStart("dao", list[0], port)
		os.Exit(0)
		return
	}
	if helper.InSlice(os.Args, "restart") {
		lib.AppRestart("dao", list[0], port)
		os.Exit(0)
		return
	}
	if helper.InSlice(os.Args, "stop") {
		lib.AppStop("dao", list[0], port)
		os.Exit(0)
		return
	}
}
