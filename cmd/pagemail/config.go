package main

import "github.com/mr55p-dev/gonk"

type Config struct {
	App struct {
		Host          string
		Port          int
		CookieKeyFile string
	}
	Mail struct {
		Host     string
		Port     int
		Username string
		Password string
		PoolSize int
	}
	Db struct {
		Path string
	}
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

	return config
}
