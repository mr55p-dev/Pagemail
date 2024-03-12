package main

import (
	"log/slog"

	"github.com/mr55p-dev/pagemail/internal/tools"
)

type AppConfig struct {
	Environment string `config:"app.environment,required" log:"environment"`
	LogLevel    string `config:"app.log-level" log:"logLevel"`
	Host        string `config:"app.host" log:"host"`
	Port        int    `config:"app.port" log:"port"`
	DBPath      string `config:"db.path" log:"db-path"`
}

func (config *AppConfig) LogValue() slog.Value {
	vals := tools.LogValue(config)
	return slog.GroupValue(vals...)
}
