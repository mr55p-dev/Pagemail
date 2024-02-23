package logging

import (
	"context"
	"log/slog"
)

type Logger struct {
	*slog.Logger
}

func (l Logger) Err(msg string, err error) {
	l.Error(msg, "error", err.Error())
}

func (l Logger) Errc(ctx context.Context, msg string, err error) {
	l.ErrorContext(ctx, msg, "error", err.Error())
}

func New(log *slog.Logger) *Logger {
	return &Logger{
		Logger: log,
	}
}
