// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TableIdColumnName = "id"

func GetTableName[T any](ctx context.Context, meta *TableMeta[T]) string {
	tableName := meta.Table

	if meta.ShardingFunc != nil {
		shardingKey := GetTableShardingKeyFromCtx(ctx)
		tableName = meta.ShardingFunc(tableName, shardingKey)
	}
	return tableName
}

func buildSelectBase[T any](meta *TableMeta[T], viewColumns []string, ctx context.Context) string {
	columnsStr := ""
	if len(viewColumns) == 0 {
		columnsStr = strings.Join(meta.Columns, ",")
	} else {
		columnsStr = strings.Join(viewColumns, ",")
	}
	return "select " + columnsStr + " from " + GetTableName(ctx, meta)
}

func selectQuery[T any](meta *TableMeta[T], ctx context.Context, matcher Matcher, pager *Pager, orders []*Order, viewColumns []string) (string, []any) {
	base := buildSelectBase(meta, viewColumns, ctx)
	if matcher == nil {
		return base, nil
	}
	var args []any
	condi, args := matcher.ToSQL(args)

	if condi == "" {
		return base + buildQuerySuffix(pager, orders), nil
	}

	return base + " where " + condi + buildQuerySuffix(pager, orders), args
}

func countQuery[T any](meta *TableMeta[T], ctx context.Context, matcher Matcher) (string, []any) {
	var base string
	if meta.AutoColumn == "" {
		base = "SELECT COUNT(*) FROM " + GetTableName(ctx, meta)
	} else {
		base = "SELECT COUNT(" + meta.AutoColumn + ") FROM " + GetTableName(ctx, meta)
	}

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

func buildQuerySuffix(pager *Pager, orders []*Order) string {
	ordStat := ""
	last := len(orders) - 1

	for i, order := range orders {
		if i == 0 {
			ordStat = " order by " + order.ColumnName
		} else {
			ordStat = ordStat + order.ColumnName
		}
		if order.Desc {
			ordStat = ordStat + " desc"
		}
		if i < last {
			ordStat = ordStat + ","
		}
	}
	if pager == nil {
		return ordStat
	}
	limitStat := ""

	if pager.PageNumber == 1 {
		limitStat = " limit " + strconv.Itoa(pager.PageSize)
	} else {
		startPos := int64(pager.PageNumber-1) * int64(pager.PageSize)
		limitStat = " limit " + strconv.FormatInt(startPos, 10) + "," + strconv.Itoa(pager.PageSize)
	}
	return ordStat + limitStat
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

func buildInsInfoOfRow[T any](meta *TableMeta[T], viewColumns []string) (*T, []any) {
	ins := new(T)
	if len(viewColumns) == 0 {
		return ins, meta.ExtractFieldValues(ins, true, nil)
	}
	return ins, meta.ExtractFieldValuesByColumns(ins, true, viewColumns)
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

func traceLogSQLBefore(ctx context.Context, sql string, args []any) string {
	var argJson []byte
	md5data := []byte(sql)
	argJson, err := json.Marshal(args)
	if err != nil {
		LogError(ctx, err)
	} else {
		md5data = append(md5data, argJson...)
	}
	sqlMd5 := fmt.Sprintf("%X", md5.Sum(md5data))
	LogExecSQLBefore(ctx, sql, argJson, sqlMd5)
	return sqlMd5
}

func traceLogSQLAfter(ctx context.Context, sqlMd5 string, startTime int64) {
	LogExecSQLAfter(ctx, sqlMd5, time.Now().UnixMilli()-startTime)
}
