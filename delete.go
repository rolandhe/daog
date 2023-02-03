// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

func DeleteById[T any](tc *TransContext,id int64, meta *TableMeta[T]) (int64, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return DeleteByMatcher(tc,m, meta)
}

func DeleteByIds[T any](tc *TransContext,ids []int64, meta *TableMeta[T]) (int64, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, ConvertToAnySlice(ids))
	return DeleteByMatcher(tc,m, meta)
}

func DeleteByMatcher[T any](tc *TransContext,matcher Matcher, meta *TableMeta[T]) (int64, error) {
	base := "delete from " + GetTableName(tc.ctx, meta)
	if matcher == nil {
		LogInfo(tc.ctx, "delete must has condition")
		return 0, nil
	}
	var args []any
	condi, args := matcher.ToSQL(args)
	if condi == "" {
		LogInfo(tc.ctx, "delete must has condition")
		return 0, nil
	}

	sql := base + " where " + condi

	return execSQLCore(tc, sql, args)
}
