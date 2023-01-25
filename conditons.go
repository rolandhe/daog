// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

import (
	"fmt"
	"strings"
)

// define "and","or" operand
// define like style, %xx%, %s,s%
const (
	logicOpAnd     = "and"
	logicOpOr      = "or"
	LikeStyleAll   = 0
	LikeStyleLeft  = 1
	LikeStyleRight = 2
)

func NewMatcher() Matcher {
	return NewAndMatcher()
}

func NewAndMatcher() Matcher {
	return &compositeCond{logicOp: logicOpAnd}
}

func NewOrMatcher() Matcher {
	return &compositeCond{logicOp: logicOpOr}
}

type SQLCond interface {
	ToSQL(args []any) (string, []any)
}

type Matcher interface {
	SQLCond
	Add(matcher Matcher) Matcher
	AddCond(cond SQLCond) Matcher
	Eq(column string, value any) Matcher
	Ne(column string, value any) Matcher
	Lt(column string, value any) Matcher
	Lte(column string, value any) Matcher
	Gt(column string, value any) Matcher
	Gte(column string, value any) Matcher
	In(column string, values []any) Matcher
	NotIn(column string, values []any) Matcher
	Like(column string, value string, likeStyle int) Matcher
	Null(column string, not bool) Matcher
	Between(column string, start any, end any) Matcher
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

func (cc *compositeCond) ToSQL(args []any) (string, []any) {
	var condSegs []string

	if len(cc.conds) == 0 {
		return "", args
	}

	for _, cond := range cc.conds {
		s, a := cond.ToSQL(args)
		if s == "" {
			continue
		}
		condSegs = append(condSegs, s)
		args = a
	}

	l := len(condSegs)
	if l == 0 {
		return "", args
	}
	sql := strings.Join(condSegs, fmt.Sprintf(" %s ", cc.logicOp))
	if cc.logicOp == logicOpOr && l > 1 {
		sql = "(" + sql + ")"
	}
	return sql, args
}

type simpleCond struct {
	op     string
	column string
	value  any
}

func (sc *simpleCond) ToSQL(args []any) (string, []any) {
	return fmt.Sprintf("%s %s ?", sc.column, sc.op), append(args, sc.value)
}

type inCond struct {
	column string
	values []any
	not    bool
}

func (ic *inCond) ToSQL(args []any) (string, []any) {
	if len(ic.values) == 0 {
		return "", nil
	}
	holders := make([]string, len(ic.values))
	for i := 0; i < len(ic.values); i++ {
		holders[i] = "?"
	}

	sfmt := "%s in (%s)"
	if ic.not {
		sfmt = "%s not in (%s)"
	}
	return fmt.Sprintf(sfmt, ic.column, strings.Join(holders, ",")), append(args, ic.values...)
}

type betweenCond struct {
	column string
	start  any
	end    any
}

func (btc *betweenCond) ToSQL(args []any) (string, []any) {
	if btc.start == nil && btc.end == nil {
		return "", nil
	}

	if btc.start != nil && btc.end == nil {
		return fmt.Sprintf("%s >= ?", btc.column), append(args, btc.start)
	}

	if btc.start == nil && btc.end != nil {
		return fmt.Sprintf("%s <= ?", btc.column), append(args, btc.end)
	}

	return fmt.Sprintf("%s between ? and ?", btc.column), append(args, btc.start, btc.end)
}

type nullCond struct {
	column string
	not    bool
}

func (nc *nullCond) ToSQL(args []any) (string, []any) {
	if nc.not {
		return fmt.Sprintf("%s is not null", nc.column), args
	}
	return fmt.Sprintf("%s is null", nc.column), args
}

type likeCond struct {
	column    string
	value     string
	likeStyle int
}

func (likec *likeCond) ToSQL(args []any) (string, []any) {
	if likec.value == "" {
		return "", args
	}
	v := likec.value
	switch likec.likeStyle {
	case LikeStyleRight:
		v = "%" + v
	case LikeStyleLeft:
		v = v + "%"
	case LikeStyleAll:
		v = "%" + v + "%"
	default:
		return "", args
	}

	return fmt.Sprintf("%s like ?", likec.column), append(args, v)
}
