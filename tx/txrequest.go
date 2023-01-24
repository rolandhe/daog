// Package txrequest,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package txrequest

type RequestStyle int

const (
	RequestNone     = RequestStyle(0)
	RequestReadonly = RequestStyle(1)
	RequestWrite    = RequestStyle(2)
)
