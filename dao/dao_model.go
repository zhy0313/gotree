package dao

import (
	"fmt"

	"github.com/8treenet/gotree/dao/orm"
	"github.com/8treenet/gotree/helper"
	"github.com/8treenet/gotree/lib"
	"github.com/8treenet/gotree/lib/chart"

	// SQL Server 数据库支持
	_ "github.com/denisenkom/go-mssqldb"
	// MySQL 数据库支持
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/8treenet/gotree/lib/rpc"

	_ "github.com/go-sql-driver/mysql"
)

type DaoModel struct {
	lib.Object
	open    bool
	daoName string
	node    *chart.Node
}

var (
	node               = new(chart.Node).Node(50)
	modelProfiler      = false
	modelProfilerCount = map[string]int{}
	modelProfilerSync  sync.Mutex
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDriver("mssql", orm.DRMySQL)
}

func (self *DaoModel) DaoModel(child interface{}) *DaoModel {
	self.Object.Gotree(self)
	self.AddChild(self, child)
	self.open = false
	self.daoName = ""

	self.AddSubscribe("DaoTelnet", self.daoTelnet)
	self.AddSubscribe("ModelOn", self.modelOn)
	return self
}

type Conn interface {
	Raw(query string, args ...interface{}) orm.RawSeter
}

//Orm 获取orm
func (self *DaoModel) Conn() Conn {
	if !self.open {
		helper.Log().WriteError("model error: 未开启dao:" + self.daoName)
		panic("model error: 未开启dao")
	}
	if self.daoName == "" {
		helper.Log().WriteError("这是一个未注册的dao")
		return nil
	}
	o := orm.New(self.daoName)
	if modelProfiler {
		o.RawCallBack(func(sql string, args []interface{}) {
			self.profiler(sql, args...)
		})
	}
	return o
}

//TestOn 单元测试 开启
func (self *DaoModel) TestOn() {
	mode := helper.Config().String("sys::Mode")
	if mode == "prod" {
		helper.Log().WriteError("生产环境不可以使用单元测试model")
		panic("生产环境不可以使用单元测试model")
	}
	rpc.GoDict().Set("bseq", "ModelUnit")
	self.DaoInit()
	self.ormOn()
}

//daoOn 开启回调
func (self *DaoModel) daoTelnet(args ...interface{}) {
	for _, arg := range args {
		dao := arg.(daoNode)
		if dao.Name == self.daoName {
			self.ormOn()
			return
		}
	}
}

//modelOn 开启回调
func (self *DaoModel) modelOn(arg ...interface{}) {
	daoName := arg[0].(string)
	if daoName == self.daoName {
		self.ormOn()
	}
}

//Open
func (self *DaoModel) Open() bool {
	return self.open
}

//ormOn 开启orm
func (self *DaoModel) ormOn() {
	self.open = true

	if !connectDao(self.daoName + "model") {
		return
	}
	//处理连接
	driver := "mysql"
	dbconfig := helper.Config().String("mysql::" + self.daoName)
	dbMaxIdleConns := helper.Config().String("mysql::" + self.daoName + "MaxIdleConns")
	dbMaxOpenConns := helper.Config().String("mysql::" + self.daoName + "MaxOpenConns")
	if dbconfig == "" {
		driver = "mssql"
		dbconfig = helper.Config().String("mssql::" + self.daoName)
		dbMaxIdleConns = helper.Config().String("mssql::" + self.daoName + "MaxIdleConns")
		dbMaxOpenConns = helper.Config().String("mssql::" + self.daoName + "MaxOpenConns")
	}

	if dbconfig == "" {
		helper.Log().WriteError(self.daoName + ":数据库配置信息不存在")
		panic(self.daoName + ":数据库配置信息不存在")
	}
	_, err := orm.GetDB(self.daoName)
	if err == nil {
		//已注册
		return
	}
	if dbMaxIdleConns == "" {
		dbMaxIdleConns = helper.Config().DefaultString("sys::SqlMaxIdleConns", "1")
	}
	if dbMaxOpenConns == "" {
		dbMaxOpenConns = helper.Config().DefaultString("sys::SqlMaxOpenConns", "2")
	}
	maxIdleConns, ei := strconv.Atoi(dbMaxIdleConns)
	maxOpenConns, eo := strconv.Atoi(dbMaxOpenConns)
	if ei != nil || eo != nil || maxIdleConns == 0 || maxOpenConns == 0 || maxIdleConns > maxOpenConns {
		helper.Log().WriteError("连接dao sql:" + self.daoName + "失败, 错误原因: MaxIdleConns 或MaxOpenConns 参数错误")
		panic("连接dao sql:" + self.daoName + "失败, 错误原因: MaxIdleConns 或MaxOpenConns 参数错误")
	}
	helper.Log().WriteInfo("connect database: MaxIdleConns:" + dbMaxIdleConns + ", MaxOpenConns:" + dbMaxOpenConns + ", config:" + dbconfig)
	err = orm.RegisterDataBase(self.daoName, driver, dbconfig, maxIdleConns, maxOpenConns)
	if err != nil {
		helper.Log().WriteError("连接dao sql:" + self.daoName + "失败, 错误原因:" + err.Error())
		panic("连接dao sql:" + self.daoName + "失败, 错误原因:" + err.Error())
	}
}

