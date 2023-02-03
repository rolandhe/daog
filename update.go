// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

import (
	txrequest "github.com/rolandhe/daog/tx"
	"time"
)

func Update[T any](tc *TransContext,ins *T, meta *TableMeta[T]) (int64, error) {
	idValue := meta.LookupFieldFunc(TableIdColumnName, ins, false)
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, idValue)
	sql, args := updateExec(meta, ins, tc.ctx, m)

	return execSQLCore(tc, sql, args)
}

func UpdateList[T any](tc *TransContext,insList []*T, meta *TableMeta[T]) (int64, error) {
	var affectRow int64

	for _, ins := range insList {
		n, err := Update(tc,ins, meta)
		if err != nil {
			if tc.txRequest == txrequest.RequestNone {
				return affectRow, err
			}
			return 0, err
		}
		affectRow += n
	}
	return affectRow, nil
}

func UpdateById[T any](tc *TransContext,modifier Modifier, id int64, meta *TableMeta[T]) (int64, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return UpdateByModifier(tc,modifier, m, meta)
}

func UpdateByIds[T any](tc *TransContext,modifier Modifier, ids []int64, meta *TableMeta[T]) (int64, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, ConvertToAnySlice(ids))
	return UpdateByModifier(tc,modifier, m, meta)
}

func UpdateByModifier[T any](tc *TransContext,modifier Modifier, matcher Matcher, meta *TableMeta[T]) (int64, error) {
	sql, args := buildModifierExec(meta, tc.ctx, modifier, matcher)
	if sql == "" {
		return 0, nil
	}
	return execSQLCore(tc, sql, args)
}

func ExecRawSQL(tc *TransContext, sql string, args ...any) (int64, error) {
	return execSQLCore(tc, sql, args)
}

func execSQLCore(tc *TransContext, sql string, args []any) (int64, error) {
	err := tc.check()
	if err != nil {
		return 0, err
	}
	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	result, err := tc.conn.ExecContext(tc.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	affectRow, err := result.RowsAffected()
	return affectRow, err
}
