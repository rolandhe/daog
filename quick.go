// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

type QuickDao[T any] interface {
	GetAll(tc *TransContext, viewColumns ...string) ([]*T, error)
	GetById(tc *TransContext, id int64, viewColumns ...string) (*T, error)
	GetByIds(tc *TransContext, ids []int64, viewColumns ...string) ([]*T, error)
	QueryListMatcher(tc *TransContext, m Matcher, orders ...*Order) ([]*T, error)
	QueryPageListMatcher(tc *TransContext, m Matcher, pager *Pager, orders ...*Order) ([]*T, error)
	QueryPageListMatcherWithViewColumns(tc *TransContext, m Matcher, viewColumns []string, pager *Pager, orders ...*Order) ([]*T, error)
	QueryListMatcherByBatchHandle(tc *TransContext, m Matcher, totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error
	QueryListMatcherWithViewColumnsByBatchHandle(tc *TransContext, m Matcher, viewColumns []string, totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error
	QueryOneMatcher(tc *TransContext, m Matcher, viewColumns ...string) (*T, error)
	QueryRawSQL(tc *TransContext, extract ExtractScanFieldPoints[T], sql string, args ...any) ([]*T, error)
	QueryRawSQLByBatchHandle(tc *TransContext, batchSize int, handler BatchHandler[T], extract ExtractScanFieldPoints[T], sql string, args ...any) error
	Count(tc *TransContext, m Matcher) (int64, error)

	Insert(tc *TransContext, ins *T) (int64, error)

	Update(tc *TransContext, ins *T) (int64, error)
	UpdateList(tc *TransContext, insList []*T) (int64, error)
	UpdateById(tc *TransContext, modifier Modifier, id int64) (int64, error)
	UpdateByIds(tc *TransContext, modifier Modifier, ids []int64) (int64, error)
	UpdateByModifier(tc *TransContext, modifier Modifier, matcher Matcher) (int64, error)
	ExecRawSQL(tc *TransContext, sql string, args ...any) (int64, error)

	DeleteById(tc *TransContext, id int64) (int64, error)
	DeleteByIds(tc *TransContext, ids []int64) (int64, error)
	DeleteByMatcher(tc *TransContext, matcher Matcher) (int64, error)
}

func NewBaseQuickDao[T any](meta *TableMeta[T]) QuickDao[T] {
	return &baseQuickDao[T]{meta}
}

type baseQuickDao[T any] struct {
	meta *TableMeta[T]
}

func (dao *baseQuickDao[T]) GetAll(tc *TransContext, viewColumns ...string) ([]*T, error) {
	return GetAll(tc, dao.meta, viewColumns...)
}

func (dao *baseQuickDao[T]) GetById(tc *TransContext, id int64, viewColumns ...string) (*T, error) {
	return GetById(tc, id, dao.meta, viewColumns...)
}

func (dao *baseQuickDao[T]) GetByIds(tc *TransContext, ids []int64, viewColumns ...string) ([]*T, error) {
	return GetByIds(tc, ids, dao.meta, viewColumns...)
}

func (dao *baseQuickDao[T]) QueryListMatcher(tc *TransContext, m Matcher, orders ...*Order) ([]*T, error) {
	return QueryListMatcher(tc, m, dao.meta, orders...)
}

func (dao *baseQuickDao[T]) QueryPageListMatcher(tc *TransContext, m Matcher, pager *Pager, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcher(tc, m, dao.meta, pager, orders...)
}

func (dao *baseQuickDao[T]) QueryPageListMatcherWithViewColumns(tc *TransContext, m Matcher, viewColumns []string, pager *Pager, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcherWithViewColumns(tc, m, dao.meta, viewColumns, pager, orders...)
}

func (dao *baseQuickDao[T]) QueryListMatcherByBatchHandle(tc *TransContext, m Matcher, totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error {
	return QueryListMatcherByBatchHandle(tc, m, dao.meta, totalLimit, batchSize, handler, orders...)
}

func (dao *baseQuickDao[T]) QueryListMatcherWithViewColumnsByBatchHandle(tc *TransContext, m Matcher, viewColumns []string, totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error {
	return QueryListMatcherWithViewColumnsByBatchHandle(tc, m, dao.meta, viewColumns, totalLimit, batchSize, handler, orders...)
}

func (dao *baseQuickDao[T]) QueryOneMatcher(tc *TransContext, m Matcher, viewColumns ...string) (*T, error) {
	return QueryOneMatcher(tc, m, dao.meta, viewColumns...)
}

func (dao *baseQuickDao[T]) QueryRawSQL(tc *TransContext, extract ExtractScanFieldPoints[T], sql string, args ...any) ([]*T, error) {
	return QueryRawSQL(tc, extract, sql, args...)
}

func (dao *baseQuickDao[T]) QueryRawSQLByBatchHandle(tc *TransContext, batchSize int, handler BatchHandler[T], extract ExtractScanFieldPoints[T], sql string, args ...any) error {
	return QueryRawSQLByBatchHandle(tc, batchSize, handler, extract, sql, args...)
}

func (dao *baseQuickDao[T]) Count(tc *TransContext, m Matcher) (int64, error) {
	return Count(tc, m, dao.meta)
}

func (dao *baseQuickDao[T]) Insert(tc *TransContext, ins *T) (int64, error) {
	return Insert(tc, ins, dao.meta)
}

func (dao *baseQuickDao[T]) Update(tc *TransContext, ins *T) (int64, error) {
	return Update(tc, ins, dao.meta)
}

func (dao *baseQuickDao[T]) UpdateList(tc *TransContext, insList []*T) (int64, error) {
	return UpdateList(tc, insList, dao.meta)
}

func (dao *baseQuickDao[T]) UpdateById(tc *TransContext, modifier Modifier, id int64) (int64, error) {
	return UpdateById(tc, modifier, id, dao.meta)
}

func (dao *baseQuickDao[T]) UpdateByIds(tc *TransContext, modifier Modifier, ids []int64) (int64, error) {
	return UpdateByIds(tc, modifier, ids, dao.meta)
}

func (dao *baseQuickDao[T]) UpdateByModifier(tc *TransContext, modifier Modifier, matcher Matcher) (int64, error) {
	return UpdateByModifier(tc, modifier, matcher, dao.meta)
}

func (dao *baseQuickDao[T]) ExecRawSQL(tc *TransContext, sql string, args ...any) (int64, error) {
	return ExecRawSQL(tc, sql, args...)
}
func (dao *baseQuickDao[T]) DeleteById(tc *TransContext, id int64) (int64, error) {
	return DeleteById(tc, id, dao.meta)
}

func (dao *baseQuickDao[T]) DeleteByIds(tc *TransContext, ids []int64) (int64, error) {
	return DeleteByIds(tc, ids, dao.meta)
}

func (dao *baseQuickDao[T]) DeleteByMatcher(tc *TransContext, matcher Matcher) (int64, error) {
	return DeleteByMatcher(tc, matcher, dao.meta)
}
