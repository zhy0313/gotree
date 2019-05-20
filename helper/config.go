package helper

import (
	"github.com/8treenet/gotree/lib/config"
)

var _conf config.Configer

func LoadConfig(project string) (e error) {
	defer func() {
		if e != nil {
			panic("Config file not Found \n ./conf/app.conf \n /usr/local/etc/+project+/app.conf \n /etc/+project+/app.conf")
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
	_conf, e = config.NewConfig("ini", "/usr/local/etc/"+project+"/app.conf")
	if e == nil {
		return
	}

	_conf, e = config.NewConfig("ini", "/etc/"+project+"/app.conf")
	return
}

func Config() config.Configer {
	return _conf
}
