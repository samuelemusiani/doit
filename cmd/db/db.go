package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/samuelemusiani/doit/cmd/config"
	"github.com/samuelemusiani/doit/cmd/doit"
)

var global_db *SQLiteRepository = nil

func Init() error {
	slog.Debug("Init db connection")

	config := config.GetConfig()

	rawDB, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		return err
	}

	global_db = newSQLiteRepository(rawDB)
	err = global_db.migrate()
	if err != nil {
		return err
	}
	return nil
}

func Close() error {
	slog.Debug("Closing DB connection")
	err := global_db.db.Close()
	return err
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

func CreateUser(user doit.User) (*doit.User, error) {
	return global_db.createUser(user)
}

func AllUsers() ([]doit.User, error) {
	return global_db.allUsers()
}

func GetUserById(id int64) (*doit.User, error) {
	return global_db.getUserByID(id)
}

func GetUserByUsername(username string) (*doit.User, error) {
	return global_db.getUserByUsername(username)
}

func GetUserByEmail(email string) (*doit.User, error) {
	return global_db.getUserByEmail(email)
}

func DeleteUserByID(id int64) error {
	return global_db.deleteUserByID(id)
}
