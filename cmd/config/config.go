package config

import (
	"github.com/pelletier/go-toml/v2"
	"log/slog"
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
	DBPath:           "./db",
}

func ParseConfig(path string) {
	configBuff, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Error during config reading. Could not continue: %w", err)
		os.Exit(1)
	}

	err = toml.Unmarshal(configBuff, &config)
	if err != nil {
		slog.Error("Error during config parsing. Could not continue: %w", err)
		os.Exit(1)
	}
}

func GetConfig() *Config {
	return &config
}
