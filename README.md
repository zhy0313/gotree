# gotree

gotree 是一个垂直分布式框架。 gotree 的目标是轻松开发分布式服务，解放开发者心智。

## 特性
* 熔断
* 平滑升级
* rpc 通信(c50k)
* 定时器
* SQL 慢查询监控
* SQL 冗余监控
* 分层
* 强制垂直分库
* 基于 seq 串行的全局日志，多机器之间
* 单元测试
* 督程
* 一致性哈希、主从、随机、均衡等负载方式

## 快速使用

1. 获取 gotree。
```sh
$ go get -u github.com/8treenet/gotree
```

2. 安装 gotree。
```sh
$ cd $GOPATH/src/github.com/8treenet/gotree/
$ go install
```

3. 创建 learning 项目。
```sh
$ gotree new learning
```

4. learning 项目数据库安装、数据库用户密码配置。使用 source 或工具安装 learning.sql。
```sh
$ mysql > source $GOPATH/src/learning/learning.sql
$ vi $GOPATH/src/learning/dao/conf/dev/db.conf
```

5. 启动 dao服务、 business 服务。
```sh
$ cd $GOPATH/src/learning/dao
$ go run main.go
$ cd $GOPATH/src/learning/business
$ go run main.go
```

6. 模拟网关执行调用，请查看代码。 代码位于 $GOPATH/src/learning/business/unit/gateway_test.go
```sh
$ go test -v -count=1 -run TestUserRegister $GOPATH/src/learning/business/unit/gateway_test.go
$ go test -v -count=1 -run TestStore $GOPATH/src/learning/business/unit/gateway_test.go
$ go test -v -count=1 -run TestShopping $GOPATH/src/learning/business/unit/gateway_test.go
$ go test -v -count=1 -run TestUserOrder $GOPATH/src/learning/business/unit/gateway_test.go
```

7. qps 压测
```sh
$ go run $GOPATH/src/learning/business/unit/qps_press/main.go
```


## 快速入门  

### 使用 BusinessController  
```go
    //定义 business 控制器，继承自 business.BusinessController
    type ProductController struct {
	    business.BusinessController
    }

    func init() {
        business.RegisterController(new(ProductController).Gotree())
    }

    func (self *ProductController) Gotree() *ProductController {
        self.BusinessController.Gotree(self)
        return self
    }

    func (self *ProductController) Prepare(method string, argv interface{}) {
        self.BusinessController.Prepare(method, argv)
        helper.Log().WriteInfo("Prepare:", method, argv)
    }

    func (self *ProductController) Finish(method string, reply interface{}, e error) {
        self.BusinessController.Finish(method, reply, e)
        helper.Log().WriteInfo("Finish:", method, fmt.Sprint(reply), e)
    }
```