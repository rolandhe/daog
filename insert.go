// A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

import (
	"fmt"
	"strings"
	"time"
)

func Insert[T any](tc *TransContext, ins *T, meta *TableMeta[T]) (int64, error) {
	tableName := GetTableName(tc.ctx, meta)
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

	sql := fmt.Sprintf("insert into %s(%s) values(%s)", tableName, strings.Join(insertColumns, ","), strings.Join(holder, ","))
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
