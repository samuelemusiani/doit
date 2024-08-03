package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	ListeningAddress string
	ListeningPort    uint16
	DBPath           string
}

var config Config = Config{
	ListeningAddress: "0.0.0.0",
	ListeningPort:    8080,
	DBPath:           "./doit.db",
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
