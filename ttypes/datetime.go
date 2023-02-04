// Package ttypes,A quickly mysql access component.
//
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
	DatetimeFormat = "2006-01-02 15:04:05"
)

type NormalDatetime time.Time

func (ndt NormalDatetime) Value() (driver.Value, error) {
	return time.Time(ndt), nil
}

func (ndt NormalDatetime)String() string{
	return time.Time(ndt).Format(DatetimeFormat)
}

func (ndt *NormalDatetime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Compare(b,nullJsonValue) == 0 {
		return nil
	}

	value := strings.Trim(string(b), `"`) //get rid of "
	t, err := time.Parse(DatetimeFormat, value) //parse time
	if err != nil {
		daog.SimpleLogError(err)
		return err
	}
	*ndt = NormalDatetime(t) //set result using the pointer
	return nil
}

func (ndt NormalDatetime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(ndt).Format(DatetimeFormat) + `"`), nil
}
