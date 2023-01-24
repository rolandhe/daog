// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

func DeleteById[T any](id int64, meta *TableMeta[T], tc *TransContext) (int64, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return DeleteByMatcher(m, meta, tc)
}

func DeleteByIds[T any](ids []int64, meta *TableMeta[T], tc *TransContext) (int64, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, ConvertToAnySlice(ids))
	return DeleteByMatcher(m, meta, tc)
}

func DeleteByMatcher[T any](matcher Matcher, meta *TableMeta[T], tc *TransContext) (int64, error) {
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
