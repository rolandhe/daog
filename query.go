// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

import "time"

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
	defer func() {
		if err != nil {
			forError(tc)
		}
	}()
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

func QueryOneMatcher[T any](m Matcher, meta *TableMeta[T], tc *TransContext) (*T, error) {
	var err error
	err = tc.check()
	if err != nil {
		return nil, err
	}
	sql, args := selectQuery(meta, tc.ctx, m)
	defer func() {
		if err != nil {
			forError(tc)
		}
	}()
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

type RowFieldPointMapper[T any] interface {
	ExtractScanFieldPoints(ins *T) []any
}

// QueryRawSQL 执行原生select sql语句,返回行数据数组，行数据使用T struct描述
// mapper, 它T的各个field指针提取出来并按照顺序生成一个slice，用于Row.Scan方法，把sql字段映射到T对象的各个Field上
func QueryRawSQL[T any](tc *TransContext, mapper RowFieldPointMapper[T], sql string, args ...any) ([]*T, error) {
	var err error
	err = tc.check()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			forError(tc)
		}
	}()
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
		scanFields := mapper.ExtractScanFieldPoints(ins)
		if err = rows.Scan(scanFields...); err != nil {
			return nil, err
		}
		inses = append(inses, ins)
	}

	return inses, nil
}
