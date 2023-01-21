package daog

import (
	txrequest "github.com/roland/daog/tx"
)

func Update[T any](ins *T, meta *TableMeta[T], tc *TransContext) (int64, error) {
	idValue := meta.LookupFieldFunc("id", ins, false)
	m := NewMatcher()
	fieldId := "id"
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, idValue)
	sql, args := updateExec(meta, ins, tc.ctx, m)

	return execSQLCore(tc, sql, args)
}

func UpdateList[T any](insList []*T, meta *TableMeta[T], tc *TransContext) (int64, error) {
	var affectRow int64

	for _, ins := range insList {
		n, err := Update(ins, meta, tc)
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

func UpdateById[T any](modifier Modifier, id int64, meta *TableMeta[T], tc *TransContext) (int64, error) {
	m := NewMatcher()
	fieldId := "id"
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return UpdateByModifier(modifier, m, meta, tc)
}

func UpdateByIds[T any](modifier Modifier, ids []int64, meta *TableMeta[T], tc *TransContext) (int64, error) {
	m := NewMatcher()
	fieldId := "id"
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, ConvertToAnySlice(ids))
	return UpdateByModifier(modifier, m, meta, tc)
}

func UpdateByModifier[T any](modifier Modifier, matcher Matcher, meta *TableMeta[T], tc *TransContext) (int64, error) {
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
	var err error
	defer func() {
		if err != nil {
			forError(tc)
		}
	}()

	if tc.LogSQL {
		DaogLogExecSQL(tc.ctx, sql, args)
	}
	result, err := tc.conn.ExecContext(tc.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	affectRow, err := result.RowsAffected()
	return affectRow, err
}
