# 轻量、高性能
daog是轻量级的数据库访问组件，它并不能称之为orm组件，仅仅提供了一组函数用以实现常用的数据库访问功能。
它是高性能的，与原生的使用sql包函数相比，没有性能损耗，这是因为，它并没有使用反射技术，而是使用编译技术把create table sql语句编译成daog需要的go代码。
它目前仅支持mysql。

设计思路来源于java的[orm框架sampleGenericDao](https://github.com/tiandarwin/simpleGenericDao)和protobuf的编译思路。之所以选择编译
而没有使用反射，是因为基于编译的抽象没有性能损耗。

## 编译组件

[compilex](https://github.com/rolandhe/compilex) 是编译create table语句文件的工具，使用如下语句可以编译：

```
  ./compilex -i="sql file" -pkg packageName -o xxx/xx
```

编译完成后，把整个packageName文件夹copy到你的项目中即可。每一个表生成两个文件：
 * 主文件，以表名的驼峰格式命名，包含映射表的struct，和两个元数据对象，主文件不要修改
 * 扩展文件，以表名的驼峰格式+"-ext" 命名，这个文件可以修改，开发者可以自由扩展针对该表的数据访问功能。daog支持分表，分表函数需要在该文件的init函数中设置。

# 使用

## 核心概念
使用之前先要了解Datasource、TransContext、TableMeta

### Datasource
数据源，用于描述一个mysql database，这儿的database指的是您使用create database创建出来的逻辑库。Datasource提供获链接、关闭库函数，也可以配置在改数据源上操作数据是否要输出执行sql日志。
* 使用NewDatasource或NewShardingDatasource函数来创建Datasource对象
* 数据库相关配置使用DbConf描述

### TransContext
事务的执行上下文，所有的数据库操作都应该在一个数据上下文中执行，所有操作完成后必须调用Complete函数来结束事务上下文，一旦结束该上下文将不能再被使用。
* 使用NewTransContext和NewTransContextWithSharding函数来创建事务上下文，二者的区别是是否支持分库分表，分库分表不必同时进行，可以只分库，也可以只分表，不需要的sharding Key传入nil即可。
* 支持3中事务类型：没有事务、只读事务、写事务，txrequest包定义了对应的常量
* 必须调用Complete方法来结束事务上下文,一般使用defer语句来结束事务上下文

### TableMeta
go struct，用来描述数据表及对应go 对象信息信息，在go程序中一张数据库表需要对应的一个struct来描述，包括：
* 表名及对应的struct的名称
* 字段信息：字段名，对应struct的属性信息，数据类型信息
* 用于根据字段名查找struct对象的属性值或者属性指针，采用该方法避免使用反射来为属性赋值，或者读取属性的值

每张表的TableMeta对象是通过compliex来生成的，无需手工创建

### 其他概念
#### Matcher
用于拼接sql条件的工具，直接进行字符串拼接往往会产生错误，而且错误只能在运行时被发现，提供工具来避免这种情况。
Matcher至支持多个条件组合.
* Matcher内置了eq,like,between,gt,lt等过个快捷条件生成，支持组合新的Matcher，也支持您自己实现新的条件，直接实现 SQLCond接口即可。
* 使用NewMatcher、NewAndMatcher、NewOrMatcher来创建对象

#### TableFields
这是逻辑概念，compilex会在每张表对应的主go文件中创建一个匿名struct对象。该对象记录了数据库的字段名称，以便于利用Matcher拼接sql

#### Modifier
顾名思义，用于update表字段，它描述了一组字段名与对应值对，用于拼接update语句


## 使用方式
daog提供了两种数据方式形式：
* 直接使用daog提供的函数，比如：
```
func Insert[T any](tc *TransContextins *T, meta *TableMeta[T]) (int64, error) 
func QueryListMatcher[T any](tc *TransContext,m Matcher, meta *TableMeta[T]) ([]*T, error)
```
* 使用QuickDao接口，该接口支持模板参数，每个编译好的主文件中都有类似GroupInfoDao的变量，该变量是QuickDao[GroupInfo]类型，包含一个实现QuickDao[GroupInfo]的匿名struct对象，可以直接使用它来操作数据库，相对于使用函数，它少传递了TableMeta对象

## 使用实例
可以参照代码的example

### 编译create table语句

```
create table group_info (
    id bigint(20) not null AUTO_INCREMENT primary key,
    `name` varchar(200) not null comment 'user name',
    main_data json not null,
    create_at datetime not null,
    total_amount decimal(10,2) not null
) ENGINE=innodb CHARACTER SET utf8mb4 comment 'group info';
```

编译出的主代码：

```
package dal

import (
	"github.com/rolandhe/daog"
	dbtime "github.com/rolandhe/daog/time"
	"github.com/shopspring/decimal"
)

var GroupInfoFields = struct {
	Id string
	Name string
	MainData string
	CreateAt string
	TotalAmount string

}{
	"id",
	"name",
	"main_data",
	"create_at",
	"total_amount",

}

var  GroupInfoMeta = &daog.TableMeta[GroupInfo]{
	InstanceFunc: func() *GroupInfo{
		return &GroupInfo{}
	},
	Table: "group_info",
	Columns: []string {
		"id",
		"name",
		"main_data",
		"create_at",
		"total_amount",

	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string,ins *GroupInfo,point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "name" == columnName {
			if point {
				return &ins.Name
			}
			return ins.Name
		}
		if "main_data" == columnName {
			if point {
				return &ins.MainData
			}
			return ins.MainData
		}
		if "create_at" == columnName {
			if point {
				return &ins.CreateAt
			}
			return ins.CreateAt
		}
		if "total_amount" == columnName {
			if point {
				return &ins.TotalAmount
			}
			return ins.TotalAmount
		}

		return nil
	},
}

var GroupInfoDao daog.QuickDao[GroupInfo] = &struct {
	daog.QuickDao[GroupInfo]
}{
	daog.NewBaseQuickDao(GroupInfoMeta),
}

type GroupInfo struct {
	Id int64
	Name string
	MainData string
	CreateAt dbtime.NormalDatetime
	TotalAmount decimal.Decimal

}
```

copy到您的工程中

### 构建全局的Datasource对象

```
var datasource daog.Datasource

func init() {
	conf := &daog.DbConf{
		DbUrl:  "root:12345678@tcp(localhost:3306)/daog?parseTime=true&timeout=1s&readTimeout=2s&writeTimeout=2s",
		LogSQL: true,
	}
	var err error
	datasource, err = daog.NewDatasource(conf)
	if err != nil {
		log.Fatalln(err)
	}
}
```

### 构建事务上下文

```
// 构建写事务
tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-100099")
if err != nil {
    fmt.Println(err)
    return
}
// 构建读事务
tc, err := daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-100099")
if err != nil {
    fmt.Println(err)
    return
}

// 构建无事务
tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-100099")
if err != nil {
    fmt.Println(err)
    return
}

```

### 函数模式
#### 写表

```
func create() {
	amount, err := decimal.NewFromString("128.0")
	if err != nil {
		fmt.Println(err)
		return
	}
	t := &dal.GroupInfo{
		Name:        "roland-one",
		MainData:    `{"a":102}`,
		Content:     "hello world!!",
		BinData:     []byte("byte data"),
		CreateAt:    ttypes.NormalDatetime(time.Now()),
		TotalAmount: amount,
		Remark:      *ttypes.FromString("haha"),
	}

	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	}
	daog.AutoTrans(tcCreate, func(tc *daog.TransContext) error {
		affect, err := daog.Insert(tc, t, dal.GroupInfoMeta)
		fmt.Println(affect, t.Id, err)
		if err != nil {
			return err
		}
		t.Name = "rolandx"
		af, err := daog.Update(tc, t, dal.GroupInfoMeta)
		fmt.Println(af, err)
		if err != nil {
			return err
		}
		return nil
	})
}
```

#### 读取数据

```
func queryByIds() {
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-1001")
	}
	gs, err := daog.AutoTransWithResult(tcCreate, func(tc *daog.TransContext) ([]*dal.GroupInfo, error) {
		return daog.GetByIds(tc, []int64{1, 2}, dal.GroupInfoMeta)
	})

	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByIds", string(j))
	fmt.Println(gs)
}
```

```
func queryAll() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// tc 创建后必须马上跟上 defer func， 如果这之间有return或者panic，连接将被泄露
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	// 必须使用匿名函数，不能使用 tc.Complete(err)， 因为defer 后面函数的参数在执行defer语句是就会被确定
	defer func() {
		// 注意：后面代码的error都要使用err变量来接收，否则在发生错误的情况下，事务不会被回滚
		tc.CompleteWithPanic(err, recover())
	}()
	gs, err := daog.GetAll(tc, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryAll", string(j))
	fmt.Println(gs)
}
```

#### 根据Matcher读取
```
func queryByMatcher() {
	matcher := daog.NewMatcher().Like(dal.GroupInfoFields.Name, "roland", daog.LikeStyleRight).Lt(dal.GroupInfoFields.Id, 4)
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	}
	gs, err := daog.AutoTransWithResult(tcCreate, func(tc *daog.TransContext) ([]*dal.GroupInfo, error) {
		return daog.QueryListMatcher(tc, matcher, dal.GroupInfoMeta)
	})

	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByMatcher", string(j))
	fmt.Println(gs)
}
```

#### 先读再写

```
func update() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-100099")
	if err != nil {
		fmt.Println(err)
		return
	}
	// tc 创建后必须马上跟上 defer func， 如果这之间有return或者panic，连接将被泄露
	// 必须使用匿名函数，不能使用 tc.CompleteWithPanic(err)， 因为defer 后面函数的参数在执行defer语句是就会被确定
	defer func() {
		// 注意：后面代码的error都要使用err变量来接收，否则在发生错误的情况下，事务不会被回滚
		tc.CompleteWithPanic(err, recover())
	}()
	g, err := daog.GetById(tc, 5, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(g)
	fmt.Println("query", string(j))

	g.Name = "Eric"
	af, err := daog.Update(tc, g, dal.GroupInfoMeta)
	fmt.Println(af, err)

}
```

#### 删除

```
func deleteById() {
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	}
	daog.AutoTrans(tcCreate, func(tc *daog.TransContext) error {
		g, err := daog.DeleteById(tc, 2, dal.GroupInfoMeta)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("delete", g)
		return err
	})
}
```

### 两种种事务管理
#### 自动委模式
把所有的事务处理代码都内置到一个函数里，当然这个函数可以是匿名函数，也可以是命名函数，然后调用daog.AutoTrans 或者daog.AutoTransWithResult 来执行业务处理函数，同时事务上下文的构建也被内置到一个函数里，这个函数可以是匿名，也可以是命名的。
在这种方式下，事务的创建及最终的提交或回滚不需要再额外处理。使用方式如下，如果你不需要业务返回值，可以调用daog.AutoTrans。

```
func createUserUseAutoTrans() {
	t := &dal.UserInfo{
		Name:     "roland",
		CreateAt: ttypes.NormalDatetime(time.Now()),
		ModifyAt: *ttypes.FromDatetime(time.Now()),
	}
	affect, err := daog.AutoTransWithResult[int64](func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	}, func(tc *daog.TransContext) (int64, error) {
		return daog.Insert(tc, t, dal.UserInfoMeta)
	})
	fmt.Println(affect, t.Id, err)
}
```
#### 自行处理式
在创建TransContext后，需要手工处理事务的结束，必须通过一个匿名deffer函数来结束事务，匿名函数里调用 tc.CompleteWithPanic(err, recover()) 来最终结束事务。

需要注意的是：
* err在当前的函数里必须是全局的
* 每一个会产生err的函数调用，都必须使用该全局err来接收，比如： err = XXX()
* tc 创建后必须马上跟上 defer func， 如果这之间有return或者panic，连接将被泄露

```
func queryByMatcherOrder() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// tc 创建后必须马上跟上 defer func， 如果这之间有return或者panic，连接将被泄露
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	// 必须使用匿名函数，不能使用 tc.Complete(err)， 因为defer 后面函数的参数在执行defer语句是就会被确定
	defer func() {
		// 注意：后面代码的error都要使用err变量来接收，否则在发生错误的情况下，事务不会被回滚
		tc.CompleteWithPanic(err, recover())
	}()
	matcher := daog.NewMatcher().Like(dal.GroupInfoFields.Name, "roland", daog.LikeStyleLeft).Lt(dal.GroupInfoFields.Id, 4)
	gs, err := daog.QueryListMatcher(tc, matcher, dal.GroupInfoMeta, daog.NewDescOrder(dal.GroupInfoFields.Id))
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByMatcherOrder", string(j))
	fmt.Println(gs)
}

```

### QuickDao接口模式
使用方式和函数模式非常类似，只是少传递TableMeta参数，以下以一个query示例来说明一下。

### 查询数据

```
func queryByIdsUsingDao() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 注意： 创建好tc到调用WrapTransWithResult或者 WrapTrans之间不能返回或者panic，否则，会导致连接不释放
	// 如果你的场景是需要返回或者panic，你可以是使用 queryAll 中的事务处理方式
	gs, err := daog.WrapTransWithResult(tc, func(tc *daog.TransContext) ([]*dal.GroupInfo, error) {
		return dal.GroupInfoDao.GetByIds(tc, []int64{1, 2})
	})

	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByIdsUsingDao", string(j))
	fmt.Println(gs)
}
```

# 其他
## 分库分表
daog缺省支持分表，您需要为每个表指定分表函数，这需要设置TableMeta.ShardingFunc，具体需要在编译出的-ext.go文件的init函数中设置。
daog缺省支持分库，分库策略需要您实现DatasourceShardingPolicy接口，并在NewShardingDatasource是传入。ModInt64ShardingDatasourcePolicy是一个简单实现。GetDatasourceShardingKeyFromCtx函数实现了从context.Context中读取datasource sharding key的能力。

## 日志输出
通过DbConf.LogSQL可以设置该数据源是否需要输出执行的sql及参数，可以为数据源指定，每个一个TransContext执行时会继承这个配置，您也可以设置
TransContext.LogSQL属性为每个事务上下文设置，更细粒度的控制日志输出。

日志通过调用 daog.GLogger,它的数据类型是 SQLLogger 接口， 输出实现，缺省是调用标准库的log包，您也可以自行实现SQLLogger接口，并构建对象赋值给 GLogger 全局变量：

GetTraceIdFromContext函数可以从context.Context中读取traceId
GetGoroutineIdFromContext 函数可以从 context.Context中读取创建TransContext的goroutine id
## 日期

golang的time.Time支持纳秒级别，但数据库支持秒级别即可，因此提供ttypes.NormalDate和ttypes.NormalDatetime来支持。
他们都内置了对json序列化的支持。序列化格式通过ttypes.DateFormat和ttypes.DatetimeFormat来设置，他们缺省是yyyy-MM-dd格式。

* NormalDate.ToTimePointer 方法可以返回 NormalDate 包含的*time.Time
* NormalDatetime.ToTimePointer 方法可以返回 NormalDatetime 包含的*time.Time

## null字段值
golang sql包支持NullString, NullTime, NullFloat64, Nullxxx类型，但这些类型没有实现json序列化、反序列化接口。daog仅仅支持NullString和NullTime,其他的不支持，
这是因为实际的业务中，大部分情况要求字段是非空，尤其是数字数据类型。daog封装了NilableDate、NilableDatetime、NilableString三种类型，并提供一些函数用于简化开发，同时提供json序列化支持。

* FromDatetime、FromDate、FromString函数用于把Time\string转换成Nilable对象；
* NilableDate{},NilableDatetime{},NilableString{}表示null对象
* NilableDate.ToTimePointer 方法可以返回 NilableDate 包含的*time.Time, 如果 NilableDate 包含nil，那返回nil
* NilableDatetime.ToTimePointer 方法可以返回 NilableDatetime 包含的*time.Time， 如果 NilableDatetime 包含nil，那返回nil
