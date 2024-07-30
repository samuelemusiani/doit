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
	DBPath:           "./doit.db",
}

func ParseConfig(path string) {
	configBuff, err := os.ReadFile(path)
	if err != nil {
		slog.With("err", err).Error("Error during config reading")
		slog.With("config", config).Warn("Using default config values")
		return
	}

	err = toml.Unmarshal(configBuff, &config)
	if err != nil {
		slog.With("err", err).Error("Error during config parsing")
		slog.Warn("Stopping now to avoid confusion. Either provide a valid config or nothing at all :)")
		os.Exit(1)
	}
}

func GetConfig() *Config {
	return &config
}
