// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

import (
	"fmt"
	"strings"
)

type pair struct {
	column string
	value  any
}

type Modifier struct {
	modifies []*pair
}

func (m *Modifier) Add(column string, value any) *Modifier {
	m.modifies = append(m.modifies, &pair{column, value})
	return m
}

func (m *Modifier) toSQL(tableName string) (string, []any) {
	l := len(m.modifies)
	if l == 0 {
		return "", nil
	}
	modStmt := make([]string, l)
	args := make([]any, l)
	for i, p := range m.modifies {
		modStmt[i] = p.column + "=?"
		args[i] = p.value
	}
	return fmt.Sprintf("update %s set %s", tableName, strings.Join(modStmt, ",")), args
}
