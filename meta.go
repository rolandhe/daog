package daog

import "context"

type TableMeta[T any] struct {
	InstanceFunc    func() *T
	LookupFieldFunc func(columnName string, ins *T,point bool) any
	ShardingFunc    func(tableName string, ctx context.Context) string
	Table           string
	Columns         []string
	AutoColumn      string
}

func (meta *TableMeta[T]) ExtractFieldValues(ins *T,point bool, exclude map[string]int) []any {
	var ret []any

	for _, column := range meta.Columns {
		if exclude != nil && exclude[column] != 0 {
			continue
		}
		ret = append(ret, meta.LookupFieldFunc(column, ins,point))
	}
	return ret
}


