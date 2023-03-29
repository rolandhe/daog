// A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

import (
	"strings"
)

const (
	selfAdd   = 1
	selfMinus = 2
)

type pair struct {
	column string
	value  any
	self   int
}

// NewModifier 创建 Modifier 对象
func NewModifier() Modifier {
	return &internalModifier{}
}

// Modifier 描述 update 语义中set cause的生成，通过 Modifier 来避免自己拼接sql片段，降低出错概率，
// 最终生成 update tab set xx=?,bb=? 的 sql 片段
type Modifier interface {
	// Add 增加一个字段的修改，比如 id = 100
	Add(column string, value any) Modifier
	SelfAdd(column string, value any) Modifier
	SelfMinus(column string, value any) Modifier
	toSQL(tableName string) (string, []any)
}

type internalModifier struct {
	modifies []*pair
}

func (m *internalModifier) Add(column string, value any) Modifier {
	m.modifies = append(m.modifies, &pair{column, value, 0})
	return m
}

func (m *internalModifier) SelfAdd(column string, value any) Modifier {
	m.modifies = append(m.modifies, &pair{column, value, 1})
	return m
}

func (m *internalModifier) SelfMinus(column string, value any) Modifier {
	m.modifies = append(m.modifies, &pair{column, value, 2})
	return m
}

func (m *internalModifier) toSQL(tableName string) (string, []any) {
	l := len(m.modifies)
	if l == 0 {
		return "", nil
	}
	modStmt := make([]string, l)
	args := make([]any, l)
	for i, p := range m.modifies {
		if p.self == selfAdd {
			modStmt[i] = p.column + "=" + p.column + "+?"
		} else if p.self == selfMinus {
			modStmt[i] = p.column + "=" + p.column + "-?"
		} else {
			modStmt[i] = p.column + "=?"
		}
		args[i] = p.value
	}
	return "update " + tableName + " set " + strings.Join(modStmt, ","), args
}
