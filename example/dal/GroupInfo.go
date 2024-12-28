package dal

import (
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
	"github.com/shopspring/decimal"
)

var GroupInfoFields = struct {
	Id          string
	Name        string
	MainData    string
	Content     string
	BinData     string
	CreateAt    string
	TotalAmount string
}{
	"id",
	"name",
	"main_data",
	"content",
	"bin_data",
	"create_at",
	"total_amount",
}

var GroupInfoMeta = &daog.TableMeta[GroupInfo]{
	Table: "group_info",
	Columns: []string{
		"id",
		"name",
		"main_data",
		"content",
		"bin_data",
		"create_at",
		"total_amount",
	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string, ins *GroupInfo, point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "name" == columnName {
			if point {
				return &ins.Name
			}
			return ins.Name
		}
		if "main_data" == columnName {
			if point {
				return &ins.MainData
			}
			return ins.MainData
		}
		if "content" == columnName {
			if point {
				return &ins.Content
			}
			return ins.Content
		}
		if "bin_data" == columnName {
			if point {
				return &ins.BinData
			}
			return ins.BinData
		}
		if "create_at" == columnName {
			if point {
				return &ins.CreateAt
			}
			return ins.CreateAt
		}
		if "total_amount" == columnName {
			if point {
				return &ins.TotalAmount
			}
			return ins.TotalAmount
		}

		return nil
	},
	StampColumns: nil,
}

var GroupInfoDao daog.QuickDao[GroupInfo] = &struct {
	daog.QuickDao[GroupInfo]
}{
	daog.NewBaseQuickDao(GroupInfoMeta),
}

type GroupInfo struct {
	Id          int64
	Name        string
	MainData    string
	Content     string
	BinData     []byte
	CreateAt    ttypes.NormalDatetime
	TotalAmount decimal.Decimal
}
