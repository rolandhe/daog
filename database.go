// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

// Package daog, 是轻量级的数据库访问组件，它并不能称之为orm组件，仅仅提供了一组函数用以实现常用的数据库访问功能。
// 它是高性能的，与原生的使用sql包函数相比，没有性能损耗，这是因为，它并没有使用反射技术，而是使用编译技术把create table sql语句编译成daog需要的go代码。
// 它目前仅支持mysql。
//
// 设计思路来源于java的[orm框架sampleGenericDao](https://github.com/tiandarwin/simpleGenericDao)和protobuf的编译思路。之所以选择编译
// 而没有使用反射，是因为基于编译的抽象没有性能损耗。
package daog

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rolandhe/daog/utils"
	"log"
	"strings"
	"time"
)

var invalidShardingDatasourceKey error = errors.New("invalid shard key")

// DbConf 数据源配置, 包括数据库url和连接池相关配置，特别注意，它支持按数据源在日志中输出执行的sql
type DbConf struct {
	// 数据库url
	DbUrl string
	// 最大连接数
	Size int
	// 连接的最大生命周期，单位是秒
	Life int
	// 最大空闲连接数
	IdleCons int
	// 最大空闲时间，单位是秒
	IdleTime int
	// 该在该数据源上执行sql是是否需要把待执行的sql输出到日志
	LogSQL bool
}

// NewDatasource 按照配置创建单个数据源对象
func NewDatasource(conf *DbConf) (Datasource, error) {
	dbUrl := conf.DbUrl
	if -1 == strings.Index(conf.DbUrl, "interpolateParams") {
		if strings.Index(conf.DbUrl, "?") != -1 {

			dbUrl = dbUrl + "&interpolateParams=true"
		} else {
			dbUrl = dbUrl + "?interpolateParams=true"
		}
	}
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Printf("goid=%d, %v\n", utils.QuickGetGoroutineId(), err)
		return nil, err
	}
	if conf.Size > 0 {
		db.SetMaxOpenConns(conf.Size)
	}
	if conf.IdleCons > 0 {
		db.SetMaxIdleConns(conf.IdleCons)
	}
	if conf.IdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(int64(conf.IdleTime) * 1e9))
	}
	if conf.Life > 0 {
		db.SetConnMaxLifetime(time.Duration(int64(conf.Life) * 1e9))
	}

	return &singleDatasource{db, conf.LogSQL}, nil
}

// NewShardingDatasource 创建多分片数据源,创建好的数据源是复合数据源，内含confs参数指定的多个数据源，也包含一个分片策略，
// 使用 NewShardingDatasource 数据源时要求使用 NewTransContextWithSharding 来创建事务上下文
func NewShardingDatasource(confs []*DbConf, policy DatasourceShardingPolicy) (Datasource, error) {
	var dbs []Datasource
	for _, conf := range confs {
		ds, err := NewDatasource(conf)
		if err != nil {
			for _, sds := range dbs {
				sds.Shutdown()
			}
			return nil, err
		}
		dbs = append(dbs, ds)
	}
	if len(dbs) == 0 {
		return nil, errors.New("no db confs")
	}
	return &shardingDatasource{dbs, policy}, nil
}

// Datasource 描述一个数据源，确切的说是一个数据源分片，它对应一个mysql database
type Datasource interface {
	getDB(ctx context.Context) *sql.DB
	// Shutdown 关闭数据源
	Shutdown()
	// IsLogSQL 本数据源是否需要输出执行的sql到日志
	IsLogSQL() bool
}

// DatasourceShardingPolicy 数据源分片策略
type DatasourceShardingPolicy interface {
	// Shard 根据分片key和分片总数来路由分片数据源
	Shard(shardKey any, count int) (int, error)
}

// ModInt64ShardingDatasourcePolicy 分片key是int64直接对分片总数取模路由策略，这是最简单的方式
type ModInt64ShardingDatasourcePolicy int64

func (h ModInt64ShardingDatasourcePolicy) Shard(shardKey any, count int) (int, error) {
	key, ok := shardKey.(int64)
	if !ok {
		return 0, invalidShardingDatasourceKey
	}
	return int(key % int64(count)), nil
}

type singleDatasource struct {
	db     *sql.DB
	logSQL bool
}

func (db *singleDatasource) getDB(ctx context.Context) *sql.DB {
	return db.db
}
func (db *singleDatasource) Shutdown() {
	db.db.Close()
}

func (db *singleDatasource) IsLogSQL() bool {
	return db.logSQL
}

type shardingDatasource struct {
	singleDatasource []Datasource
	policy           DatasourceShardingPolicy
}

func (db *shardingDatasource) getDB(ctx context.Context) *sql.DB {
	key := getDatasourceShardingKeyFromCtx(ctx)
	index, err := db.policy.Shard(key, len(db.singleDatasource))
	if err != nil {
		LogError(ctx, err)
		return nil
	}
	return db.singleDatasource[index].getDB(ctx)
}
func (db *shardingDatasource) Shutdown() {
	for _, sds := range db.singleDatasource {
		sds.Shutdown()
	}
}

func (db *shardingDatasource) IsLogSQL() bool {
	return db.singleDatasource[0].IsLogSQL()
}
