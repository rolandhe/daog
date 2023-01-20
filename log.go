package daog

import (
	"encoding/json"
	"log"
)

type DaogLogErrorFunc func(tc *TransContext, err error)
type DaogLogInfoFunc func(tc *TransContext, content string)
type DaogLogExecSQLFunc func(tc *TransContext, sql string, args []any)

var (
	DaogLogError   DaogLogErrorFunc
	DaogLogInfo    DaogLogInfoFunc
	DaogLogExecSQL DaogLogExecSQLFunc
)

func init() {
	DaogLogError = func(tc *TransContext, err error) {
		log.Printf("tid=%s,err: %v\n", getTraceId(tc), err)
	}
	DaogLogInfo = func(tc *TransContext, content string) {
		log.Printf("tid=%s,content: %s\n", getTraceId(tc), content)
	}
	DaogLogExecSQL = func(tc *TransContext, sql string, args []any) {
		traceId := getTraceId(tc)
		jargs, err := json.Marshal(args)
		if err != nil {
			log.Printf("[Trace SQL] tid=%s, log sql but to json error:%v\n", traceId, err)
		}
		log.Printf("[Trace SQL] tid=%s,sql: %s, args:%s\n", traceId, sql, jargs)
	}
}

func getTraceId(tc *TransContext) string {
	values := tc.ctx.Value(CTXVALUES)
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
