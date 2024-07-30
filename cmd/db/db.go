package db

import (
	"database/sql"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/samuelemusiani/doit/cmd/config"
	"github.com/samuelemusiani/doit/cmd/doit"
)

var global_db *SQLiteRepository = nil

func Init() {
	slog.Debug("Init db connection")

	config := config.GetConfig()

	rawDB, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		slog.Error("Opening DB at path %s: %w", config.DBPath, err)
		log.Fatal("Could not continue")
	}

	global_db = newSQLiteRepository(rawDB)
	err = global_db.migrate()
	if err != nil {
		slog.Error("Migrating DB: %w", err)
		log.Fatal("Could not continue")
	}
}

func Close() {
	slog.Debug("Closing DB connection")

	err := global_db.db.Close()
	if err != nil {
		slog.Error("Closing db: %w", err)
	}
}

func CreateNote(note doit.Note) (*doit.Note, error) {
	return global_db.createNote(note)
}

func AllNotes() ([]doit.Note, error) {
	return global_db.allNotes()
}

func GetNoteByID(id int64) (*doit.Note, error) {
	return global_db.getNoteByID(id)
}

func DeleteNoteByID(id int64) error {
	return global_db.deleteNoteByID(id)
}
