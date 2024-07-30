package db

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/mattn/go-sqlite3"
	"github.com/samuelemusiani/doit/cmd/doit"
	// "errors"
	// "github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func newSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) migrate() error {
	query := `
  CREATE TABLE IF NOT EXISTS notes(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL
  );
  `
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) createNote(note doit.Note) (*doit.Note, error) {
	res, err := r.db.Exec("INSERT INTO notes(title, description) values(?, ?)", note.Title, note.Description)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	note.ID = id
	return &note, nil
}

func (r *SQLiteRepository) allNotes() ([]doit.Note, error) {
	rows, err := r.db.Query("SELECT * FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []doit.Note
	for rows.Next() {
		var note doit.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Description); err != nil {
			return nil, err
		}

		all = append(all, note)
	}

	return all, nil
}

func (r *SQLiteRepository) getNoteByID(id int64) (*doit.Note, error) {
	row := r.db.QueryRow("SELECT * FROM notes WHERE id = ?", id)

	var note doit.Note
	if err := row.Scan(&note.ID, &note.Title, &note.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &note, nil
}

func (r *SQLiteRepository) deleteNoteByID(id int64) error {
	res, err := r.db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		slog.With("err", err, "id", id).Error("Deleting note from DB")
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.With("err", err).Error("Getting rowsAffected by DELETE in db")
		return err
	}

	if rowsAffected == 0 {
		slog.With("id", id).Debug("Deleting note from DB. 0 Rows affected")
		return ErrDeleteFailed
	}

	return err
}
