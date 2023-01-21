package daog

import (
	"context"
	"encoding/json"
	"log"
)

type DaogLogErrorFunc func(ctx context.Context, err error)
type DaogLogInfoFunc func(ctx context.Context, content string)
type DaogLogExecSQLFunc func(ctx context.Context, sql string, args []any)

var (
	DaogLogError   DaogLogErrorFunc
	DaogLogInfo    DaogLogInfoFunc
	DaogLogExecSQL DaogLogExecSQLFunc
)

func init() {
	DaogLogError = func(ctx context.Context, err error) {
		log.Printf("tid=%s,err: %v\n", getTraceId(ctx), err)
	}
	DaogLogInfo = func(ctx context.Context, content string) {
		log.Printf("tid=%s,content: %s\n", getTraceId(ctx), content)
	}
	DaogLogExecSQL = func(ctx context.Context, sql string, args []any) {
		traceId := getTraceId(ctx)
		jargs, err := json.Marshal(args)
		if err != nil {
			log.Printf("[Trace SQL] tid=%s, log sql but to json error:%v\n", traceId, err)
		}
		log.Printf("[Trace SQL] tid=%s,sql: %s, args:%s\n", traceId, sql, jargs)
	}
}

func getTraceId(ctx context.Context) string {
	values := ctx.Value(CTXVALUES)
	if values == nil {
		return ""
	}

	v, ok := values.(map[string]any)
	if !ok {
		return ""
	}
	data, ok := v[TRACEID]
	if !ok {
		return ""
	}
	trace, ok := data.(string)
	if !ok {
		return ""
	}
	return trace
}
