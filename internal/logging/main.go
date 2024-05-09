package logging

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/mr55p-dev/htmx-utils/pkg/trace"
)

var logger *Logger = &Logger{slog.Default()}

// Callee stack depth to get out of when deferring a panic recover
const PANIC_STACK_DEPTH = 5 // tested on one function, might be different for different go versions

func SetHandler(h slog.Handler) {
	logger = &Logger{
		out: slog.New(h),
	}
}

func Info(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.out.InfoContext(ctx, msg, keyvals...)
}

func Debug(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.out.DebugContext(ctx, msg, keyvals...)
}

func Error(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.out.ErrorContext(ctx, msg, keyvals...)
}

func WithRequest(r *http.Request) *Logger {
	return logger.WithRequest(r)
}

func WithError(err error) *Logger {
	_, file, line, _ := runtime.Caller(1)
	return &Logger{
		out: logger.out.With("error", err.Error(), "file", file, "lineno", line),
	}
}

func WithRecover(r any) *Logger {
	return logger.WithRecover(r)
}

func Warn(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.out.WarnContext(ctx, msg, keyvals...)
}

type Logger struct {
	out *slog.Logger
}

func NewLogger(name string) *Logger {
	return &Logger{
		out: logger.out.With("name", name),
	}
}

func (l *Logger) WithRequest(r *http.Request) *Logger {
	return &Logger{
		out: l.out.With(
			"trace-id", trace.GetTrace(r.Context()),
		),
	}
}

func (l *Logger) WithError(err error) *Logger {
	_, file, line, _ := runtime.Caller(1)
	return &Logger{
		out: l.out.With("error", err.Error(), "file", file, "lineno", line),
	}
}

func (l *Logger) WithRecover(r any) *Logger {
	debug.PrintStack()
	return &Logger{
		out: l.out.With("recover", fmt.Sprintf("%v", r)),
	}
}

func (l *Logger) With(keyvals ...any) *Logger {
	return &Logger{
		out: l.out.With(keyvals...),
	}
}

func (l *Logger) InfoCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	l.out.InfoContext(ctx, msg, keyvals...)
}

func (l *Logger) DebugCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	l.out.DebugContext(ctx, msg, keyvals...)
}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	l.out.ErrorContext(ctx, msg, keyvals...)
}

func (l *Logger) WarnCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	l.out.WarnContext(ctx, msg, keyvals...)
}

func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.InfoCtx(context.Background(), msg, keyvals...)
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.DebugCtx(context.Background(), msg, keyvals...)
}

func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.ErrorCtx(context.Background(), msg, keyvals...)
}

func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.WarnCtx(context.Background(), msg, keyvals...)
}
