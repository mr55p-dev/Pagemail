package logging

import (
	"context"
	"log/slog"
)

var logger *slog.Logger

func init() {
	logger = slog.Default()
}

func SetHandler(h slog.Handler) {
	logger = slog.New(h)
}

func Info(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.InfoContext(ctx, msg, keyvals...)
}

func Debug(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.DebugContext(ctx, msg, keyvals...)
}

func Error(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.ErrorContext(ctx, msg, keyvals...)
}

func Warn(ctx context.Context, msg string, keyvals ...interface{}) {
	logger.WarnContext(ctx, msg, keyvals...)
}

type Logger struct {
	logger *slog.Logger
}

func NewLogger(name string) *Logger {
	return &Logger{
		logger: logger.With("name", name),
	}
}

func (l *Logger) InfoCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	l.logger.InfoContext(ctx, msg, keyvals...)
}

func (l *Logger) DebugCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	l.logger.DebugContext(ctx, msg, keyvals...)
}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	l.logger.ErrorContext(ctx, msg, keyvals...)
}

func (l *Logger) WarnCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	l.logger.WarnContext(ctx, msg, keyvals...)
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
