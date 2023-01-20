package entities

import (
	"github.com/roland/daog"
	dbtime "github.com/roland/daog/time"
	"github.com/shopspring/decimal"
)

var GroupInfoFields = struct {
	Id          string
	Name        string
	MainData    string
	CreateAt    string
	TotalAmount string
}{
	"id",
	"name",
	"main_data",
	"create_at",
	"total_amount",
}

var GroupInfoMeta = &daog.TableMeta[GroupInfo]{
	InstanceFunc: func() *GroupInfo {
		return &GroupInfo{}
	},
	Table: "group_info",
	Columns: []string{
		"id",
		"name",
		"main_data",
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
}

type GroupInfo struct {
	Id          int64
	Name        string
	MainData    string
	CreateAt    dbtime.NormalDatetime
	TotalAmount decimal.Decimal
}
