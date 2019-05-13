# gotree
特性:熔断、基于子进程热更新\优雅关闭、rpc长连接池通信5万/qps、定时器、Queue、Sync、SQL慢查询、SQL冗余查询、各种监控、分层、强制垂直分库、基于会话的seq全局日志、等


文档和示例正在迭代中，qq:4932004


1.下载
go get -u github.com/8treenet/gotree


2.安装
cd $GOPATH/src/github.com/8treenet/gotree/

go install


3.创建项目

gotree new learning

4.学习示例的数据库连接和安装

导入数据库文件: 
source $GOPATH/src/learning/learning.sql

编辑数据库用户和密码: 
vi $GOPATH/src/learning/dao/conf/dev/db.conf


5.启动

cd $GOPATH/src/learning/dao

go run main.go

cd $GOPATH/src/learning/business

go run main.go



6.模拟网关调用查看代码

vi $GOPATH/src/learning/business/unit/gateway_test.go

7.模拟网关调用执行

go test -v -count=1 -run TestUserRegister $GOPATH/src/learning/business/unit/gateway_test.go

go test -v -count=1 -run TestStore $GOPATH/src/learning/business/unit/gateway_test.go

go test -v -count=1 -run TestShopping $GOPATH/src/learning/business/unit/gateway_test.go

go test -v -count=1 -run TestUserOrder $GOPATH/src/learning/business/unit/gateway_test.go


8.qps压测

go run $GOPATH/src/learning/business/unit/qps_press/main.go

