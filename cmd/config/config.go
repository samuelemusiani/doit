package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Sever struct {
	Listen string
}

type Databse struct {
	Path string
}

type Log struct {
	Log_level string
}

type Config struct {
	Server  Sever
	Databse Databse
	Log     Log
}

var config Config = Config{
	Server: Sever{
		Listen: "0.0.0.0:8080",
	},
	Databse: Databse{
		Path: "./doit.db",
	},
	Log: Log{
		Log_level: "info",
	},
}

func ParseConfig(path string) error {
	configBuff, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(configBuff, &config)
	if err != nil {
		return err
	}

	return nil
}

func GetConfig() *Config {
	return &config
}
