// Package daog,A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.
package daog

import (
	"context"
	"log"
)

type LogErrorFunc func(ctx context.Context, err error)
type LogInfoFunc func(ctx context.Context, content string)
type LogExecSQLBeforeFunc func(ctx context.Context, sql string, argsJson []byte, sqlMd5 string)
type LogExecSQLAfterFunc func(ctx context.Context, sqlMd5 string, cost int64)

var (
	LogError         LogErrorFunc
	LogInfo          LogInfoFunc
	LogExecSQLBefore LogExecSQLBeforeFunc
	LogExecSQLAfter  LogExecSQLAfterFunc
)

func init() {
	LogError = func(ctx context.Context, err error) {
		log.Printf("goid=%d,tid=%s,err: %v\n", GetGoRoutineIdFromContext(ctx), GetTraceIdFromContext(ctx), err)
	}
	LogInfo = func(ctx context.Context, content string) {
		log.Printf("goid=%d,tid=%s,content: %s\n", GetGoRoutineIdFromContext(ctx), GetTraceIdFromContext(ctx), content)
	}
	LogExecSQLBefore = func(ctx context.Context, sql string, argJson []byte, sqlMd5 string) {
		traceId := GetTraceIdFromContext(ctx)
		log.Printf("[Trace SQL] goid=%d,tid=%s,sqlMd5=%s,sql: %s, args:%s\n", GetGoRoutineIdFromContext(ctx), traceId, sqlMd5, sql, argJson)
	}
	LogExecSQLAfter = func(ctx context.Context, sqlMd5 string, cost int64) {
		traceId := GetTraceIdFromContext(ctx)
		log.Printf("[Trace SQL] goid=%d,tid=%s,sqlMd5=%s,cost %d ms\n", GetGoRoutineIdFromContext(ctx), traceId, sqlMd5, cost)
	}
}
