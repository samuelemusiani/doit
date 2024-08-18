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

type FirstUser struct {
	Username string
	Email    string
}

type Users struct {
	First_User FirstUser
}

type Config struct {
	Server  Sever
	Databse Databse
	Log     Log
	Users   Users
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
	Users: Users{
		First_User: FirstUser{
			Username: "admin",
			Email:    "admin@mail.com",
		},
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
