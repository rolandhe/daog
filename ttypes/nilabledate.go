// Package ttypes,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package ttypes

import (
	"bytes"
	"database/sql"
	"github.com/rolandhe/daog"
	"strings"
	"time"
)

type NilableDate struct {
	sql.NullTime
}

func FromDate(d time.Time) *NilableDate {
	return &NilableDate{
		sql.NullTime{d,true},
	}
}

func (ndt NilableDate)String() string{
	if !ndt.Valid {
		return "<nil>"
	}
	return ndt.Time.Format(DateFormat)
}

func (d *NilableDate) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		d.Valid = false
		return nil
	}
	if bytes.Compare(b,nullJsonValue) == 0{
		d.Valid = false
		return nil
	}
	value := strings.Trim(string(b), `"`) //get rid of "
	t, err := time.Parse(DateFormat, value) //parse time
	if err != nil {
		daog.SimpleLogError(err)
		return err
	}
	d.Time = t
	d.Valid = true
	return nil
}

func (d NilableDate) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"),nil
	}
	return []byte(`"` + d.Time.Format(DateFormat) + `"`), nil
}