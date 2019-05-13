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
package api

import (
	"github.com/8treenet/gotree/lib"
	"github.com/8treenet/gotree/remote_call"
)

var _scl *lib.ServiceLocator //controller服务定位器
func init() {
	_scl = new(lib.ServiceLocator).Gotree()
	//注册相关组件
	_scl.AddComponent(new(remote_call.InnerMaster).Gotree())
	_scl.AddComponent(new(remote_call.InnerClient).Gotree())
}

func AppendBusiness(remoteAddr string) {
	var ic *remote_call.InnerClient
	_scl.GetComponent(&ic)
	ic.AddRemoteAddr(remoteAddr)
}

//启动连接器 args[0]=最大并发数, args[1]=call business 超时时间
func Run(args ...int) {
	var (
		ic          *remote_call.InnerClient
		concurrency int = 8192
		callTimeout int = 12
	)
	if len(args) > 0 {
		concurrency = args[0]
	}
	if len(args) > 1 {
		callTimeout = args[1]
	}
	_scl.GetComponent(&ic)
	_scl.AddComponent(new(remote_call.RpcClient).Gotree(concurrency, callTimeout))
	rpcBreak := new(remote_call.RpcBreak).Gotree()
	_scl.AddComponent(rpcBreak)
	rpcBreak.RunTick()
	go ic.Run()
}

func RpcClient() *remote_call.RpcClient {
	var rc *remote_call.RpcClient
	_scl.GetComponent(&rc)
	return rc
}
