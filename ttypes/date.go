// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

// Package ttypes 定义特殊类型支持，比如日期类型， golang的日期类型在转换成json时不能指定日期的格式，ttypes.NormalDate 可以按照 DateFormat指定的格式输出到json中或者fmt.Println.
// sql包中的NullTime也不支持格式和json格式输出，对应的NilableDateTime支持。本包包括:
//
// # NormalDate 支持格式和json输出的日期类型， 格式由 DateFormat 变量指定
//
// # NormalDatetime 支持格式和json输出的时间类型，格式由 DatetimeFormat 指定
//
// # NilableDate 可以为空的日期类型，扩展 sql.NullTime, 支持格式及json输出，格式由 DateFormat 变量指定
//
// # NilableDatetime 可以为空的时间类型，扩展 sql.NullTime, 支持格式及json输出 ，格式由 DatetimeFormat 指定
//
// NilableString， 可以为空的String类型，扩展 sql.NullString，支持json输出，提供与string的转换及其他字符串操作功能
package ttypes

import (
	"bytes"
	"database/sql/driver"
	"strings"
	"time"
)

// NormalDate 支持按日期格式输出的日期类型, 格式由 DateFormat 全局变量指定, 实现 fmt.Stringer, driver.Valuer, json.Unmarshaler, json.Marshaler 接口
type NormalDate time.Time

func ParseNormalDate(sDate string) (*NormalDate, error) {
	t, err := time.ParseInLocation(DateFormat, sDate, time.Local)
	if err != nil {
		return nil, err
	}
	ndt := NormalDate(t)
	return &ndt, nil
}

func ParseNormalDateInUTC(sDate string) (*NormalDate, error) {
	t, err := time.Parse(DateFormat, sDate)
	if err != nil {
		return nil, err
	}
	ndt := NormalDate(t)
	return &ndt, nil
}

func ParseNormalDateInLocal(sDate string, loc *time.Location) (*NormalDate, error) {
	t, err := time.ParseInLocation(DateFormat, sDate, loc)
	if err != nil {
		return nil, err
	}
	ndt := NormalDate(t)
	return &ndt, nil
}

func (ndt *NormalDate) MonthFmtString() string {
	return ndt.ToTimePointer().Format("2006-01")
}

// Value 实现 driver.Valuer
func (ndt NormalDate) Value() (driver.Value, error) {
	return *ndt.ToTimePointer(), nil
}

// String 实现 fmt.Stringer 接口
func (ndt NormalDate) String() string {
	return ndt.ToTimePointer().Format(DateFormat)
}

// UnmarshalJSON 实现 json.Unmarshaler
func (ndt *NormalDate) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Compare(b, nullJsonValue) == 0 {
		return nil
	}
	value := strings.Trim(string(b), `"`)                         //get rid of "
	t, err := time.ParseInLocation(DateFormat, value, time.Local) //parse time
	if err != nil {
		return err
	}
	*ndt = NormalDate(t) //set result using the pointer
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (ndt NormalDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ndt.ToTimePointer().Format(DateFormat) + `"`), nil
}
func (ndt *NormalDate) ToTimePointer() *time.Time {
	return (*time.Time)(ndt)
}
