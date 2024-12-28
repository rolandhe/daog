// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

package daog

// build an equals condition, e.g, name = ?
func newEqCond(column string, value any) SQLCond {
	return newSimpleCond("=", column, value)
}

// build not equals condition, e.g, name != ?
func newNeCond(column string, value any) SQLCond {
	return newSimpleCond("!=", column, value)
}

// build greater than condition, e.g, total > ?
func newGtCond(column string, value any) SQLCond {
	return newSimpleCond(">", column, value)
}

// build greater than or equals condition, e.g, total >= ?
func newGteCond(column string, value any) SQLCond {
	return newSimpleCond(">=", column, value)
}

func newLtCond(column string, value any) SQLCond {
	return newSimpleCond("<", column, value)
}

func newLteCond(column string, value any) SQLCond {
	return newSimpleCond("<=", column, value)
}

func newInCond(column string, values []any) SQLCond {
	return &inCond{
		column: column,
		values: values,
	}
}

func newNotInCond(column string, values []any) SQLCond {
	return &inCond{
		column: column,
		values: values,
		not:    true,
	}
}

func newLikeCond(column string, value string, likeStyle int) SQLCond {
	return &likeCond{
		column, value, likeStyle,
	}
}

func newNullCond(column string, not bool) SQLCond {
	return &nullCond{
		column, not,
	}
}

func newBetweenCond(column string, start any, end any) SQLCond {
	return &betweenCond{
		column, start, end,
	}
}

func newSimpleCond(op string, column string, value any) SQLCond {
	return &simpleCond{
		op, column, value,
	}
}

func newScalarCond(cond string) SQLCond {
	return &scalarCond{
		cond: cond,
	}
}
func newBitwiseAndCond(column string, mask any, target any) SQLCond {
	return &bitwiseAndCond{
		column: column,
		mask:   mask,
		target: target,
	}
}
