package daog

func newEqCond(column string, value any) SQLCond {
	return createSimpleCond("=", column, value)
}

func newNeCond(column string, value any) SQLCond {
	return createSimpleCond("!=", column, value)
}

func newGtCond(column string, value any) SQLCond {
	return createSimpleCond(">", column, value)
}

func newGteCond(column string, value any) SQLCond {
	return createSimpleCond(">=", column, value)
}

func newLtCond(column string, value any) SQLCond {
	return createSimpleCond("<", column, value)
}

func newLteCond(column string, value any) SQLCond {
	return createSimpleCond("<=", column, value)
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

func createSimpleCond(op string, column string, value any) SQLCond {
	return &simpleCond{
		op, column, value,
	}
}
