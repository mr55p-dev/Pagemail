package pagemailLog

import "log/slog"

type Logger slog.Logger

func New(name string) *slog.Logger {
	h := slog.NewHandler
	return slog.New()
}
