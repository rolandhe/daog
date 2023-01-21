# 轻量、高性能，仅支持mysql
daog是轻量级的数据库访问组件，它并不能称之为orm组件，提供了一组函数用以实现常用的数据库访问功能。
它是高性能的，与原生的使用sql包函数相比，没有性能损耗，这是因为，它并没有使用反射技术，而是使用编译技术把create table sql语句编译成daog需要的go代码。

## 编译组件

[complix](https://github.com/rolandhe/compilex) 是编译create table语句文件的工具，使用如下语句可以编译：

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
package entities

import (
	"github.com/roland/daog"
	dbtime "github.com/roland/daog/time"
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

### 写表

```
func create() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	amount, err := decimal.NewFromString("128.0")
	if err != nil {
		fmt.Println(err)
		return
	}
	t := &entities.GroupInfo{
		Name:        "roland",
		MainData:    `{"a":102}`,
		CreateAt:    dbtime.NormalDatetime(time.Now()),
		TotalAmount: amount,
	}
	affect, err := daog.Insert(t, entities.GroupInfoMeta, tc)
	fmt.Println(affect, t.Id, err)

	t.Name = "roland he"
	af, err := daog.Update(t, entities.GroupInfoMeta, tc)
	fmt.Println(af, err)
}
```

### 读取数据

```
func queryByIds() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	gs, err := daog.GetByIds([]int64{1, 2}, entities.GroupInfoMeta, tc)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByIds", string(j))
	fmt.Println(gs)
}
```

### 根据Matcher读取
```
func queryByMatcher() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	matcher := daog.NewMatcher().Eq(entities.GroupInfoFields.Name, "xiufeg").Lt(entities.GroupInfoFields.Id, 3)
	gs, err := daog.QueryListMatcher(matcher, entities.GroupInfoMeta, tc)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByMatcher", string(j))
	fmt.Println(gs)
}
```

### 先读再写

```
func update() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-100099")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	g, err := daog.GetById(1, entities.GroupInfoMeta, tc)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(g)
	fmt.Println("query", string(j))

	g.Name = "Eric"
	af, err := daog.Update(g, entities.GroupInfoMeta, tc)
	fmt.Println(af, err)
}
```

# 其他
## 分库分表
daog缺省支持分表，您需要为每个表指定分表函数，这需要设置TableMeta.ShardingFunc，具体需要在编译出的-ext.go文件的init函数中设置。
daog缺省支持分库，分库策略需要您实现DatasourceShardingPolicy接口，并在NewShardingDatasource是传入。ModInt64ShardingDatasourcePolicy是一个简单实现。GetDatasourceShardingKeyFromCtx函数实现了从context.Context中读取datasource sharding key的能力。

## 日志输出
通过DbConf.LogSQL可以设置该数据源是否需要输出执行的sql及参数，可以为数据源指定，每个一个TransContext执行时会继承这个配置，您也可以设置
TransContext.LogSQL属性为每个事务上下文设置，更细粒度的控制日志输出。

日志的输出实现，缺省是调用标准库的log包，您也可以通过配置daog包的3个全局函数来修改：
* DaogLogErrorFunc
* DaogLogInfoFunc
* DaogLogExecSQLFunc

## 日期

golang的time.Time支持纳秒级别，但数据库支持秒级别即可，因此提供dbtime.NormalDate和dbtime.NormalDatetime来支持。
他们都内置了对json序列化的支持。序列化格式通过dbtime.DateFormat和dbtime.DatetimeFormat来设置，他们缺省是yyyy-MM-dd格式。