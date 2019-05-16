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

### 描述
+ Business 主要用于逻辑功能处理等。均衡负载部署多台，为网关提供服务。 目录结构在 learning/business。
+ Dao 主要用于数据功能处理，组织低级数据提供给上游 business。负载方式较多，可根据数据分布来设计。可通过配置来开启组件。目录结构在 learning/dao。
+ Protocol 通信协议 business_cmd/value 作用于 Api网关和 Business 通信。 dso_cmd/value 作用于 Business 和 Dao 通信。 目录结构在 learning/protocol。

> 3台网关、2台business、3台dao 组成的集群
>
> 服务器 | 服务器 | 服务器
> -------------|-------------|-------------
> APIGateway-1 | APIGateway-2 | TcpGateway-1
> Business-1   |              | Business-2
> Dao-1        | Dao-2        | Dao-3

### 分层
架构主要分为4层。第一层基类 __BusinessController__，作为 Business 的入口控制器, 主要职责有组织和协调Service、逻辑处理。 第二层基类 __BusinessService__, 作为 __BusinessController__ 的下沉层， 主要下沉的职责有拆分、治理、解耦、复用、使用Dao。 第三次基类 __DaoController__ ，作为 Dao 的入口控制器，主要职责有组织数据、解耦数据和逻辑、抽象数据源、使用数据源。 第四层数据源基类 __DaoModel__ 数据库表模型数据源基类、 __DaoMemory__ 内存数据源基类、 __DaoCache__ redis数据源基类、 __DaoApi__ Http数据源基类。

### 使用 gateway
```go
    // 1. 模拟api网关调用，等同 beego、gin 等api gateway, 以及 tcp 网关项目.
	// 2. 实际应用中 business 分布在多个物理机器.  gateway.AppendBusiness 因填写多机器的内网ip.
	func main() {
		gateway.AppendBusiness("192.168.1.1:8888")
		gateway.AppendBusiness("192.168.1.2:8888")
		gateway.AppendBusiness("192.168.1.3:8888")
        gateway.Run()
        
        //创建 business 调用命令
        cmd := new(business_cmd.Store).Gotree([]int64{1, 2})
        //创建 business 返回数据
        value := business_value.Store{}
        
        //设置自定义头，透传数据
        cmd.SetHeader("", "")
        //可直接设置http头
        cmd.SetHttpHeader(head) 
        gateway.RpcClient().Call(cmd, &value)
        //response value ....
	}
```

### 使用 BusinessController  
```go
    /* 
         learning/business/controllers/controllers.go
    */
    func init() {
        //注册 ProductController 控制器
        business.RegisterController(new(ProductController).Gotree())
    }

    //定义一个电商的商品控制器。
    type ProductController struct {
        //继承 business 控制器的基类
	    business.BusinessController
    }

    //这个是 gotree 风格的构造函数，底层通过指针原型链，可以实现多态，和基础类的支持。
    func (self *ProductController) Gotree() *ProductController {
        self.BusinessController.Gotree(self)
        return self
    }

    //每一个 APIGateway 触发的 rpc 动作调用 都会创造一个 ProductController 对象， 并且调用 Prepare。
    func (self *ProductController) Prepare(method string, argv interface{}) {
        //调用父类 Prepare
        self.BusinessController.Prepare(method, argv)
        helper.Log().WriteInfo("Prepare:", method, argv)
    }

    //每一个 APIGateway 触发的 rpc 动作调用结束 都会触发 Finish。
    func (self *ProductController) Finish(method string, reply interface{}, e error) {
        self.BusinessController.Finish(method, reply, e)
        //打印日志
        helper.Log().WriteInfo("Finish:", method, fmt.Sprint(reply), e)
    }

    /*
        这是一个查看商品列表的 Action, cmd 是入参，result 是出参， 在 protocol中定义，下文详细介绍。
    */
    func (self *ProductController) Store(cmd business_cmd.Store, result *business_value.Store) (e error) {
        var (
            productSer *service.Product
        )
        *result = business_value.Store{}

        //通过 父类 Service 方法 取出 service.Product 类型的服务对象。
        //因为 go 没有泛型，实现服务定位器模式，只可依赖二级指针，不用管原理，直接取。
        self.Service(&productSer)

        //使用服务的 Store 方法 读取出商品数据， 并且赋值给出参 result 
        result.List, e = productSer.Store()
        return
    }
```

### 使用 BusinessService
```go
    /* 
         learning/business/service/product.go
    */
    func init() {
        //注册 service
        business.RegisterService(new(Product).Gotree())
    }

    type Product struct {
        //继承 BusinessService 基类
        business.BusinessService
    }

    // gotree 风格构造
    func (self *Product) Gotree() *Product {
        self.BusinessService.Gotree(self)
        return self
    }

    // 读取商品服务 返回一个商品信息匿名结构体数组。
    func (self *Product) Store() (result []struct {
        Id    int64 //商品 id
        Price int64 //商品价格
        Desc  string //商品描述
    }, e error) {

        //创建 dao调用命令
        cmdPt := new(dao_cmd.ProductGetList).Gotree([]int64{1, 2})
        //创建 dao返回数据
        store := dao_value.ProductGetList{}

        //CallDao 调用远程服务器的dao 入参cmdPt 出参store
        e = self.CallDao(cmdPt, &store)
        if e == helper.ErrBreaker {
            //熔断处理
            helper.Log().WriteInfo("Store ErrBreaker")
            return
        }
        result = store.List
        return
    }
```

### 使用  Protocol
```go
    /* 
        business 远程调用 learning/protocol/business_cmd/product.go
        dao 远程调用 learning/protocol/dao_cmd/product.go
    */

```