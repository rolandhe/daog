// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

package daog

import (
	"errors"
	"strings"
)

// define "and","or" operand
// define like style, %xx%, %s,s%
const (
	logicOpAnd = "and"
	logicOpOr  = "or"

	// LikeStyleAll ,like "%value%"
	LikeStyleAll = 0
	// LikeStyleLeft ,like "%value"
	LikeStyleLeft = 1
	// LikeStyleRight ,like "value%"
	LikeStyleRight = 2
)

// NewMatcher 构建一个以 and 连接的匹配条件构建器
func NewMatcher() Matcher {
	return NewAndMatcher()
}

// NewAndMatcher 构建一个以 and 连接的匹配条件构建器
func NewAndMatcher() Matcher {
	return &compositeCond{logicOp: logicOpAnd}
}

// NewOrMatcher 构建一个以 or 连接的匹配条件构建器
func NewOrMatcher() Matcher {
	return &compositeCond{logicOp: logicOpOr}
}

// SQLCond 抽象描述一个sql的条件，可以是 单个字段的条件，比如 name=?, 也可以是通过连接操作符(and/or)连接的多个条件。
// 每一个条件以 [字段 操作符 值占位符] 的方式组成，比如 id = ?,生成条件时需要传入每个占位符对应一个参数值
// 也可以直接给一个标量条件，没有参数，比如 status = 0
type SQLCond interface {
	// ToSQL 生成包含?占位符的sql，并返回对应的参数数组
	// 输入参数 args是已经收集到的参数
	ToSQL(args []any) (string, []any, error)
}

// Matcher sql where条件的构建器，用以构造以 and 或者 or 连接的各种条件, 最后拼接生成一个可用的、包含?占位符的where条件,并且收集所有对应?的参数数组
type Matcher interface {
	SQLCond

	// Add 这个方法是个冗余方法，其实可以直接使用AddCond, 因为Matcher继承自SQLCond，冗余这个方法仅仅是让使用者更好的理解
	Add(matcher Matcher) Matcher

	// AddCond 加入一个新的到条件
	AddCond(cond SQLCond) Matcher

	// Eq 快速生成一个等于语义的条件，比如 id = 100
	// column 是数据库表的字段名，value是条件值
	Eq(column string, value any) Matcher

	// Ne 快速生成一个 not equals 条件语义， Eq 的反向
	Ne(column string, value any) Matcher

	// Lt 快速生成一个 less than 条件语义
	Lt(column string, value any) Matcher

	// Lte 快速生成一个 less than or equals 条件语义
	Lte(column string, value any) Matcher

	// Gt 快速生成一个greater than 条件语义
	Gt(column string, value any) Matcher

	// Gte 快速生成一个 greater than or equals 条件语义
	Gte(column string, value any) Matcher

	// In 快速生成in 条件语义，比如 xx in(?,?,...)
	In(column string, values []any) Matcher

	// NotIn 快速生成 not in 语义，比如 xx not in(?,?,...)
	NotIn(column string, values []any) Matcher

	// Like 快速生成 like 条件语义， 参数 likeStyle对应 枚举值： LikeStyleAll/ LikeStyleLeft / LikeStyleRight
	Like(column string, value string, likeStyle int) Matcher

	// Null 快速生成是否为空的条件语义，比如 name is null,  参数not表示是否为not null， 如果为true， 则生成条件 name is not null
	Null(column string, not bool) Matcher

	// Between 快速生成between 语义, start 和end可以有一个为 nil， 如果 start = nil，则退化成 column <= end, 如果 end = nil，则退化成 column >= end
	// 两个都不为nil，生成标准的between语义
	Between(column string, start any, end any) Matcher

	// AddScalar 增加一个标量条件，即增加一个条件字符串，比如 "id = 100", 或者 "name = 'Joe'",
	// 注意 尽量不要使用这个方法，因为它容易引起sql注入，如果你需要使用这个方法，你一定要使用转义来防止sql注入
	AddScalar(cond string) Matcher
}

type compositeCond struct {
	conds   []SQLCond
	logicOp string
}

func (cc *compositeCond) Add(matcher Matcher) Matcher {
	cc.conds = append(cc.conds, matcher)
	return cc
}

func (cc *compositeCond) AddCond(cond SQLCond) Matcher {
	cc.conds = append(cc.conds, cond)
	return cc
}

func (cc *compositeCond) Eq(column string, value any) Matcher {
	cc.conds = append(cc.conds, newEqCond(column, value))
	return cc
}
func (cc *compositeCond) Ne(column string, value any) Matcher {
	cc.conds = append(cc.conds, newNeCond(column, value))
	return cc
}

