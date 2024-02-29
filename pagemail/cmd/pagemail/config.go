package main

type AppConfig struct {
	Environment string `config:"app.environment"`
	LogLevel    string `config:"app.log-level"`
	Host        string `config:"app.host"`
	Port        string `config:"app.port"`
	DBPath      string `config:"db.path"`
}
