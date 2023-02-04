// Package dbtime,A quickly mysql access component.
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

type NormalDate time.Time

func (ndt NormalDate) Value() (driver.Value, error) {
	return time.Time(ndt), nil
}

func (ndt NormalDate)String() string{
	return time.Time(ndt).Format(DateFormat)
}

func (ndt *NormalDate) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Compare(b,nullJsonValue) == 0 {
		return nil
	}
	value := strings.Trim(string(b), `"`) //get rid of "
	t, err := time.Parse(DateFormat, value) //parse time
	if err != nil {
		daog.SimpleLogError(err)
		return err
	}
	*ndt = NormalDate(t) //set result using the pointer
	return nil
}

func (ndt NormalDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(ndt).Format(DateFormat) + `"`), nil
}
