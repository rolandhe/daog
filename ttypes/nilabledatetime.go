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

type NilableDatetime struct {
	sql.NullTime
}

func FromDatetime(d time.Time) *NilableDatetime {
	return &NilableDatetime{
		sql.NullTime{d,true},
	}
}

func (ndt NilableDatetime)String() string{
	if !ndt.Valid {
		return "<nil>"
	}
	return ndt.Time.Format(DatetimeFormat)
}

func (s *NilableDatetime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		s.Valid = false
		return nil
	}
	if bytes.Compare(b,nullJsonValue) == 0{
		s.Valid = false
		return nil
	}
	value := strings.Trim(string(b), `"`) //get rid of "
	t, err := time.Parse(DatetimeFormat, value) //parse time
	if err != nil {
		daog.SimpleLogError(err)
		return err
	}
	s.Time = t
	s.Valid = true
	return nil
}

func (s NilableDatetime) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"),nil
	}
	return []byte(`"` + s.Time.Format(DatetimeFormat) + `"`), nil
}