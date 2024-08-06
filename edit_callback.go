// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

package daog

// define callback func

var BeforeInsertCallback func(tableName string, ins any) error
var BeforeUpdateCallback func(tableName string, ins any) error
var BeforeModifyCallback func(tableName string, modi Modifier, columns []string, values []any) error
