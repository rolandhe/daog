package dal

import (
	"github.com/rolandhe/daog"
)

var BitsSampleFields = struct {
	Id     string
	V      string
	Status string
}{
	"id",
	"v",
	"status",
}

var BitsSampleMeta = &daog.TableMeta[BitsSample]{
	Table: "bits_sample",
	Columns: []string{
		"id",
		"v",
		"status",
	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string, ins *BitsSample, point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "v" == columnName {
			if point {
				return &ins.V
			}
			return ins.V
		}
		if "status" == columnName {
			if point {
				return &ins.Status
			}
			return ins.Status
		}

		return nil
	},
	StampColumns: nil,
}

var BitsSampleDao daog.QuickDao[BitsSample] = &struct {
	daog.QuickDao[BitsSample]
}{
	daog.NewBaseQuickDao(BitsSampleMeta),
}

type BitsSample struct {
	Id     int64
	V      int32
	Status int32
}
