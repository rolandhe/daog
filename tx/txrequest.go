// A quickly mysql access component.
// Copyright 2023 The daog Authors. All rights reserved.

// Package txrequest, 定义了事务的级别，包括三个级别：
//
// # RequestNone 没有特别的事务要求，每一条sql的执行都在mysql缺省事务中
//
// # RequestReadonly 只读事务, 事务内只有读取数据sql，不能有写数据的操作，可以有效的提升性能
//
// RequestWrite 写事务，事务内支持写操作，当然也支持读操作
package txrequest

type RequestStyle int

const (
	RequestNone     = RequestStyle(0)
	RequestReadonly = RequestStyle(1)
	RequestWrite    = RequestStyle(2)
)
