package trace

import "context"

type traceKeyType string

const traceKey traceKeyType = "traceId"

func SetTrace(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, traceKey, traceId)
}

func GetTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value(traceKey).(string)
	if !ok {
		return ""
	}
	return traceId
}
