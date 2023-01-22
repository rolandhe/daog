package daog

import (
	"context"
	"log"
)

type LogErrorFunc func(ctx context.Context, err error)
type LogInfoFunc func(ctx context.Context, content string)
type LogExecSQLFunc func(ctx context.Context, sql string, argsJson []byte, sqlMd5 string)
type LogExecSQLAfterFunc func(ctx context.Context, sqlMd5 string, cost int64)

var (
	LogError        LogErrorFunc
	LogInfo         LogInfoFunc
	LogExecSQL      LogExecSQLFunc
	LogExecSQLAfter LogExecSQLAfterFunc
)

func init() {
	LogError = func(ctx context.Context, err error) {
		log.Printf("tid=%s,err: %v\n", GetTraceIdFromContext(ctx), err)
	}
	LogInfo = func(ctx context.Context, content string) {
		log.Printf("tid=%s,content: %s\n", GetTraceIdFromContext(ctx), content)
	}
	LogExecSQL = func(ctx context.Context, sql string, argJson []byte, sqlMd5 string) {
		traceId := GetTraceIdFromContext(ctx)
		log.Printf("[Trace SQL] tid=%s,sqlMd5=%s,sql: %s, args:%s\n", traceId, sqlMd5, sql, string(argJson))
	}
	LogExecSQLAfter = func(ctx context.Context, sqlMd5 string, cost int64) {
		traceId := GetTraceIdFromContext(ctx)
		log.Printf("[Trace SQL] tid=%s,sqlMd5=%s,cost %d ms\n", traceId, sqlMd5, cost)
	}
}
