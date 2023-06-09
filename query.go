// A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

import (
	"errors"
	"time"
)

var invalidBatchSizeError = errors.New("page size must be greater than 0")

// GetAll 查询表的所有数据
// 可变参数 viewColumns：
//
//	可以指定需要查询的表字段，可以指定多个或者不指定，如果不指定表示要查询所有的表字段
//
//	需要注意的是，表字段指的是数据库表的列名，不是描述表的struct里的属性名
//
//	compile生成的文件中会有表字段的常量，比如 GroupInfo.go 文件中的 GroupInfoFields.Id, 直接使用它，避免手动写字符串
func GetAll[T any](tc *TransContext, meta *TableMeta[T], viewColumns ...string) ([]*T, error) {
	return QueryPageListMatcherWithViewColumns(tc, nil, meta, viewColumns, nil)
}

// GetById 根据指定的主键返回单条数据
// 可变参数 viewColumns：
//
//	可以指定需要查询的表字段，可以指定多个或者不指定，如果不指定表示要查询所有的表字段
//
//	需要注意的是，表字段指的是数据库表的列名，不是描述表的struct里的属性名
//
//	compile生成的文件中会有表字段的常量，比如 GroupInfo.go 文件中的 GroupInfoFields.Id, 直接使用它，避免手动写字符串
func GetById[T any](tc *TransContext, id int64, meta *TableMeta[T], viewColumns ...string) (*T, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return QueryOneMatcher(tc, m, meta, viewColumns...)
}

// GetByIds 根据主键数组返回多条数据
// 可变参数 viewColumns：
//
//	可以指定需要查询的表字段，可以指定多个或者不指定，如果不指定表示要查询所有的表字段
//
//	需要注意的是，表字段指的是数据库表的列名，不是描述表的struct里的属性名
//
//	compile生成的文件中会有表字段的常量，比如 GroupInfo.go 文件中的 GroupInfoFields.Id, 直接使用它，避免手动写字符串
func GetByIds[T any](tc *TransContext, ids []int64, meta *TableMeta[T], viewColumns ...string) ([]*T, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.In(fieldId, ConvertToAnySlice(ids))
	return QueryPageListMatcherWithViewColumns(tc, m, meta, viewColumns, nil)
}

// QueryListMatcher 根据查询条件 Matcher 返回多条数据， 通过与 Matcher 有关的相关函数来构建查询条件
// orders 可变参数：
//
//	可以传入一个、多个或者零个排序条件
//
//	每个条件可以指定排序表字段名及是否是升序要求
func QueryListMatcher[T any](tc *TransContext, m Matcher, meta *TableMeta[T], orders ...*Order) ([]*T, error) {
	return QueryPageListMatcher(tc, m, meta, nil, orders...)
}

// QueryListMatcherWithViewColumns 根据查询条件 Matcher 返回多条数据，每条数据可以是一个视图， 通过与 Matcher 有关的相关函数来构建查询条件, viewColumns 指定需要查询的表字段名，表示一个视图
// orders 可变参数：
//
//	可以传入一个、多个或者零个排序条件
//
//	每个条件可以指定排序表字段名及是否是升序要求
func QueryListMatcherWithViewColumns[T any](tc *TransContext, m Matcher, meta *TableMeta[T], viewColumns []string, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcherWithViewColumns(tc, m, meta, viewColumns, nil, orders...)
}

// QueryPageListMatcher 根据查询条件 Matcher 及 Pager 返回一页数据， 通过与 Matcher 有关的相关函数来构建查询条件， 根据 Pager 相关函数来构建分页条件
// orders 可变参数：
//
//	可以传入一个、多个或者零个排序条件
//
//	每个条件可以指定排序表字段名及是否是升序要求
//
// pager 参数，可以为nil，如果为nil，不分页
func QueryPageListMatcher[T any](tc *TransContext, m Matcher, meta *TableMeta[T], pager *Pager, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcherWithViewColumns(tc, m, meta, nil, pager, orders...)
}

// QueryPageListMatcherWithViewColumns 根据查询条件 Matcher 及 Pager 返回一页数据， 通过与 Matcher 有关的相关函数来构建查询条件， 根据 Pager 相关函数来构建分页条件， viewColumns 指定需要查询的表字段名，表示一个视图
// orders 可变参数：
//
//	可以传入一个、多个或者零个排序条件
//
//	每个条件可以指定排序表字段名及是否是升序要求
//
// pager 参数，可以为nil，如果为nil，不分页
//
// viewColumns 指定需要查询的表字段名，表示一个视图, 可以传入 nil， 表示读取所有字段
func QueryPageListMatcherWithViewColumns[T any](tc *TransContext, m Matcher, meta *TableMeta[T], viewColumns []string, pager *Pager, orders ...*Order) ([]*T, error) {
	sql, args, err := selectQuery(meta, tc.ctx, m, pager, orders, viewColumns)
	if err != nil {
		return nil, err
	}
	return queryRawSQLCore(tc, func() (*T, []any) {
		return buildInsInfoOfRow(meta, viewColumns)
	}, sql, args...)
}

