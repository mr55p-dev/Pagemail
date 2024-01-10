package logging

import (
	"context"
	"log/slog"
	"os"
	"runtime"
)

type Log struct{ *slog.Logger }
type LogKey string

const (
	UserId   = "user-id"
	UserMail = "user-email"
	PageId   = "page-id"
	Error    = "error"
	File     = "file"
	Line     = "lineno"
)

var BaseLog *slog.Logger

func init() {
	BaseLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

func GetLogger(name string) Log {
	return Log{BaseLog.With("module", name)}
}

func (l *Log) Err(msg string, err error, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	callerArgs := []any{Error, err.Error(), File, file, Line, line}
	callerArgs = append(callerArgs, args...)
	l.Error(msg, callerArgs...)
}

func (l *Log) ErrContext(ctx context.Context, msg string, err error, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	callerArgs := []any{Error, err.Error(), File, file, Line, line}
	callerArgs = append(callerArgs, args...)
	l.ErrorContext(ctx, msg, callerArgs...)
}
