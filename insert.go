// A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

import (
	"strings"
	"time"
)

// Insert 插入一条数据到表中，如果表有自增id，那么生成的id赋值到ins对象中
//
// 参数: ins 表实体对象，对应表的struct由 compile生成，比如 GroupInfo, meta 表的元数据，由compile编译生成，比如  GroupInfo.GroupInfoMeta
//
// 返回值: 插入的记录数，是否出错
func Insert[T any](tc *TransContext, ins *T, meta *TableMeta[T]) (int64, error) {
	tableName := GetTableName(tc.ctx, meta)
	if BeforeInsertCallback != nil {
		if err := BeforeInsertCallback(tableName, ins); err != nil {
			return 0, err
		}
	}

	var insertColumns []string
	var holder []string
	var exclude map[string]int
	if meta.AutoColumn == "" {
		insertColumns = meta.Columns
		for range insertColumns {
			holder = append(holder, "?")
		}
	} else {
		exclude = map[string]int{
			meta.AutoColumn: 1,
		}
		for _, column := range meta.Columns {
			if column == meta.AutoColumn {
				continue
			}
			insertColumns = append(insertColumns, column)
			holder = append(holder, "?")
		}
	}

	if err := auoFillField(tc, ins, meta); err != nil {
		return 0, err
	}

	var builder strings.Builder
	builder.WriteString("insert into ")
	builder.WriteString(tableName)
	builder.WriteString("(")
	builder.WriteString(strings.Join(insertColumns, ","))
	builder.WriteString(") values(")
	builder.WriteString(strings.Join(holder, ","))
	builder.WriteString(")")
	sql := builder.String()
	args := meta.ExtractFieldValues(ins, false, exclude)
	affect, lastId, err := execInsert(tc, sql, args, meta.AutoColumn != "")
	if err != nil {
		return 0, err
	}

	if meta.AutoColumn != "" {
		autoAddr := meta.LookupFieldFunc(meta.AutoColumn, ins, true)
		*(autoAddr.(*int64)) = lastId
	}

	return affect, err
}

func execInsert(tc *TransContext, sql string, args []any, auto bool) (int64, int64, error) {
	err := tc.check()
	if err != nil {
		return 0, 0, err
	}
	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	result, err := tc.conn.ExecContext(tc.ctx, sql, args...)
	if err != nil {
		return 0, 0, err
	}

	affectRow, err := result.RowsAffected()
	if !auto || err != nil {
		return affectRow, 0, err
	}

	id, err := result.LastInsertId()
	return affectRow, id, err
}

func auoFillField[T any](tc *TransContext, ins *T, meta *TableMeta[T]) error {
	if ChangeFieldOfInsBeforeWrite != nil {
		err := ChangeFieldOfInsBeforeWrite(tc.ExtInfo, &fieldExtractor[T]{
			ins:  ins,
			meta: meta,
		})
		return err
	}
	return nil
}
