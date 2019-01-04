package utils

import (
	"jryghq.cn/lib/config"
)

var _conf config.Configer

func LoadConfig(project string) (e error) {
	defer func() {
		if e != nil {
			panic("找不到配置文件 \n ./conf/app.conf \n /usr/local/etc/jryg/+project+/app.conf \n /etc/jryg/+project+/app.conf")
		}
	}()
	_conf, e = config.NewConfig("ini", "conf/app.conf")
	if e == nil {
		return
	}
	_conf, e = config.NewConfig("ini", "../conf/app.conf")
	if e == nil {
		return
	}
	_conf, e = config.NewConfig("ini", "/usr/local/etc/jryg/"+project+"/app.conf")
	if e == nil {
		return
	}

	_conf, e = config.NewConfig("ini", "/etc/jryg/"+project+"/app.conf")
	return
}

func Config() config.Configer {
	return _conf
}
