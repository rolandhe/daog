// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

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

type DbConf struct {
	DbUrl    string
	Size     int
	Life     int
	IdleCons int
	IdleTime int
	LogSQL   bool
}

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
		log.Printf("goid=%d, %v\n", utils.QuickGetGoRoutineId(), err)
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

type Datasource interface {
	getDB(ctx context.Context) *sql.DB
	Shutdown()
	IsLogSQL() bool
}

type DatasourceShardingPolicy interface {
	Shard(shardKey any, count int) (int, error)
}

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
	key := GetDatasourceShardingKeyFromCtx(ctx)
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
