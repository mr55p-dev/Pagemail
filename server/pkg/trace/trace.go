package trace

import "context"

type traceKey string

var key traceKey = "trace-id"
var TraceHeader = "X-Trace-Id"

func SetTrace(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, traceKey("trace-id"), val)
}

func GetTrace(ctx context.Context) string {
	val, _ := ctx.Value(key).(string)
	return val
}
