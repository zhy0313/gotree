package dao

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"jryghq.cn/lib"
	"jryghq.cn/utils"
)

func telnet() {
	if !utils.InArray(os.Args, "telnet") {
		return
	}
	utils.Log().Debug()
	_msl.NotitySubscribe("DaoTelnet", daos()...)
	baddrs := utils.Config().String("BusinessAddrs")
	if baddrs == "" {
		utils.Log().WriteWarn("BusinessAddrs地址为空")
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}

	list := strings.Split(baddrs, ",")
	for _, item := range list {
		_, err := net.DialTimeout("tcp", item, time.Duration(2*time.Second))
		if err != nil {
			utils.Log().WriteWarn("BusinessAddrs连接失败", item)
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		}
	}
	os.Exit(0)
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
		lib.AppStart("dao", list[0], port)
		os.Exit(0)
		return
	}
	if utils.InArray(os.Args, "restart") {
		lib.AppRestart("dao", list[0], port)
		os.Exit(0)
		return
	}
	if utils.InArray(os.Args, "stop") {
		lib.AppStop("dao", list[0], port)
		os.Exit(0)
		return
	}
}
