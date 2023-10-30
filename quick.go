// A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

// QuickDao 类似于java的Dao接口，用于仅仅访问一张表。
// 针对每一张表生成一个实现 QuickDao 的struct的实例，使用这个实例来操作这张表。相对于直接调用诸如 GetAll 、 GetById 等函数可以少传入 TableMeta 对象，并且让代码看起来更面向对象。
// 针对每一张表的 QuickDao 不用自行实现， compile 会自动生成，比如 GroupInfo.go文件中GroupInfoDao
type QuickDao[T any] interface {
	// GetAll 封装 GetAll 函数
	GetAll(tc *TransContext, viewColumns ...string) ([]*T, error)

	// GetAllWithViewObj 封装 GetAllWithViewObj 函数
	GetAllWithViewObj(tc *TransContext, view *View) ([]*T, error)
	// GetById 封装 GetById 函数
	GetById(tc *TransContext, id int64, viewColumns ...string) (*T, error)

	// GetByIdWithViewObj 封装 GetByIdWithViewObj 函数
	GetByIdWithViewObj(tc *TransContext, id int64, view *View) (*T, error)

	// GetByIdForUpdate 封装 GetByIdForUpdate 函数
	GetByIdForUpdate(tc *TransContext, id int64, skipLocked bool, viewColumns ...string) ([]*T, error)

	// GetByIds 封装 GetByIds 函数
	GetByIds(tc *TransContext, ids []int64, viewColumns ...string) ([]*T, error)

	// GetByIdsWithViewObj 封装 GetByIdsWithViewObj 函数
	GetByIdsWithViewObj(tc *TransContext, ids []int64, view *View) ([]*T, error)

	// GetByIdsForUpdate 封装 GetByIdsForUpdate 函数
	GetByIdsForUpdate(tc *TransContext, ids []int64, skipLocked bool, viewColumns ...string) ([]*T, error)
	// QueryListMatcher 封装 QueryListMatcher 函数
	QueryListMatcher(tc *TransContext, m Matcher, orders ...*Order) ([]*T, error)
	// QueryListMatcherWithViewColumns 封装 QueryListMatcherWithViewColumns 函数
	QueryListMatcherWithViewColumns(tc *TransContext, m Matcher, viewColumns []string, orders ...*Order) ([]*T, error)
	// QueryListMatcherWithViewObj 封装 QueryListMatcherWithViewObj 函数
	QueryListMatcherWithViewObj(tc *TransContext, m Matcher, view *View, orders ...*Order) ([]*T, error)

	// QueryListMatcherWithViewColumnsForUpdate 封装 QueryListMatcherWithViewColumnsForUpdate 函数
	QueryListMatcherWithViewColumnsForUpdate(tc *TransContext, m Matcher, viewColumns []string, skipLocked bool, orders ...*Order) ([]*T, error)
	// QueryPageListMatcher 封装 QueryPageListMatcher 函数
	QueryPageListMatcher(tc *TransContext, m Matcher, pager *Pager, orders ...*Order) ([]*T, error)
	// QueryPageListMatcherForUpdate 封装 QueryPageListMatcherForUpdate 函数
	QueryPageListMatcherForUpdate(tc *TransContext, m Matcher, pager *Pager, skipLocked bool, orders ...*Order) ([]*T, error)
	// QueryListMatcherForUpdate 封装 QueryListMatcherForUpdate 函数
	QueryListMatcherForUpdate(tc *TransContext, m Matcher, skipLocked bool, orders ...*Order) ([]*T, error)
	// QueryPageListMatcherWithViewColumns 封装 QueryPageListMatcherWithViewColumns 函数
	QueryPageListMatcherWithViewColumns(tc *TransContext, m Matcher, viewColumns []string, pager *Pager, orders ...*Order) ([]*T, error)

	// QueryPageListMatcherWithViewObj 封装 QueryPageListMatcherWithViewObj 函数
	QueryPageListMatcherWithViewObj(tc *TransContext, m Matcher, view *View, pager *Pager, orders ...*Order) ([]*T, error)

	// QueryPageListMatcherWithViewColumnsForUpdate 封装 QueryPageListMatcherWithViewColumnsForUpdate 函数
	QueryPageListMatcherWithViewColumnsForUpdate(tc *TransContext, m Matcher, viewColumns []string, pager *Pager, skipLocked bool, orders ...*Order) ([]*T, error)
	// QueryListMatcherByBatchHandle 封装 QueryListMatcherByBatchHandle 函数
	QueryListMatcherByBatchHandle(tc *TransContext, m Matcher, totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error
	// QueryListMatcherWithViewColumnsByBatchHandle 封装 QueryListMatcherWithViewColumnsByBatchHandle 函数
	QueryListMatcherWithViewColumnsByBatchHandle(tc *TransContext, m Matcher, viewColumns []string, totalLimit int, batchSize int, handler BatchHandler[T], orders ...*Order) error
	// QueryOneMatcher 封装 QueryOneMatcher 函数
	QueryOneMatcher(tc *TransContext, m Matcher, viewColumns ...string) (*T, error)

	// QueryOneMatcherWithViewObj 封装 QueryOneMatcherWithViewObj 函数
	QueryOneMatcherWithViewObj(tc *TransContext, m Matcher, view *View) (*T, error)

	// QueryOneMatcherForUpdate 封装 QueryOneMatcherForUpdate 函数
	QueryOneMatcherForUpdate(tc *TransContext, m Matcher, skipLocked bool, viewColumns ...string) (*T, error)
	// QueryRawSQL 封装 QueryRawSQL 函数
	QueryRawSQL(tc *TransContext, extract ExtractScanFieldPoints[T], sql string, args ...any) ([]*T, error)
	// QueryRawSQLByBatchHandle 封装 QueryRawSQLByBatchHandle 函数
	QueryRawSQLByBatchHandle(tc *TransContext, batchSize int, handler BatchHandler[T], extract ExtractScanFieldPoints[T], sql string, args ...any) error

	// Count 封装 Count 函数
	Count(tc *TransContext, m Matcher) (int64, error)

	// Insert 封装 Insert 函数
	Insert(tc *TransContext, ins *T) (int64, error)

	// Update 封装 Update 函数
	Update(tc *TransContext, ins *T) (int64, error)
	// UpdateList 封装 UpdateList 函数
	UpdateList(tc *TransContext, insList []*T) (int64, error)
	// UpdateById 封装 UpdateById 函数
	UpdateById(tc *TransContext, modifier Modifier, id int64) (int64, error)
	// UpdateByIds 封装 UpdateByIds 函数
	UpdateByIds(tc *TransContext, modifier Modifier, ids []int64) (int64, error)
	// UpdateByModifier 封装 UpdateByModifier 函数
	UpdateByModifier(tc *TransContext, modifier Modifier, matcher Matcher) (int64, error)
	// ExecRawSQL 封装 ExecRawSQL 函数
	ExecRawSQL(tc *TransContext, sql string, args ...any) (int64, error)

	// DeleteById 封装 DeleteById 函数
	DeleteById(tc *TransContext, id int64) (int64, error)
	// DeleteByIds 封装 DeleteByIds 函数
	DeleteByIds(tc *TransContext, ids []int64) (int64, error)
	// DeleteByMatcher 封装 GetById 函数
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

func (dao *baseQuickDao[T]) GetAllWithViewObj(tc *TransContext, view *View) ([]*T, error) {
	return GetAllWithViewObj(tc, dao.meta, view)
}

func (dao *baseQuickDao[T]) GetById(tc *TransContext, id int64, viewColumns ...string) (*T, error) {
	return GetById(tc, id, dao.meta, viewColumns...)
}

func (dao *baseQuickDao[T]) GetByIdWithViewObj(tc *TransContext, id int64, view *View) (*T, error) {
	return GetByIdWithViewObj(tc, id, dao.meta, view)
}

func (dao *baseQuickDao[T]) GetByIdForUpdate(tc *TransContext, id int64, skipLocked bool, viewColumns ...string) ([]*T, error) {
	return GetByIdForUpdate(tc, id, dao.meta, skipLocked, viewColumns...)
}

func (dao *baseQuickDao[T]) GetByIds(tc *TransContext, ids []int64, viewColumns ...string) ([]*T, error) {
	return GetByIds(tc, ids, dao.meta, viewColumns...)
}

func (dao *baseQuickDao[T]) GetByIdsWithViewObj(tc *TransContext, ids []int64, view *View) ([]*T, error) {
	return GetByIdsWithViewObj(tc, ids, dao.meta, view)
}

func (dao *baseQuickDao[T]) GetByIdsForUpdate(tc *TransContext, ids []int64, skipLocked bool, viewColumns ...string) ([]*T, error) {
	return GetByIdsForUpdate(tc, ids, dao.meta, skipLocked, viewColumns...)
}

func (dao *baseQuickDao[T]) QueryListMatcher(tc *TransContext, m Matcher, orders ...*Order) ([]*T, error) {
	return QueryListMatcher(tc, m, dao.meta, orders...)
}

func (dao *baseQuickDao[T]) QueryListMatcherForUpdate(tc *TransContext, m Matcher, skipLocked bool, orders ...*Order) ([]*T, error) {
	return QueryListMatcherForUpdate(tc, m, dao.meta, skipLocked, orders...)
}

func (dao *baseQuickDao[T]) QueryListMatcherWithViewColumns(tc *TransContext, m Matcher, viewColumns []string, orders ...*Order) ([]*T, error) {
	return QueryListMatcherWithViewColumns(tc, m, dao.meta, viewColumns, orders...)
}

func (dao *baseQuickDao[T]) QueryListMatcherWithViewObj(tc *TransContext, m Matcher, view *View, orders ...*Order) ([]*T, error) {
	return QueryListMatcherWithViewObj(tc, m, dao.meta, view, orders...)
}

func (dao *baseQuickDao[T]) QueryListMatcherWithViewColumnsForUpdate(tc *TransContext, m Matcher, viewColumns []string, skipLocked bool, orders ...*Order) ([]*T, error) {
	return QueryListMatcherWithViewColumnsForUpdate(tc, m, dao.meta, viewColumns, skipLocked, orders...)
}

func (dao *baseQuickDao[T]) QueryPageListMatcher(tc *TransContext, m Matcher, pager *Pager, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcher(tc, m, dao.meta, pager, orders...)
}

func (dao *baseQuickDao[T]) QueryPageListMatcherForUpdate(tc *TransContext, m Matcher, pager *Pager, skipLocked bool, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcherForUpdate(tc, m, dao.meta, pager, skipLocked, orders...)
}

func (dao *baseQuickDao[T]) QueryPageListMatcherWithViewColumns(tc *TransContext, m Matcher, viewColumns []string, pager *Pager, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcherWithViewColumns(tc, m, dao.meta, viewColumns, pager, orders...)
}

func (dao *baseQuickDao[T]) QueryPageListMatcherWithViewObj(tc *TransContext, m Matcher, view *View, pager *Pager, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcherWithViewObj(tc, m, dao.meta, view, pager, orders...)
}

func (dao *baseQuickDao[T]) QueryPageListMatcherWithViewColumnsForUpdate(tc *TransContext, m Matcher, viewColumns []string, pager *Pager, skipLocked bool, orders ...*Order) ([]*T, error) {
	return QueryPageListMatcherWithViewColumnsForUpdate(tc, m, dao.meta, viewColumns, pager, skipLocked, orders...)
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

func (dao *baseQuickDao[T]) QueryOneMatcherWithViewObj(tc *TransContext, m Matcher, view *View) (*T, error) {
	return QueryOneMatcherWithViewObj(tc, m, dao.meta, view)
}

func (dao *baseQuickDao[T]) QueryOneMatcherForUpdate(tc *TransContext, m Matcher, skipLocked bool, viewColumns ...string) (*T, error) {
	return QueryOneMatcherForUpdate(tc, m, dao.meta, skipLocked, viewColumns...)
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
