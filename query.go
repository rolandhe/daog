// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

import (
	"errors"
	"time"
)

var invalidBatchSizeError = errors.New("page size must be greater than 0")

func GetAll[T any](tc *TransContext, meta *TableMeta[T], viewColumns ...string) ([]*T, error) {
	return QueryPageListMatcherWithViewColumns(tc, nil, meta, viewColumns, nil)
}

func GetById[T any](tc *TransContext, id int64, meta *TableMeta[T], viewColumns ...string) (*T, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return QueryOneMatcher(tc, m, meta, viewColumns...)
}

func GetByIds[T any](tc *TransContext, ids []int64, meta *TableMeta[T], viewColumns ...string) ([]*T, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.In(fieldId, ConvertToAnySlice(ids))
	return QueryPageListMatcherWithViewColumns(tc, m, meta, viewColumns, nil)
}

func QueryListMatcher[T any](tc *TransContext, m Matcher, meta *TableMeta[T], orders ...*Order) ([]*T, error) {
	return QueryPageListMatcher(tc, m, meta, nil, orders...)
}

func QueryPageListMatcher[T any](tc *TransContext, m Matcher, meta *TableMeta[T], pager *Pager, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcherWithViewColumns(tc, m, meta, nil, pager, orders...)
}

func QueryPageListMatcherWithViewColumns[T any](tc *TransContext, m Matcher, meta *TableMeta[T], viewColumns []string, pager *Pager, orders ...*Order) ([]*T, error) {
	sql, args := selectQuery(meta, tc.ctx, m, pager, orders, viewColumns)
	return queryRawSQLCore(tc, func() (*T, []any) {
		return buildInsInfoOfRow(meta, viewColumns)
	}, sql, args...)
}

type BatchHandler[T any] func(batch []*T) error

// QueryListMatcherByBatchHandle 读取数据并且分批处理数据，当读取的数据量巨大时非常有用，如果数据都读入内存，容易打爆内存，分批量处理就非常有用
// batchSize 每批处理数据的最大容量，必须大于0，但不要设置太大，当设置为1时，退化成每条处理
// handler 用于处理每批数据的函数
// 查询数据最大上限数， 0 表示无上限
func QueryListMatcherByBatchHandle[T any](tc *TransContext, m Matcher, meta *TableMeta[T], totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error {
	return QueryListMatcherWithViewColumnsByBatchHandle(tc, m, meta, nil, totalLimit, batchSize, handler, orders...)
}

func QueryListMatcherWithViewColumnsByBatchHandle[T any](tc *TransContext, m Matcher, meta *TableMeta[T], viewColumns []string, totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error {
	if batchSize <= 0 {
		return invalidBatchSizeError
	}
	var pager *Pager
	if totalLimit > 0 {
		pager = &Pager{0, totalLimit}
	}
	sql, args := selectQuery(meta, tc.ctx, m, pager, orders, viewColumns)

	return queryRawSQLByBatchHandleCore(tc, batchSize, handler, func() (*T, []any) {
		return buildInsInfoOfRow(meta, viewColumns)
	}, sql, args...)
}

func QueryOneMatcher[T any](tc *TransContext, m Matcher, meta *TableMeta[T], viewColumns ...string) (*T, error) {
	var err error
	err = tc.check()
	if err != nil {
		return nil, err
	}
	sql, args := selectQuery(meta, tc.ctx, m, nil, nil, viewColumns)

	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	rows, err := tc.conn.QueryContext(tc.ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	ins, scanFields := buildInsInfoOfRow(meta, viewColumns)
	if err = rows.Scan(scanFields...); err != nil {
		return nil, err
	}
	return ins, nil
}

type ExtractScanFieldPoints[T any] func(ins *T) []any

// QueryRawSQL 执行原生select sql语句,返回行数据数组，行数据使用T struct描述
// mapper, 它T的各个field指针提取出来并按照顺序生成一个slice，用于Row.Scan方法，把sql字段映射到T对象的各个Field上
func QueryRawSQL[T any](tc *TransContext, extract ExtractScanFieldPoints[T], sql string, args ...any) ([]*T, error) {
	return queryRawSQLCore(tc, func() (*T, []any) {
		ins := new(T)
		return ins, extract(ins)
	}, sql, args...)
}

func QueryRawSQLByBatchHandle[T any](tc *TransContext, batchSize int, handler BatchHandler[T], extract ExtractScanFieldPoints[T], sql string, args ...any) error {
	return queryRawSQLByBatchHandleCore(tc, batchSize, handler, func() (*T, []any) {
		ins := new(T)
		return ins, extract(ins)
	}, sql, args...)
}

func Count[T any](tc *TransContext, m Matcher, meta *TableMeta[T]) (int64, error) {
	var err error
	err = tc.check()
	if err != nil {
		return 0, err
	}
	sql, args := countQuery(meta, tc.ctx, m)

	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	rows, err := tc.conn.QueryContext(tc.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, nil
	}
	var countValue int64
	if err = rows.Scan(&countValue); err != nil {
		return 0, err
	}
	return countValue, nil
}

type rowInsCreate[T any] func() (*T, []any)

func queryRawSQLByBatchHandleCore[T any](tc *TransContext, batchSize int, handler BatchHandler[T], creatorFunc rowInsCreate[T], sql string, args ...any) error {
	var err error
	err = tc.check()
	if err != nil {
		return err
	}

	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	rows, err := tc.conn.QueryContext(tc.ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	batch := make([]*T, batchSize)
	index := 0

	for rows.Next() {
		ins, scanFields := creatorFunc()
		if err = rows.Scan(scanFields...); err != nil {
			return err
		}
		batch[index] = ins
		index++
		if index == batchSize {
			if err = handler(batch); err != nil {
				return err
			}
			index = 0
		}
	}
	if index == 0 {
		return nil
	}

	batch = batch[:index]
	if err = handler(batch); err != nil {
		return err
	}

	return nil
}

func queryRawSQLCore[T any](tc *TransContext, creatorFunc rowInsCreate[T], sql string, args ...any) ([]*T, error) {
	var err error
	err = tc.check()
	if err != nil {
		return nil, err
	}

	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	rows, err := tc.conn.QueryContext(tc.ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var inses []*T
	for rows.Next() {
		ins, scanFields := creatorFunc()
		if err = rows.Scan(scanFields...); err != nil {
			return nil, err
		}
		inses = append(inses, ins)
	}

	return inses, nil
}
