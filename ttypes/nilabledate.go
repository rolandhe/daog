// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

package ttypes

import (
	"bytes"
	"database/sql"
	"github.com/rolandhe/daog"
	"strings"
	"time"
)

// NilableDate 可以为空的日期类型，扩展NullTime, 支持格式及json输出，格式由 DateFormat 变量指定，实现 fmt.Stringer,  json.Unmarshaler, json.Marshaler 接口
type NilableDate struct {
	sql.NullTime
}

// FromDate 转换time.Time类型为 NilableDate 类型，返回值为指针，如果接收变量为 NilableDate 类型，需要使用加*号来解引用: *nilableDate
func FromDate(d time.Time) *NilableDate {
	return &NilableDate{
		sql.NullTime{d, true},
	}
}

// String 实现 fmt.Stringer 接口
func (ndt NilableDate) String() string {
	if !ndt.Valid {
		return "<nil>"
	}
	return ndt.Time.Format(DateFormat)
}

// UnmarshalJSON 实现 json.Unmarshaler
func (d *NilableDate) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		d.Valid = false
		return nil
	}
	if bytes.Compare(b, nullJsonValue) == 0 {
		d.Valid = false
		return nil
	}
	value := strings.Trim(string(b), `"`)   //get rid of "
	t, err := time.Parse(DateFormat, value) //parse time
	if err != nil {
		daog.GLogger.SimpleLogError(err)
		return err
	}
	d.Time = t
	d.Valid = true
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (d NilableDate) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format(DateFormat) + `"`), nil
}
