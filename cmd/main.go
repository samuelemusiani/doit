package main

import (
	"log/slog"

	"github.com/samuelemusiani/doit/cmd/config"
	"github.com/samuelemusiani/doit/cmd/db"
	"github.com/samuelemusiani/doit/cmd/http_server"
)

func main() {
	slog.Info("Starting DOIT")
	config.ParseConfig("config.toml")

	slog.SetLogLoggerLevel(slog.LevelDebug)

	db.Init()
	http_server.Init()
	http_server.ListenAndServe()

	db.Close()
	slog.Info("Terminating DOIT")
}
