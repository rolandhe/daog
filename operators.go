package daog

import (
	"context"
	"fmt"
	txrequest "github.com/roland/daog/tx"
	"strings"
)

func GetTableName[T any](ctx context.Context, meta *TableMeta[T]) string {
	tableName := meta.Table

	if meta.ShardingFunc != nil {
		shardingKey := getTableShardingKeyFromCtx(ctx)
		tableName = meta.ShardingFunc(tableName, shardingKey)
	}
	return tableName
}

func buildSelectBase[T any](meta *TableMeta[T], ctx context.Context) string {
	sfmt := "select %s from %s"

	columnsStr := strings.Join(meta.Columns, ",")

	return fmt.Sprintf(sfmt, columnsStr, GetTableName(ctx, meta))
}

func selectQuery[T any](meta *TableMeta[T], ctx context.Context, matcher Matcher) (string, []any) {
	base := buildSelectBase(meta, ctx)
	if matcher == nil {
		return base, nil
	}
	if matcher == nil {
		return base, nil
	}
	var args []any
	condi, args := matcher.ToSQL(args)
	if condi == "" {
		return base, nil
	}

	return base + " where " + condi, args
}

func buildUpdateBase[T any](meta *TableMeta[T], ctx context.Context) string {
	sfmt := "update %s set %s"

	var upConds []string
	for _, v := range meta.Columns {
		if v == meta.AutoColumn {
			continue
		}
		upConds = append(upConds, v+" = ?")
	}
	upCondStmt := strings.Join(upConds, ",")

	return fmt.Sprintf(sfmt, GetTableName(ctx, meta), upCondStmt)
}

func updateExec[T any](meta *TableMeta[T], ins *T, ctx context.Context, matcher Matcher) (string, []any) {
	base := buildUpdateBase(meta, ctx)
	if matcher == nil {
		return base, nil
	}
	var exclude map[string]int
	if meta.AutoColumn != "" {
		exclude = map[string]int{
			meta.AutoColumn: 1,
		}
	}
	args := meta.ExtractFieldValues(ins, false, exclude)
	if matcher == nil {
		return base, args
	}
	condi, args := matcher.ToSQL(args)
	if condi == "" {
		return base, args
	}

	return base + " where " + condi, args
}

func buildModifierExec[T any](meta *TableMeta[T], ctx context.Context, modifier Modifier, matcher Matcher) (string, []any) {
	tableName := GetTableName(ctx, meta)
	base, args := modifier.toSQL(tableName)
	if base == "" {
		return "", nil
	}

	if matcher == nil {
		return base, args
	}

	condi, args := matcher.ToSQL(args)
	if condi == "" {
		return base, args
	}
	return base + " where " + condi, args
}

func buildInsInfoOfRow[T any](meta *TableMeta[T]) (*T, []any) {
	ins := meta.InstanceFunc()
	scanFields := meta.ExtractFieldValues(ins, true, nil)
	return ins, scanFields
}

func forError(tc *TransContext) {
	DaogLogInfo(tc.ctx, "met for Error")
	if tc.txRequest == txrequest.RequestNone {
		return
	}
	tc.rollback()
}

func ConvertToAnySlice[T any](data []T) []any {
	l := len(data)
	if l == 0 {
		return nil
	}
	target := make([]any, l)
	for i, v := range data {
		target[i] = v
	}
	return target
}
