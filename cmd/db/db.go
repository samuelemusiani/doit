package db

import (
	"database/sql"
	"errors"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/samuelemusiani/doit/cmd/config"
	"github.com/samuelemusiani/doit/cmd/doit"
)

var global_db *SQLiteRepository = nil

func Init() error {
	slog.Debug("Init db connection")

	config := config.GetConfig()

	rawDB, err := sql.Open("sqlite3", config.Databse.Path+"?_foreign_keys=on")
	if err != nil {
		return err
	}

	global_db = newSQLiteRepository(rawDB)
	err = global_db.migrate()
	if err != nil {
		return err
	}

	err = fillDB()
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

func AllNotes(userID int64) ([]doit.Note, error) {
	return global_db.allNotes(userID)
}

func GetNoteByID(id int64) (*doit.Note, error) {
	return global_db.getNoteByID(id)
}

// Delte note with id noteID only if userID match
func DeleteNoteByID(noteID int64, userID int64) error {
	return global_db.deleteNoteByID(noteID, userID)
}

func DeleteNotesByUserID(userID int64) error {
	return global_db.deleteNotesByUserID(userID)
}

func UpdateNote(noteID int64, note doit.Note, userID int64) (*doit.Note, error) {
	return global_db.updateNote(noteID, note, userID)
}

func CreateUser(user doit.User) (*doit.User, error) {
	return global_db.createUser(user)
}

func AllUsers() ([]doit.User, error) {
	return global_db.allUsers()
}

func GetUserByID(id int64) (*doit.User, error) {
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

func UpdateUser(id int64, user doit.User) (*doit.User, error) {
	return global_db.updateUser(id, user)
}

func fillDB() error {
	err := global_db.insertNoteStates(doit.States)
	if err != nil {
		return errors.Join(err, errors.New("Inserting states into db"))
	}

	err = global_db.insertNotePriorities(doit.Priorities)
	if err != nil {
		return errors.Join(err, errors.New("Inserting states into db"))
	}

	err = global_db.insertNoteColors(doit.Colors)
	if err != nil {
		return errors.Join(err, errors.New("Inserting states into db"))
	}

	return nil
}
