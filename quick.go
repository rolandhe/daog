// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

type QuickDao[T any] interface {
	GetById(id int64, tc *TransContext) (*T, error)
	GetByIds(ids []int64, tc *TransContext) ([]*T, error)
	QueryListMatcher(m Matcher, tc *TransContext) ([]*T, error)
	QueryListMatcherPageHandle(m Matcher, pageSize int, handler PageHandler[T], tc *TransContext) error
	QueryOneMatcher(m Matcher, tc *TransContext) (*T, error)
	QueryRawSQL(tc *TransContext, extract ExtractScanFieldPoints[T], sql string, args ...any) ([]*T, error)

	Insert(ins *T, tc *TransContext) (int64, error)

	Update(ins *T, tc *TransContext) (int64, error)
	UpdateList(insList []*T, tc *TransContext) (int64, error)
	UpdateById(modifier Modifier, id int64, tc *TransContext) (int64, error)
	UpdateByIds(modifier Modifier, ids []int64, tc *TransContext) (int64, error)
	UpdateByModifier(modifier Modifier, matcher Matcher, tc *TransContext) (int64, error)
	ExecRawSQL(tc *TransContext, sql string, args ...any) (int64, error)

	DeleteById(id int64, tc *TransContext) (int64, error)
	DeleteByIds(ids []int64, tc *TransContext) (int64, error)
	DeleteByMatcher(matcher Matcher, tc *TransContext) (int64, error)
}

func NewBaseQuickDao[T any](meta *TableMeta[T]) QuickDao[T] {
	return &baseQuickDao[T]{meta}
}

type baseQuickDao[T any] struct {
	meta *TableMeta[T]
}

func (dao *baseQuickDao[T]) GetById(id int64, tc *TransContext) (*T, error) {
	return GetById(id, dao.meta, tc)
}

func (dao *baseQuickDao[T]) GetByIds(ids []int64, tc *TransContext) ([]*T, error) {
	return GetByIds(ids, dao.meta, tc)
}

func (dao *baseQuickDao[T]) QueryListMatcher(m Matcher, tc *TransContext) ([]*T, error) {
	return QueryListMatcher(m, dao.meta, tc)
}

func (dao *baseQuickDao[T]) QueryListMatcherPageHandle(m Matcher, pageSize int, handler PageHandler[T], tc *TransContext) error {
	return QueryListMatcherPageHandle(m, dao.meta, pageSize, handler, tc)
}

func (dao *baseQuickDao[T]) QueryOneMatcher(m Matcher, tc *TransContext) (*T, error) {
	return QueryOneMatcher(m, dao.meta, tc)
}

func (dao *baseQuickDao[T]) QueryRawSQL(tc *TransContext, extract ExtractScanFieldPoints[T], sql string, args ...any) ([]*T, error) {
	return QueryRawSQL(tc, extract, sql, args...)
}

func (dao *baseQuickDao[T]) Insert(ins *T, tc *TransContext) (int64, error) {
	return Insert(ins, dao.meta, tc)
}

func (dao *baseQuickDao[T]) Update(ins *T, tc *TransContext) (int64, error) {
	return Update(ins, dao.meta, tc)
}

func (dao *baseQuickDao[T]) UpdateList(insList []*T, tc *TransContext) (int64, error) {
	return UpdateList(insList, dao.meta, tc)
}

func (dao *baseQuickDao[T]) UpdateById(modifier Modifier, id int64, tc *TransContext) (int64, error) {
	return UpdateById(modifier, id, dao.meta, tc)
}

func (dao *baseQuickDao[T]) UpdateByIds(modifier Modifier, ids []int64, tc *TransContext) (int64, error) {
	return UpdateByIds(modifier, ids, dao.meta, tc)
}

func (dao *baseQuickDao[T]) UpdateByModifier(modifier Modifier, matcher Matcher, tc *TransContext) (int64, error) {
	return UpdateByModifier(modifier, matcher, dao.meta, tc)
}

func (dao *baseQuickDao[T]) ExecRawSQL(tc *TransContext, sql string, args ...any) (int64, error) {
	return ExecRawSQL(tc, sql, args...)
}
func (dao *baseQuickDao[T]) DeleteById(id int64, tc *TransContext) (int64, error) {
	return DeleteById(id, dao.meta, tc)
}

func (dao *baseQuickDao[T]) DeleteByIds(ids []int64, tc *TransContext) (int64, error) {
	return DeleteByIds(ids, dao.meta, tc)
}

func (dao *baseQuickDao[T]) DeleteByMatcher(matcher Matcher, tc *TransContext) (int64, error) {
	return DeleteByMatcher(matcher, dao.meta, tc)
}
