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
)

type tcStatus int

const (
	TraceId               = "trace-id"
	CtxValues             = "values"
	ShardingKey           = "shardingKey"
	DatasourceShardingKey = "datasourceSharingKey"

	tcStatusInit      = tcStatus(0)
	tcStatusCommitted = tcStatus(1)
	tcStatusRollback  = tcStatus(2)
	tcStatusFailed    = tcStatus(3)
	tcStatusInvalid   = tcStatus(4)
)

var invalidTcStatus = errors.New("invalid tc status")

func NewTransContext(datasource Datasource, txRequest txrequest.RequestStyle, traceId string) (*TransContext, error) {
	var conn *sql.Conn
	var err error
	ctx := buildContext(traceId, nil, nil)

	if conn, err = datasource.getDB(ctx).Conn(context.Background()); err != nil {
		return nil, err
	}
	tc := &TransContext{
		txRequest: txRequest,
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
	ctx := buildContext(traceId, ShardingKey, nil)
	if conn, err = datasource.getDB(ctx).Conn(context.Background()); err != nil {
		return nil, err
	}
	tc := &TransContext{
		txRequest: txRequest,
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
	if tc.txRequest == txrequest.RequestNone || tc.tx != nil {
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

func (tc *TransContext) commit() error {
	if tc.txRequest == txrequest.RequestNone || tc.tx == nil {
		return nil
	}
	if tc.status != tcStatusInit {
		return errors.New(fmt.Sprintf("tc status error,%d", tc.status))
	}
	err := tc.tx.Commit()
	if err == nil {
		tc.status = tcStatusCommitted
	}
	return err
}

func (tc *TransContext) rollback() error {
	if tc.txRequest == txrequest.RequestNone || tc.tx == nil {
		return nil
	}
	if tc.status != tcStatusInit {
		return errors.New(fmt.Sprintf("tc status error,%d", tc.status))
	}
	err := tc.tx.Rollback()
	if err != nil {
		tc.status = tcStatusFailed
	} else {
		tc.status = tcStatusRollback
	}
	return err
}

func (tc *TransContext) Complete(e error) {
	LogError(tc.ctx, e)
	if tc.status == tcStatusInvalid {
		return
	}
	if tc.txRequest == txrequest.RequestNone {
		if err := tc.conn.Close(); err != nil {
			LogError(tc.ctx, err)
		}
		tc.status = tcStatusInvalid
		return
	}
	if tc.status == tcStatusInit {
		if e != nil {
			tc.rollbackAndClose()
			tc.status = tcStatusInvalid
			return
		}
		err := tc.commit()
		if err != nil {
			LogError(tc.ctx, err)
			tc.rollbackAndClose()
		}
		tc.status = tcStatusInvalid
	}

}

func (tc *TransContext) rollbackAndClose() {
	var err error
	if err = tc.rollback(); err != nil {
		LogError(tc.ctx, err)
		if err := tc.conn.Close(); err != nil {
			LogError(tc.ctx, err)
		}
	}
}
func buildContext(traceId string, shardingKey any, dataSourceSharingKey any) context.Context {
	mp := map[string]any{}
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
