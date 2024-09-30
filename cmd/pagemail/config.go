package main

import (
	"os"

	"github.com/mr55p-dev/gonk"
)

type Config struct {
	App struct {
		Host      string
		Port      int
		CookieKey string
	}
	Mail struct {
		Host     string
		Port     int
		Username string
		Password string
		PoolSize int
	}
	DB struct {
		DSN string `config:"dsn,optional"`
	} `config:"db"`
}

func MustLoadConfig() *Config {
	config := new(Config)
	yamlLoader, err := gonk.NewYamlLoader("pagemail.yaml")
	if err != nil {
		LogError("Failed to open pagemail.yaml", err)
	}
	err = gonk.LoadConfig(config, yamlLoader, gonk.EnvLoader("PM"))
	if err != nil {
		PanicError("Failed to load config", err)
	}
	if config.DB.DSN == "" {
		logger.Info("Attempting to load databse DSN from DATABASE_URL env")
		config.DB.DSN = os.Getenv("DATABASE_URL")
	}

	return config
}
