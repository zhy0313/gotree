package utils

import (
	"net"
	"os"
)

// 获取本机IP地址
func LocalIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		os.Exit(1)
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
