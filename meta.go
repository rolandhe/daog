// A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

// TableMeta daog中需要表的元数据，基于元数据来自动生成sql，把从数据库读取的数据分配给表的实体对象，TableMeta对应的实例会由compile工具生成。
// TableMeta 需要知道表名，表的列名，自增长字段名称，以及需要提供一个函数LookupFieldFunc，该函数负责根据表的字段名称找到该名称对应的属性。
type TableMeta[T any] struct {
	// 通过表的字段名称获取表实体对象中对应的field的值，或者该field的指针，取值一般用于insert or update，
	// 取指针一般用于从表中读取数据回填到表对象的field中，该函数会被compile自动生成
	LookupFieldFunc func(columnName string, ins *T, point bool) any
	// 在分表情况下，根据分表key生成分表名称的回调函数，该函数不能被compile自动生成，需要使用者在compile生成的xx-ext.go中设置
	ShardingFunc func(tableName string, shardingKey any) string
	Table        string
	Columns      []string
	// 自增长字段的名称，在insert时，表实体对象中对应的field会被自动填充
	AutoColumn string
}

// ExtractFieldValues 从给定的T对象中抽取属性值，并返回，抽取的属性值可能是属性指针，也可能是属性的值，
// 通过exclude可以指定哪些列对应的属性被排除，exclude 中key是数据库表的字段名，不是表实体对象中的属性名
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

// ExtractFieldValuesByColumns 与ExtractFieldValues，但它仅仅抽取通过 columns 参数指定的数据表列所对应的属性
func (meta *TableMeta[T]) ExtractFieldValuesByColumns(ins *T, point bool, columns []string) []any {
	var ret []any

	for _, column := range columns {
		ret = append(ret, meta.LookupFieldFunc(column, ins, point))
	}
	return ret
}
