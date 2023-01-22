package daog

import (
	"context"
	"log"
)

type DaogLogErrorFunc func(ctx context.Context, err error)
type DaogLogInfoFunc func(ctx context.Context, content string)
type DaogLogExecSQLFunc func(ctx context.Context, sql string, argsJson []byte, sqlMd5 string)
type DaogLogExecSQLAfterFunc func(ctx context.Context, sqlMd5 string, cost int64)

var (
	DaogLogError        DaogLogErrorFunc
	DaogLogInfo         DaogLogInfoFunc
	DaogLogExecSQL      DaogLogExecSQLFunc
	DaogLogExecSQLAfter DaogLogExecSQLAfterFunc
)

func init() {
	DaogLogError = func(ctx context.Context, err error) {
		log.Printf("tid=%s,err: %v\n", GetTraceIdFromContext(ctx), err)
	}
	DaogLogInfo = func(ctx context.Context, content string) {
		log.Printf("tid=%s,content: %s\n", GetTraceIdFromContext(ctx), content)
	}
	DaogLogExecSQL = func(ctx context.Context, sql string, argJson []byte, sqlMd5 string) {
		traceId := GetTraceIdFromContext(ctx)
		log.Printf("[Trace SQL] tid=%s,sqlMd5=%s,sql: %s, args:%s\n", traceId, sqlMd5, sql, string(argJson))
	}
	DaogLogExecSQLAfter = func(ctx context.Context, sqlMd5 string, cost int64) {
		traceId := GetTraceIdFromContext(ctx)
		log.Printf("[Trace SQL] tid=%s,sqlMd5=%s,cost %d ms\n", traceId, sqlMd5, cost)
	}
}
