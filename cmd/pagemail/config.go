package main

import "github.com/mr55p-dev/gonk"

type Config struct {
	App struct {
		Host          string `config:"host"`
		Port          int    `config:"port"`
		CookieKeyFile string `config:"cookie-key-file"`
	} `config:"app"`
	Mail struct {
		Host     string `config:"host"`
		Port     int    `config:"port"`
		Username string `config:"username"`
		Password string `config:"password"`
		PoolSize int    `config:"pool-size"`
	} `config:"mail"`
	DB struct {
		Path string `config:"path"`
	} `config:"db"`
}

func MustLoadConfig() *Config {
	config := new(Config)
	yamlLoader, err := gonk.NewYamlLoader("pagemail.yaml")
	if err != nil {
		PanicError("Failed to open pagemail.yaml", err)
	}
	err = gonk.LoadConfig(config, yamlLoader)
	if err != nil {
		PanicError("Failed to load config", err)
	}

	return config
}
