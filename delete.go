// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

package daog

// DeleteById 根据主键id删除记录
//
// 参数: id 主键 , meta 表的元数据，由compile编译生成，比如  GroupInfo.GroupInfoMeta
//
// 返回值: 删除记录数，是否出错
func DeleteById[T any](tc *TransContext, id int64, meta *TableMeta[T]) (int64, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return DeleteByMatcher(tc, m, meta)
}

// DeleteByIds 根据主键id删除记录
//
// 参数: ids 一批主键 , meta 表的元数据，由compile编译生成，比如  GroupInfo.GroupInfoMeta
//
// 返回值: 删除记录数及是否出错
func DeleteByIds[T any](tc *TransContext, ids []int64, meta *TableMeta[T]) (int64, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, ConvertToAnySlice(ids))
	return DeleteByMatcher(tc, m, meta)
}

// DeleteByMatcher 通过匹配条件删除数据，返回删除记录数及是否出错
func DeleteByMatcher[T any](tc *TransContext, matcher Matcher, meta *TableMeta[T]) (int64, error) {
	base := "delete from " + GetTableName(tc.ctx, meta)
	if matcher == nil {
		GLogger.Info(tc.ctx, "delete must has condition")
		return 0, nil
	}
	var args []any
	condi, args := matcher.ToSQL(args)
	if condi == "" {
		GLogger.Info(tc.ctx, "delete must has condition")
		return 0, nil
	}

	sql := base + " where " + condi

	return execSQLCore(tc, sql, args)
}
