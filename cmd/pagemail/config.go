package main

import (
	"log/slog"

	"github.com/mr55p-dev/pagemail/internal/tools"
)

type AppConfig struct {
	Environment    string `config:"app.environment"`
	LogLevel       string `config:"app.log-level" log:"logLevel"`
	Host           string `config:"app.host" log:"host"`
	DBPath         string `config:"db.path" log:"db-path"`
	CookieKeyFile  string `config:"app.cookie-key-file" log:"cookie-key-file"`
	GoogleClientId string `config:"app.google-client-id" log:"google-client-id"`

	External struct {
		Host  string `config:"host" log:"extern-host"`
		Proto string `config:"proto" log:"extern-proto"`
	} `config:"extern"`
}

func (config *AppConfig) LogValue() slog.Value {
	vals := tools.LogValue(config)
	return slog.GroupValue(vals...)
}