// BatchHandler 处理一批从表中读取的数据的回调函数
// 被 QueryListMatcherByBatchHandle 或者 QueryListMatcherWithViewColumnsByBatchHandle 回调使用， 一般用于从数据库读取大量数据的场景，
// 如果大量数据读入内存会打爆内存，一批批的处理少量数据可以有效的降低内存
type BatchHandler[T any] func(batch []*T) error

// QueryListMatcherByBatchHandle 读取数据并且分批处理数据，当读取的数据量巨大时非常有用，如果数据都读入内存，容易打爆内存，分批量处理就非常有用
// batchSize 每批处理数据的最大容量，必须大于0，但不要设置太大，当设置为1时，退化成每条处理
// handler 用于处理每批数据的函数
// 查询数据最大上限数， 0 表示无上限
func QueryListMatcherByBatchHandle[T any](tc *TransContext, m Matcher, meta *TableMeta[T], totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error {
	return QueryListMatcherWithViewColumnsByBatchHandle(tc, m, meta, nil, totalLimit, batchSize, handler, orders...)
}

// QueryListMatcherWithViewColumnsByBatchHandle 与 QueryListMatcherByBatchHandle 类似，适合分批读取少量数据并回调 BatchHandler 进行处理，与 QueryListMatcherByBatchHandle 稍有不同的是它提供 viewColumns 参数，可以只查询 viewColumns 指定的表字段
// viewColumns 指定需要查询的表字段名，表示一个视图, 可以传入 nil， 表示读取所有字段
func QueryListMatcherWithViewColumnsByBatchHandle[T any](tc *TransContext, m Matcher, meta *TableMeta[T], viewColumns []string, totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error {
	if batchSize <= 0 {
		return invalidBatchSizeError
	}
	var pager *Pager
	if totalLimit > 0 {
		pager = &Pager{0, totalLimit}
	}
	sql, args, err := selectQuery(meta, tc.ctx, m, pager, orders, viewColumns)
	if err != nil {
		return err
	}

	return queryRawSQLByBatchHandleCore(tc, batchSize, handler, func() (*T, []any) {
		return buildInsInfoOfRow(meta, viewColumns)
	}, sql, args...)
}

// QueryOneMatcher 通过 Matcher 条件来查询，但只返回单条数据
// viewColumns 可以指定需要查询的表字段，可以指定多个或者不指定，如果不指定表示要查询所有的表字段
func QueryOneMatcher[T any](tc *TransContext, m Matcher, meta *TableMeta[T], viewColumns ...string) (*T, error) {
	var err error
	err = tc.check()
	if err != nil {
		return nil, err
	}
	sql, args, err := selectQuery(meta, tc.ctx, m, nil, nil, viewColumns)
	if err != nil {
		return nil, err
	}
	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	rows, err := tc.conn.QueryContext(tc.ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	ins, scanFields := buildInsInfoOfRow(meta, viewColumns)
	if err = rows.Scan(scanFields...); err != nil {
		return nil, err
	}
	return ins, nil
}

// ExtractScanFieldPoints 从指定的 *T 类型的对象中抽取出所需要的 field的指针，它是一个回调函数，用于 QueryRawSQL 或者 QueryRawSQLByBatchHandle 函数，
// 其目的是把从数据库读取的一行数据填充到指定 *T 对象中，这对于执行一个原生的 sql 非常有用。
type ExtractScanFieldPoints[T any] func(ins *T) []any

// QueryRawSQL 执行原生select sql语句,返回行数据数组，行数据使用T struct描述 mapper, 它T的各个field指针提取出来并按照顺序生成一个slice，用于Row.Scan方法，把sql字段映射到T对象的各个Field上
// sql 必须是含有 ? 占位符的sql， args 是对应每个 ? 的实参
func QueryRawSQL[T any](tc *TransContext, extract ExtractScanFieldPoints[T], sql string, args ...any) ([]*T, error) {
	return queryRawSQLCore(tc, func() (*T, []any) {
		ins := new(T)
		return ins, extract(ins)
	}, sql, args...)
}

// QueryRawSQLByBatchHandle 与 QueryRawSQL 和 QueryListMatcherByBatchHandle 结合体，执行原生的sql语句，但通过回调 BatchHandler 进行分批业务处理
// batchSize 每批处理数据的最大容量，必须大于0，但不要设置太大，当设置为1时，退化成每条处理
// handler 用于处理每批数据的函数
func QueryRawSQLByBatchHandle[T any](tc *TransContext, batchSize int, handler BatchHandler[T], extract ExtractScanFieldPoints[T], sql string, args ...any) error {
	return queryRawSQLByBatchHandleCore(tc, batchSize, handler, func() (*T, []any) {
		ins := new(T)
		return ins, extract(ins)
	}, sql, args...)
}

func QueryQueryByIdForUpdate[T any](tc *TransContext, id int64, meta *TableMeta[T], skipLocked bool, orders ...*Order) ([]*T, error) {
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.Eq(fieldId, id)
	return QueryQueryListMatcherForUpdate(tc, m, meta, skipLocked, orders...)
}

func QueryQueryListByIdsForUpdate[T any](tc *TransContext, ids []int64, meta *TableMeta[T], skipLocked bool, orders ...*Order) ([]*T, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	m := NewMatcher()
	fieldId := TableIdColumnName
	if meta.AutoColumn != "" {
		fieldId = meta.AutoColumn
	}
	m.In(fieldId, ConvertToAnySlice(ids))
	return QueryQueryListMatcherForUpdate(tc, m, meta, skipLocked, orders...)
}

func QueryQueryListMatcherForUpdate[T any](tc *TransContext, m Matcher, meta *TableMeta[T], skipLocked bool, orders ...*Order) ([]*T, error) {
	return QueryQueryPagerListMatcherWithViewColumnsForUpdate(tc, m, meta, nil, meta.Columns, skipLocked, orders...)
}

func QueryOneMatcherForUpdate[T any](tc *TransContext, m Matcher, meta *TableMeta[T], skipLocked bool, viewColumns ...string) (*T, error) {
	rows, err := QueryQueryListMatcherWithViewColumnsForUpdate(tc, m, meta, viewColumns, skipLocked)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, err
	}
	return rows[0], nil
}

func QueryQueryListMatcherWithViewColumnsForUpdate[T any](tc *TransContext, m Matcher, meta *TableMeta[T], viewColumns []string, skipLocked bool, orders ...*Order) ([]*T, error) {
	return QueryQueryPagerListMatcherWithViewColumnsForUpdate(tc, m, meta, nil, viewColumns, skipLocked, orders...)
}

func QueryQueryPagerListMatcherWithViewColumnsForUpdate[T any](tc *TransContext, m Matcher, meta *TableMeta[T], pager *Pager, viewColumns []string, skipLocked bool, orders ...*Order) ([]*T, error) {
	tableName := GetTableName(tc.ctx, meta)
	sql, params, err := selectQueryCore(m, pager, orders, func() string {
		return getSelectStat(tableName, viewColumns, skipLocked)
	})
	if err != nil {
		return nil, err
	}
	return queryRawSQLCore(tc, func() (*T, []any) {
		return buildInsInfoOfRow(meta, viewColumns)
	}, sql, params...)
}

// Count 表达 select count(*)  语义，其条件通过 Matcher 确定
func Count[T any](tc *TransContext, m Matcher, meta *TableMeta[T]) (int64, error) {
	var err error
	err = tc.check()
	if err != nil {
		return 0, err
	}
	sql, args, err := countQuery(meta, tc.ctx, m)
	if err != nil {
		return 0, err
	}

	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	rows, err := tc.conn.QueryContext(tc.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, nil
	}
	var countValue int64
	if err = rows.Scan(&countValue); err != nil {
		return 0, err
	}
	return countValue, nil
}

type rowInsCreate[T any] func() (*T, []any)

func queryRawSQLByBatchHandleCore[T any](tc *TransContext, batchSize int, handler BatchHandler[T], creatorFunc rowInsCreate[T], sql string, args ...any) error {
	var err error
	err = tc.check()
	if err != nil {
		return err
	}

	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	rows, err := tc.conn.QueryContext(tc.ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	batch := make([]*T, batchSize)
	index := 0

	for rows.Next() {
		ins, scanFields := creatorFunc()
		if err = rows.Scan(scanFields...); err != nil {
			return err
		}
		batch[index] = ins
		index++
		if index == batchSize {
			if err = handler(batch); err != nil {
				return err
			}
			index = 0
		}
	}
	if index == 0 {
		return nil
	}

	batch = batch[:index]
	if err = handler(batch); err != nil {
		return err
	}

	return nil
}

func queryRawSQLCore[T any](tc *TransContext, creatorFunc rowInsCreate[T], sql string, args ...any) ([]*T, error) {
	var err error
	err = tc.check()
	if err != nil {
		return nil, err
	}

	if tc.LogSQL {
		sqlMd5 := traceLogSQLBefore(tc.ctx, sql, args)
		defer traceLogSQLAfter(tc.ctx, sqlMd5, time.Now().UnixMilli())
	}
	rows, err := tc.conn.QueryContext(tc.ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var inses []*T
	for rows.Next() {
		ins, scanFields := creatorFunc()
		if err = rows.Scan(scanFields...); err != nil {
			return nil, err
		}
		inses = append(inses, ins)
	}

	return inses, nil
}
