package db

import (
	"database/sql"

	"github.com/samuelemusiani/doit/cmd/doit"
	// "errors"
	// "github.com/mattn/go-sqlite3"
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
    title TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL
  );
  `
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) createNote(note doit.Note) (*doit.Note, error) {
	res, err := r.db.Exec("INSERT INTO notes(title, description) values(?, ?)", note.Title, note.Description)

	if err != nil {
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
