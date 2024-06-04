package logging

import "context"

type contextKey string

const logContextKey = "logger"

func Get(ctx context.Context) *Logger {
	var lg *Logger
	lg, ok := ctx.Value(logContextKey).(*Logger)
	if !ok {
		lg = NewLogger("logger")
	}

	return lg
}

func Set(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, logContextKey, logger)
}

func NewContext(ctx context.Context, name string) context.Context {
	logger := NewLogger(name)
	return Set(ctx, logger)
}
