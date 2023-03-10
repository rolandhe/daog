// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

type TableMeta[T any] struct {
	LookupFieldFunc func(columnName string, ins *T, point bool) any
	ShardingFunc    func(tableName string, shardingKey any) string
	Table           string
	Columns         []string
	AutoColumn      string
}

func (meta *TableMeta[T]) ExtractFieldValues(ins *T, point bool, exclude map[string]int) []any {
	var ret []any

	for _, column := range meta.Columns {
		if exclude != nil && exclude[column] != 0 {
			continue
		}
		ret = append(ret, meta.LookupFieldFunc(column, ins, point))
	}
	return ret
}

func (meta *TableMeta[T]) ExtractFieldValuesByColumns(ins *T, point bool,columns []string) []any {
	var ret []any

	for _, column := range columns {
		ret = append(ret, meta.LookupFieldFunc(column, ins, point))
	}
	return ret
}
