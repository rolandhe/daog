// Package dbtime,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package dbtime

import (
	"database/sql/driver"
	"strings"
	"time"
)

type NormalDatetime time.Time

func (ndt NormalDatetime) Value() (driver.Value, error) {
	return time.Time(ndt), nil
}
func (ndt *NormalDatetime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(DatetimeFormat, value) //parse time
	if err != nil {
		return err
	}
	*ndt = NormalDatetime(t) //set result using the pointer
	return nil
}

func (ndt NormalDatetime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(ndt).Format(DatetimeFormat) + `"`), nil
}