func (cc *compositeCond) Lt(column string, value any) Matcher {
	cc.conds = append(cc.conds, newLtCond(column, value))
	return cc
}
func (cc *compositeCond) Lte(column string, value any) Matcher {
	cc.conds = append(cc.conds, newLteCond(column, value))
	return cc
}

func (cc *compositeCond) Gt(column string, value any) Matcher {
	cc.conds = append(cc.conds, newGtCond(column, value))
	return cc
}
func (cc *compositeCond) Gte(column string, value any) Matcher {
	cc.conds = append(cc.conds, newGteCond(column, value))
	return cc
}
func (cc *compositeCond) In(column string, values []any) Matcher {
	cc.conds = append(cc.conds, newInCond(column, values))
	return cc
}
func (cc *compositeCond) NotIn(column string, values []any) Matcher {
	cc.conds = append(cc.conds, newNotInCond(column, values))
	return cc
}

func (cc *compositeCond) Like(column string, value string, likeStyle int) Matcher {
	cc.conds = append(cc.conds, newLikeCond(column, value, likeStyle))
	return cc
}

func (cc *compositeCond) Null(column string, not bool) Matcher {
	cc.conds = append(cc.conds, newNullCond(column, not))
	return cc
}

func (cc *compositeCond) Between(column string, start any, end any) Matcher {
	cc.conds = append(cc.conds, newBetweenCond(column, start, end))
	return cc
}

func (cc *compositeCond) AddScalar(cond string) Matcher {
	cc.conds = append(cc.conds, newScalarCond(cond))
	return cc
}

func (cc *compositeCond) ToSQL(args []any) (string, []any, error) {
	var condSegs []string

	if len(cc.conds) == 0 {
		return "", args, nil
	}

	for _, cond := range cc.conds {
		s, a, err := cond.ToSQL(args)
		if err != nil {
			return "", nil, err
		}
		if s == "" {
			continue
		}
		condSegs = append(condSegs, s)
		args = a
	}

	l := len(condSegs)
	if l == 0 {
		return "", args, nil
	}
	sql := strings.Join(condSegs, " "+cc.logicOp+" ")
	if cc.logicOp == logicOpOr && l > 1 {
		sql = "(" + sql + ")"
	}
	return sql, args, nil
}

type simpleCond struct {
	op     string
	column string
	value  any
}

func (sc *simpleCond) ToSQL(args []any) (string, []any, error) {
	return sc.column + " " + sc.op + " ?", append(args, sc.value), nil
}

type inCond struct {
	column string
	values []any
	not    bool
}

func (ic *inCond) ToSQL(args []any) (string, []any, error) {
	if len(ic.values) == 0 {
		return "", args, errors.New(ic.column + ": no param values")
	}
	holders := make([]string, len(ic.values))
	for i := 0; i < len(ic.values); i++ {
		holders[i] = "?"
	}

	var builder strings.Builder
	builder.WriteString(ic.column)
	if ic.not {
		builder.WriteString(" not in (")
	} else {
		builder.WriteString(" in (")
	}
	builder.WriteString(strings.Join(holders, ","))
	builder.WriteString(")")
	return builder.String(), append(args, ic.values...), nil
}

type betweenCond struct {
	column string
	start  any
	end    any
}

func (btc *betweenCond) ToSQL(args []any) (string, []any, error) {
	if btc.start == nil && btc.end == nil {
		return "", args, nil
	}

	if btc.start != nil && btc.end == nil {
		return btc.column + " >= ?", append(args, btc.start), nil
	}

	if btc.start == nil && btc.end != nil {
		return btc.column + " <= ?", append(args, btc.end), nil
	}
	return btc.column + " between ? and ?", append(args, btc.start, btc.end), nil
}

type nullCond struct {
	column string
	not    bool
}

func (nc *nullCond) ToSQL(args []any) (string, []any, error) {
	if nc.not {
		return nc.column + " is not null", args, nil
	}

	return nc.column + " is null", args, nil
}

type likeCond struct {
	column    string
	value     string
	likeStyle int
}

func (likec *likeCond) ToSQL(args []any) (string, []any, error) {
	if likec.value == "" {
		return "", args, nil
	}
	v := likec.value
	switch likec.likeStyle {
	case LikeStyleLeft:
		v = "%" + v
	case LikeStyleRight:
		v = v + "%"
	case LikeStyleAll:
		v = "%" + v + "%"
	default:
		return "", args, nil
	}

	return likec.column + " like ?", append(args, v), nil
}

type scalarCond struct {
	cond string
}

func (scalar *scalarCond) ToSQL(args []any) (string, []any, error) {
	return scalar.cond, args, nil
}
