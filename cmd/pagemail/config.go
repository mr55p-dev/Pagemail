package main

import "github.com/mr55p-dev/gonk"

type Config struct {
	App struct {
		Host string `config:"host"`
	} `config:"app"`
	DB struct {
		Path string `config:"path"`
	} `config:"path"`
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
