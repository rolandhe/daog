// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	txrequest "github.com/rolandhe/daog/tx"
	"github.com/rolandhe/daog/utils"
)

type tcStatus int

const (
	TraceId               = "trace-id"
	GoId                  = "goroutine-id"
	CtxValues             = "values"
	ShardingKey           = "shardingKey"
	DatasourceShardingKey = "datasourceSharingKey"

	tcStatusInit    = tcStatus(1)
	tcStatusInvalid = tcStatus(4)
)

var invalidTcStatus = errors.New("invalid tc status")

var metRecover = errors.New("met recover")

func NewTransContext(datasource Datasource, txRequest txrequest.RequestStyle, traceId string) (*TransContext, error) {
	var conn *sql.Conn
	var err error
	goroutineId := utils.QuickGetGoRoutineId()
	ctx := buildContext(goroutineId, traceId, nil, nil)

	if conn, err = datasource.getDB(ctx).Conn(context.Background()); err != nil {
		return nil, err
	}
	tc := &TransContext{
		txRequest: txRequest,
		status:    tcStatusInit,
		ctx:       ctx,
		conn:      conn,
		LogSQL:    datasource.IsLogSQL(),
	}
	err = tc.begin()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return tc, nil
}

func NewTransContextWithSharding(datasource Datasource, txRequest txrequest.RequestStyle, traceId string, shardingKey any, datasourceShardingKey any) (*TransContext, error) {
	var conn *sql.Conn
	var err error
	goroutineId := utils.QuickGetGoRoutineId()
	ctx := buildContext(goroutineId, traceId, ShardingKey, nil)
	if conn, err = datasource.getDB(ctx).Conn(context.Background()); err != nil {
		return nil, err
	}
	tc := &TransContext{
		txRequest: txRequest,
		status:    tcStatusInit,
		ctx:       ctx,
		conn:      conn,
		LogSQL:    datasource.IsLogSQL(),
	}
	err = tc.begin()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return tc, nil
}

func WrapTrans(tc *TransContext, workFn func(tc *TransContext) error) error {
	var err error
	defer func() {
		tc.CompleteWithPanic(err, recover())
	}()
	err = workFn(tc)
	return err
}

func WrapTransWithResult[T any](tc *TransContext, workFn func(tc *TransContext) (T, error)) (T, error) {
	var err error
	defer func() {
		tc.CompleteWithPanic(err, recover())
	}()
	ret, err := workFn(tc)
	return ret, err
}

func GetDatasourceShardingKeyFromCtx(ctx context.Context) any {
	mapAny := ctx.Value(CtxValues)
	if mapAny == nil {
		return nil
	}
	mapValue, ok := mapAny.(map[string]any)
	if !ok {
		return nil
	}
	return mapValue[DatasourceShardingKey]
}

func GetTraceIdFromContext(ctx context.Context) string {
	values := ctx.Value(CtxValues)
	if values == nil {
		return ""
	}

	v, ok := values.(map[string]any)
	if !ok {
		return ""
	}
	data, ok := v[TraceId]
	if !ok {
		return ""
	}
	trace, ok := data.(string)
	if !ok {
		return ""
	}
	return trace
}

func GetGoRoutineIdFromContext(ctx context.Context) uint64 {
	values := ctx.Value(CtxValues)
	if values == nil {
		return 0
	}

	v, ok := values.(map[string]any)
	if !ok {
		return 0
	}
	data, ok := v[GoId]
	if !ok {
		return 0
	}
	goid, ok := data.(uint64)
	if !ok {
		return 0
	}
	return goid
}

func GetTableShardingKeyFromCtx(ctx context.Context) any {
	mapAny := ctx.Value(CtxValues)
	if mapAny == nil {
		return nil
	}
	mapValue, ok := mapAny.(map[string]any)
	if !ok {
		return nil
	}
	return mapValue[ShardingKey]
}

type TransContext struct {
	txRequest txrequest.RequestStyle
	tx        driver.Tx
	conn      *sql.Conn
	status    tcStatus
	ctx       context.Context
	LogSQL    bool
}

func (tc *TransContext) begin() (err error) {
	if tc.txRequest == txrequest.RequestNone {
		return nil
	}
	tc.tx, err = tc.conn.BeginTx(context.Background(), &sql.TxOptions{
		ReadOnly: tc.txRequest == txrequest.RequestReadonly,
	})
	return err
}

func (tc *TransContext) check() error {
	if tc.status != tcStatusInit {
		return invalidTcStatus
	}
	return nil
}

func (tc *TransContext) commitAndReleaseConn() error {
	if tc.txRequest == txrequest.RequestNone {
		return nil
	}
	if tc.status != tcStatusInit {
		return errors.New(fmt.Sprintf("tc status error,%d", tc.status))
	}
	err := tc.tx.Commit()
	if err == nil {
		closeConn(tc)
	}
	return err
}

func (tc *TransContext) rollbackAndReleaseConn() error {
	if tc.txRequest == txrequest.RequestNone {
		return nil
	}
	if tc.status != tcStatusInit {
		return errors.New(fmt.Sprintf("tc status error,%d", tc.status))
	}

	err := tc.tx.Rollback()
	if err == nil {
		closeConn(tc)
	}
	return err
}

func closeConn(tc *TransContext) {
	if err := tc.conn.Close(); err != nil {
		LogError(tc.ctx, err)
	}
}

func (tc *TransContext) CompleteWithPanic(e error, fetal any) {
	if fetal != nil {
		tc.Complete(metRecover)
		panic(fetal)
	}
	tc.Complete(e)
}

func (tc *TransContext) Complete(e error) {
	LogError(tc.ctx, e)
	if tc.status == tcStatusInvalid {
		return
	}
	if tc.txRequest == txrequest.RequestNone {
		closeConn(tc)
		tc.status = tcStatusInvalid
		return
	}
	if tc.status == tcStatusInit {
		if e != nil {
			tc.rollbackAndReleaseConn()
		} else {
			tc.commitAndReleaseConn()
		}
		tc.status = tcStatusInvalid
	}
}

func buildContext(goroutineId uint64, traceId string, shardingKey any, dataSourceSharingKey any) context.Context {
	mp := map[string]any{}
	mp[GoId] = goroutineId
	mp[TraceId] = traceId
	if shardingKey != nil {
		mp[ShardingKey] = shardingKey
	}
	if dataSourceSharingKey != nil {
		mp[DatasourceShardingKey] = dataSourceSharingKey
	}

	ctx := context.WithValue(context.Background(), CtxValues, mp)
	return ctx
}
