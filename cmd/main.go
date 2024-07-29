package main

import (
	"log/slog"

	"github.com/samuelemusiani/doit/cmd/config"
	"github.com/samuelemusiani/doit/cmd/db"
)

func main() {
	slog.Info("Starting DOIT")
	config.ParseConfig("config.toml")

	db.Init()
	db.Close()

	slog.Info("Terminating DOIT")
}
