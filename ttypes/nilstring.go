// Package ttypes,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package ttypes

import (
	"bytes"
	"database/sql"
	"strings"
)

var nullJsonValue = []byte("null")

type NilableString struct {
	sql.NullString
}

func FromString(s string) *NilableString {
	return &NilableString{
		sql.NullString{s,true},
	}
}


func (s *NilableString) StringNilAsEmpty() string{
	if s.Valid {
		return s.String
	}
	return ""
}

func (s *NilableString) StringNilAsDefault(def string) string{
	if s.Valid {
		return s.String
	}
	return def
}

func (s *NilableString) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		s.Valid = false
		return nil
	}
	if bytes.Compare(b,nullJsonValue) == 0{
		s.Valid = false
		s.String = ""
		return nil
	}
	s.String = strings.Trim(string(b), `"`) //get rid of "
	s.Valid = true
	return nil
}

func (s NilableString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"),nil
	}
	return []byte(`"` + s.String + `"`),nil
}