//FormatArray 格式化in参数
func (self *DaoModel) FormatArray(args ...interface{}) (list []interface{}) {
	for _, item := range args {
		slice := reflect.ValueOf(item)
		if slice.Kind() != reflect.Slice {
			list = append(list, item)
			continue
		}

		//slice组合
		items, ok := takeSliceArg(item)
		if !ok {
			continue
		}
		for _, sliceItem := range items {
			list = append(list, sliceItem)
		}

	}
	return
}

//FormatPlaceholder 格式化in参数占位符
func (self *DaoModel) FormatPlaceholder(arg interface{}) string {
	slice := reflect.ValueOf(arg)
	if slice.Kind() != reflect.Slice {
		return ""
	}

	result := []string{}
	c := slice.Len()
	for i := 0; i < c; i++ {
		result = append(result, "?")
	}
	return strings.Join(result, ",")
}

func (self *DaoModel) Hash(data interface{}, mod int) int {
	return node.HashNodeSum(data, mod)
}

//profiler 分析
func (self *DaoModel) profiler(ssql string, args ...interface{}) {
	dict := rpc.GoDict()
	if dict == nil || dict.Get("bseq") == nil {
		return
	}
	if strings.Contains(ssql, "EXPLAIN") {
		return
	}
	sourceSql := fmt.Sprintf(strings.Replace(ssql, "?", "%v", -1), args...)
	sql := strings.ToLower(sourceSql)
	if strings.Contains(sql, "delete") || strings.Contains(sql, "update") || strings.Contains(sql, "insert") || strings.Contains(sql, "count") || strings.Contains(sql, "sum") || strings.Contains(sql, "max") {
		helper.Log().WriteInfo("sql profiler:", sourceSql)
		return
	}
	var explain []struct {
		Table string
		Type  string
	}
	explainLog := ""
	o := orm.New(self.daoName)
	_, e := o.Raw("EXPLAIN "+ssql, args...).QueryRows(&explain)
	tables := []string{}
	warn := false
	if e == nil {
		for index := 0; index < len(explain); index++ {
			if index > 0 {
				explainLog += " "
			}
			etype := strings.ToLower(explain[index].Type)
			if etype == "all" || etype == "index" {
				warn = true
			}
			explainLog += explain[index].Table + "表" + explain[index].Type + "级别"
			tables = append(tables, explain[index].Table)
		}
	}
	if explainLog != "" {
		explainLog = "explain :(" + explainLog + ")"
	}
	if warn {
		helper.Log().WriteWarn("sql profiler:", explainLog, "source :("+sourceSql+")")
	} else {
		helper.Log().WriteInfo("sql profiler:", explainLog, "source :("+sourceSql+")")
	}

	bseq := dict.Get("bseq")
	str, ok := bseq.(string)
	if !ok {
		return
	}

	for _, item := range tables {
		table := item
		if profilerSync(str + "_" + table) {
			continue
		}
		go func() {
			time.Sleep(3 * time.Second)
			var tableCount int
			var ok bool
			modelProfilerSync.Lock()
			tableCount, ok = modelProfilerCount[str+"_"+table]
			delete(modelProfilerCount, str+"_"+table)
			modelProfilerSync.Unlock()
			if tableCount > 1 && ok {
				helper.Log().WriteWarn("在一个bseq:" + str + " 内读取表'" + table + "' " + fmt.Sprint(tableCount) + "次")
			}
		}()
	}
}

func profilerSync(key string) (exist bool) {
	defer modelProfilerSync.Unlock()
	modelProfilerSync.Lock()
	exist = false
	v, ok := modelProfilerCount[key]
	if ok {
		exist = true
		modelProfilerCount[key] = v + 1
		return
	}
	modelProfilerCount[key] = 1
	return
}

func (self *DaoModel) Connections(m map[string]int) {
	if !self.open {
		return
	}
	db, err := orm.GetDB(self.daoName)
	if err != nil {
		return
	}
	m[self.daoName] = db.Stats().OpenConnections
}

func (self *DaoModel) DaoInit() {
	if self.daoName == "" {
		self.daoName = self.TopChild().(daoName).Dao()
	}
}

func takeSliceArg(arg interface{}) (out []interface{}, ok bool) {

	slice := reflect.ValueOf(arg)
	if slice.Kind() != reflect.Slice {
		return nil, false
	}

	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}
