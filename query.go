// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

import (
	"errors"
	"time"
)

var invalidPageSizeError = errors.New("page size must be greater than 0")

func GetById[T any](id int64, meta *TableMeta[T], tc *TransContext) (*T, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return QueryOneMatcher(m, meta, tc)
}

func GetByIds[T any](ids []int64, meta *TableMeta[T], tc *TransContext) ([]*T, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.In(fieldId, ConvertToAnySlice(ids))
	return QueryListMatcher(m, meta, tc)
}

func QueryListMatcher[T any](m Matcher, meta *TableMeta[T], tc *TransContext) ([]*T, error) {
	var err error
	err = tc.check()
	if err != nil {
		return nil, err
	}

	sql, args := selectQuery(meta, tc.ctx, m)
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
		ins, scanFields := buildInsInfoOfRow(meta)
		if err = rows.Scan(scanFields...); err != nil {
			return nil, err
		}
		inses = append(inses, ins)
	}

	return inses, nil
}

type PageHandler[T any] func(page []*T) error

// QueryListMatcherPageHandle 读取数据并且分页处理数据，当读取的数据量巨大时非常有用，如果数据都读入内存，容易打爆内存，分批量处理就非常有用
// pageSize 每批处理数据的最大容量，必须大于0，但不要设置太大，当设置为1时，退化成每条处理
// handler 用于处理每批数据的函数
func QueryListMatcherPageHandle[T any](m Matcher, meta *TableMeta[T],pageSize int,handler PageHandler[T],tc *TransContext) error {
	var err error
	if pageSize <= 0 {
		return invalidPageSizeError
	}
	err = tc.check()
	if err != nil {
		return err
	}

	sql, args := selectQuery(meta, tc.ctx, m)
	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	rows, err := tc.conn.QueryContext(tc.ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	page := make([]*T,pageSize)
	index := 0

	for rows.Next() {
		ins, scanFields := buildInsInfoOfRow(meta)
		if err = rows.Scan(scanFields...); err != nil {
			return  err
		}
		page[index] = ins
		index++
		if index == pageSize {
			if err = handler(page);err != nil{
				return err
			}
			index = 0
		}
	}

	if index == 0 {
		return nil
	}

	page = page[:index]
	if err = handler(page);err != nil{
		return err
	}

	return nil
}

func QueryOneMatcher[T any](m Matcher, meta *TableMeta[T], tc *TransContext) (*T, error) {
	var err error
	err = tc.check()
	if err != nil {
		return nil, err
	}
	sql, args := selectQuery(meta, tc.ctx, m)

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
	ins, scanFields := buildInsInfoOfRow(meta)
	if err = rows.Scan(scanFields...); err != nil {
		return nil, err
	}
	return ins, nil
}


type ExtractScanFieldPoints[T any] func(ins *T) []any

// QueryRawSQL 执行原生select sql语句,返回行数据数组，行数据使用T struct描述
// mapper, 它T的各个field指针提取出来并按照顺序生成一个slice，用于Row.Scan方法，把sql字段映射到T对象的各个Field上
func QueryRawSQL[T any](tc *TransContext, extract ExtractScanFieldPoints[T], sql string, args ...any) ([]*T, error) {
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
		ins := new(T)
		scanFields := extract(ins)
		if err = rows.Scan(scanFields...); err != nil {
			return nil, err
		}
		inses = append(inses, ins)
	}

	return inses, nil
}
