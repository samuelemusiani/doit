package db

import (
	"database/sql"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/samuelemusiani/doit/cmd/config"
)

var global_db *SQLiteRepository = nil

func Init() {
	config := config.GetConfig()

	rawDB, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		slog.Error("Opening DB at path %s: %w", config.DBPath, err)
		log.Fatal("Could not continue")
	}

	global_db = NewSQLiteRepository(rawDB)
	err = global_db.Migrate()
	if err != nil {
		slog.Error("Migrating DB: %w", err)
		log.Fatal("Could not continue")
	}
}

func Close() {
	err := global_db.db.Close()
	if err != nil {
		slog.Error("Closing db: %w", err)
	}
}
