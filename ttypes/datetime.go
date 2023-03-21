// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

package ttypes

import (
	"bytes"
	"database/sql/driver"
	"github.com/rolandhe/daog"
	"strings"
	"time"
)

var (
	// DatetimeFormat 指定时间格式
	DatetimeFormat = "2006-01-02 15:04:05"
)

// NormalDatetime 支持按日期格式输出的日期类型, 格式由 DatetimeFormat 全局变量指定, 实现fmt.Stringer, driver.Valuer, json.Unmarshaler, json.Marshaler 接口
type NormalDatetime time.Time

// Value 实现 driver.Valuer
func (ndt NormalDatetime) Value() (driver.Value, error) {
	return time.Time(ndt), nil
}

// String 实现 fmt.Stringer 接口
func (ndt NormalDatetime) String() string {
	return time.Time(ndt).Format(DatetimeFormat)
}

// UnmarshalJSON 实现 json.Unmarshaler
func (ndt *NormalDatetime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Compare(b, nullJsonValue) == 0 {
		return nil
	}

	value := strings.Trim(string(b), `"`)       //get rid of "
	t, err := time.Parse(DatetimeFormat, value) //parse time
	if err != nil {
		daog.SimpleLogError(err)
		return err
	}
	*ndt = NormalDatetime(t) //set result using the pointer
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (ndt NormalDatetime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(ndt).Format(DatetimeFormat) + `"`), nil
}
