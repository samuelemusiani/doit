package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"os"

	"github.com/samuelemusiani/doit/cmd/config"
	"github.com/samuelemusiani/doit/cmd/db"

	"github.com/samuelemusiani/doit/cmd/http_server"
)

//go:embed all:_front
var front_fs embed.FS

func main() {
	slog.Info("Starting DOIT")
	err := config.ParseConfig("config.toml")
	conf := config.GetConfig()
	if err != nil {
		slog.With("err", err).Error("Error during config reading")
		slog.With("config", conf).Warn("Using default config values")
		os.Exit(1)
	}

	var logLevl slog.Level
	switch conf.Log.Log_level {
	case "debug":
		logLevl = slog.LevelDebug
	case "info":
		logLevl = slog.LevelInfo
	case "warn":
		logLevl = slog.LevelWarn
	case "error":
		logLevl = slog.LevelError
	default:
		logLevl = slog.LevelInfo
	}
	slog.SetLogLoggerLevel(logLevl)

	err = db.Init()
	if err != nil {
		slog.With("path", conf.Databse.Path, "err", err).Error("Initializing database")
		os.Exit(1)
	}

	front_fs, err := fs.Sub(front_fs, "front")
	if err != nil {
		slog.With("err", err).Error("Initializing change base path for front fs")
		os.Exit(1)
	}

	http_server.Init(front_fs)
	err = http_server.ListenAndServe()
	if err != nil {
		slog.With("err", err).Error("Listening and serving")
		os.Exit(1)
	}

	err = db.Close()
	if err != nil {
		slog.With("err", err).Error("Closing database")
	}

	slog.Info("Terminating DOIT")
}
