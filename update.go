// A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

import (
	txrequest "github.com/rolandhe/daog/tx"
	"time"
)

type fieldExtractor[T any] struct {
	ins *T
	meta *TableMeta[T]
}

func (fe * fieldExtractor[T]) Extract(fieldName string) any{
	return fe.meta.LookupFieldFunc(fieldName,fe.ins,true)
}

// Update 更新一条数据，把 *T类型的 ins 更新到数据，ins中的主键必须被设置
// meta 表的元数据，由compile编译生成，比如  GroupInfo.GroupInfoMeta
// 返回值是 更新的数据的条数，是0或者1
func Update[T any](tc *TransContext, ins *T, meta *TableMeta[T]) (int64, error) {
	if BeforeUpdateCallback != nil{
		if err := BeforeUpdateCallback(ins);err != nil{
			return 0, err
		}
	}
	idValue := meta.LookupFieldFunc(TableIdColumnName, ins, false)
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, idValue)

	if err := auoFillField(tc,ins,meta);err != nil{
		return 0, err
	}

	sql, args, err := updateExec(meta, ins, tc.ctx, m)
	if err != nil {
		return 0, err
	}

	return execSQLCore(tc, sql, args)
}

// UpdateList 更新多条数据，把多个 *T类型的 ins 更新到数据，每个ins中的主键必须被设置
// meta 表的元数据，由compile编译生成，比如  GroupInfo.GroupInfoMeta
// 返回值是 更新的数据的条数，是0或者1
// 注意： 当 tc 的事务类型是 txrequest.RequestNone 时，如果某一个 ins 更新失败，会立即返回错误，但该  ins之前的更新都会成功，此时的两个返回值都不是0值
func UpdateList[T any](tc *TransContext, insList []*T, meta *TableMeta[T]) (int64, error) {
	var affectRow int64

	for _, ins := range insList {
		n, err := Update(tc, ins, meta)
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

// UpdateById 根据主键修改一条记录，需要修改的字段值通过 Modifier 指定
func UpdateById[T any](tc *TransContext, modifier Modifier, id int64, meta *TableMeta[T]) (int64, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return UpdateByModifier(tc, modifier, m, meta)
}

// UpdateByIds 根据多个主键修改多条记录，需要修改的字段值通过 Modifier 指定，表达 update table set a=?,b=? where id in(xx,xx)的语义
func UpdateByIds[T any](tc *TransContext, modifier Modifier, ids []int64, meta *TableMeta[T]) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.In(fieldId, ConvertToAnySlice(ids))
	return UpdateByModifier(tc, modifier, m, meta)
}


// UpdateByModifier 根据Matcher条件修改多条记录，需要修改的字段值通过 Modifier 指定，表达 update table set a=?,b=? where uid=? and status=0 的类似语义
func UpdateByModifier[T any](tc *TransContext, modifier Modifier, matcher Matcher, meta *TableMeta[T]) (int64, error) {
	if AddNewModifyFieldBeforeUpdate != nil{
		if err := AddNewModifyFieldBeforeUpdate(tc.ExtInfo,modifier, func(fieldName string) bool {
			for _,name := range meta.Columns {
				if name == fieldName {
					return true
				}
			}
			return false
		});err !=nil{
			return 0, err
		}
	}
	sql, args, err := buildModifierExec(meta, tc.ctx, modifier, matcher)
	if err != nil {
		return 0, err
	}
	if sql == "" {
		return 0, nil
	}
	return execSQLCore(tc, sql, args)
}

// ExecRawSQL 执行原生的sql
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
