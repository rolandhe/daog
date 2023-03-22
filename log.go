// A quickly mysql access component.
//
// Copyright 2023 The daog Authors. All rights reserved.

package daog

import (
	"context"
	"github.com/rolandhe/daog/utils"
	"log"
)

// SQLLogger 封装log类似与java slf4j的功能
// 输出日志时你可以调用 GetTraceIdFromContext 从 context.Context 读取到traceId，并在输出日志时输出它，当然，前题是你在 NewTransContext 时指定了traceId
// 你也可以调用 GetGoroutineIdFromContext 读取到调用 NewTransContext 的 goroutine 的 goroutine id, 并在日志中输出它
type SQLLogger interface {
	// Error 输出错误日志
	Error (ctx context.Context, err error)
	// Info 输出 Info级别日志
	Info(ctx context.Context, content string)
	// ExecSQLBefore 在 sql 执行前输出它，这需要你在 构建数据源时指定的 DbConf.LogSQL = true
	ExecSQLBefore (ctx context.Context, sql string, argsJson []byte, sqlMd5 string)
	// ExecSQLAfter  在 sql 执行后输出它，这需要你在 构建数据源时指定的 DbConf.LogSQL = true
	ExecSQLAfter (ctx context.Context, sqlMd5 string, cost int64)
	// SimpleLogError 输出err，应用于你还没有构建 TransContext 前，此时 context.Context 还没有被初始化，对应err直接输出即可
	SimpleLogError (err error)
}

// GLogger 全局的日志接口对象, 您可以实现自己 SQLLogger 对象并赋值给 GLogger， 则可以自行输出日志
var GLogger SQLLogger = &defaultLogger{}

type  defaultLogger struct {

}

func (logger *defaultLogger)Error(ctx context.Context, err error)  {
	log.Printf("goid=%d,tid=%s,err: %v\n", GetGoroutineIdFromContext(ctx), GetTraceIdFromContext(ctx), err)
}

func (logger *defaultLogger)Info(ctx context.Context, content string)  {
	log.Printf("goid=%d,tid=%s,content: %s\n", GetGoroutineIdFromContext(ctx), GetTraceIdFromContext(ctx), content)
}

func (logger *defaultLogger)ExecSQLBefore(ctx context.Context, sql string, argsJson []byte, sqlMd5 string)  {
	traceId := GetTraceIdFromContext(ctx)
	log.Printf("[Trace SQL] goid=%d,tid=%s,sqlMd5=%s,sql: %s, args:%s\n", GetGoroutineIdFromContext(ctx), traceId, sqlMd5, sql, argsJson)
}
func (logger *defaultLogger) ExecSQLAfter (ctx context.Context, sqlMd5 string, cost int64){
	traceId := GetTraceIdFromContext(ctx)
	log.Printf("[Trace SQL] goid=%d,tid=%s,sqlMd5=%s,cost %d ms\n", GetGoroutineIdFromContext(ctx), traceId, sqlMd5, cost)
}

func (logger *defaultLogger)  SimpleLogError (err error){
	log.Printf("goid=%d, err: %v\n", utils.QuickGetGoroutineId(), err)
}


