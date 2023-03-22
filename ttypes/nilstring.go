// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

package ttypes

import (
	"bytes"
	"database/sql"
)

var nullJsonValue = []byte("null")

// NilableString 可以为空的string类型，扩展 sql.NullString, 支持友好的json输出，提供常用的操作函数
type NilableString struct {
	sql.NullString
}

// FromString 转换string类型为 NilableString 类型，返回值为指针，如果接收变量为 NilableString 类型，需要使用加*号来解引用: *nilableString
func FromString(s string) *NilableString {
	return &NilableString{
		sql.NullString{s, true},
	}
}

// StringNilAsEmpty 返回当前对象包含的string值，如果当前对象是nil，则返回""字符串
func (s *NilableString) StringNilAsEmpty() string {
	return s.StringNilAsDefault("")
}

// StringNilAsDefault 返回当前对象包含的string值，如果当前对象是nil，返回 def 参数指定的字符串
func (s *NilableString) StringNilAsDefault(def string) string {
	if s.Valid {
		return s.String
	}
	return def
}

// UnmarshalText 实现 encoding.TextUnmarshaler 接口
func (s *NilableString) UnmarshalText(b []byte) error{
	if len(b) == 0 {
		s.Valid = false
		return nil
	}
	if bytes.Compare(b, nullJsonValue) == 0 {
		s.Valid = false
		s.String = ""
		return nil
	}
	//s.String = strings.Trim(string(b), `"`) //get rid of "
	s.String = string(b)
	s.Valid = true
	return nil
}

// MarshalText 实现 encoding.TextMarshaler 接口
func (s NilableString) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}

	return []byte(s.String), nil
}
