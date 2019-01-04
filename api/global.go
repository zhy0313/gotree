package api

import (
	"jryghq.cn/lib"
	"jryghq.cn/remote_call"
)

var _scl *lib.ServiceLocator //controller服务定位器
func init() {
	_scl = new(lib.ServiceLocator).ServiceLocator()
	//注册相关组件
	_scl.AddComponent(new(remote_call.InnerMaster).InnerMaster())
	_scl.AddComponent(new(remote_call.InnerClient).InnerClient())
}

func AppendBusiness(remoteAddr string) {
	var ic *remote_call.InnerClient
	_scl.GetComponent(&ic)
	ic.AddRemoteAddr(remoteAddr)
}

//启动连接器 args[0]=最大并发数, args[1]=失败重试次数
func Run(args ...int) {
	var (
		ic          *remote_call.InnerClient
		concurrency int = 2048
		retry       int = 3
	)
	if len(args) > 0 {
		concurrency = args[0]
	}
	if len(args) > 1 {
		retry = args[1]
	}
	_scl.GetComponent(&ic)
	_scl.AddComponent(new(remote_call.RpcClient).RpcClient(concurrency, retry))
	go ic.Run()
}

func RpcClient() *remote_call.RpcClient {
	var rc *remote_call.RpcClient
	_scl.GetComponent(&rc)
	return rc
}
