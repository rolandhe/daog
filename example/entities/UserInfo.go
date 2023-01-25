package entities

import (
    "github.com/rolandhe/daog"
    dbtime "github.com/rolandhe/daog/time"
    "github.com/shopspring/decimal"
)

var UserInfoFields = struct {
   Id string
   Name string
   Data string
   CreateAt string
   Amount string
   
}{
    "id",
    "name",
    "data",
    "create_at",
    "amount",
    
}

var  UserInfoMeta = &daog.TableMeta[UserInfo]{
    Table: "user_info",
    Columns: []string {
        "id",
        "name",
        "data",
        "create_at",
        "amount",
        
    },
    AutoColumn: "id",
    LookupFieldFunc: func(columnName string,ins *UserInfo,point bool) any {
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
        if "data" == columnName {
            if point {
                 return &ins.Data
            }
            return ins.Data
        }
        if "create_at" == columnName {
            if point {
                 return &ins.CreateAt
            }
            return ins.CreateAt
        }
        if "amount" == columnName {
            if point {
                 return &ins.Amount
            }
            return ins.Amount
        }
        
        return nil
    },
}


type UserInfo struct {
    Id int64
    Name string
    Data string
    CreateAt dbtime.NormalDatetime
    Amount decimal.Decimal
    
}
