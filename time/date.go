// Package dbtime,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package dbtime

import (
	"database/sql/driver"
	"strings"
	"time"
)

type NormalDate time.Time

func (ndt NormalDate) Value() (driver.Value, error) {
	return time.Time(ndt), nil
}

func (ndt *NormalDate) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(DateFormat, value) //parse time
	if err != nil {
		return err
	}
	*ndt = NormalDate(t) //set result using the pointer
	return nil
}

func (ndt NormalDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(ndt).Format(DateFormat) + `"`), nil
}
