package daog

import (
	"fmt"
	"strings"
)

func Insert[T any](ins *T, meta *TableMeta[T], tc *TransContext) (int64, error) {
	tableName := GetTableName(tc.ctx, meta)
	var icolums []string
	var holder []string
	var exclude map[string]int
	if meta.AutoColumn == "" {
		icolums = meta.Columns
		for range icolums {
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
			icolums = append(icolums, column)
			holder = append(holder, "?")
		}
	}

	sql := fmt.Sprintf("insert into %s(%s) values(%s)", tableName, strings.Join(icolums, ","), strings.Join(holder, ","))
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
	var err error
	defer func() {
		if err != nil {
			forError(tc)
		}
	}()

	if tc.LogSQL {
		DaogLogExecSQL(tc, sql, args)
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