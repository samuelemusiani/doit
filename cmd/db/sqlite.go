package db

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/samuelemusiani/doit/cmd/doit"
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
  CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TINYTEXT UNIQUE NOT NULL,
    email TINYTEXT UNIQUE NOT NULL,
    name TINYTEXT NOT NULL,
    surname TINYTEXT NOT NULL,
    admin BOOL NOT NULL,
    external BOOL NOT NULL,
    active BOOL NOT NULL,
    password TINYTEXT NOT NULL
  );
  CREATE TABLE IF NOT EXISTS notes(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    userID INTEGER,
    FOREIGN KEY(userID) REFERENCES users(id)
  );
  `
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) createNote(note doit.Note) (*doit.Note, error) {
	res, err := r.db.Exec("INSERT INTO notes(title, description, userID) values(?, ?, ?)", note.Title, note.Description, note.UserID)

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

func (r *SQLiteRepository) createUser(user doit.User) (*doit.User, error) {
	res, err := r.db.Exec("INSERT INTO users(username, email, name, surname, admin, external, active, password) values(?, ?, ?, ?, ?, ?, ?, ?)",
		user.Username, user.Email, user.Name, user.Surname,
		user.Admin, user.External, user.Active, user.Password)

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

	user.ID = id
	return &user, nil
}

func (r *SQLiteRepository) allNotes(userId int64) ([]doit.Note, error) {
	rows, err := r.db.Query("SELECT * FROM notes WHERE userID = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []doit.Note
	for rows.Next() {
		var note doit.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Description, &note.UserID); err != nil {
			return nil, err
		}

		all = append(all, note)
	}

	return all, nil
}

func (r *SQLiteRepository) allUsers() ([]doit.User, error) {
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []doit.User
	for rows.Next() {
		var user doit.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.Admin, &user.External, &user.Active, &user.Password)
		if err != nil {
			return nil, err
		}

		all = append(all, user)
	}

	return all, nil
}

func (r *SQLiteRepository) getNoteByID(id int64) (*doit.Note, error) {
	row := r.db.QueryRow("SELECT * FROM notes WHERE id = ?", id)

	var note doit.Note
	if err := row.Scan(&note.ID, &note.Title, &note.Description, &note.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &note, nil
}

func scanUser(row *sql.Row) (*doit.User, error) {
	var user doit.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.Admin, &user.External, &user.Active, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}

	return &user, nil
}

func (r *SQLiteRepository) getUserByID(id int64) (*doit.User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	return scanUser(row)
}

func (r *SQLiteRepository) getUserByUsername(username string) (*doit.User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE username = ?", username)
	return scanUser(row)
}

func (r *SQLiteRepository) getUserByEmail(email string) (*doit.User, error) {
	row := r.db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	return scanUser(row)
}

// Delete note with id noteID only if userID match
func (r *SQLiteRepository) deleteNoteByID(noteID int64, userID int64) error {
	res, err := r.db.Exec("DELETE FROM notes WHERE id = ? AND userID = ?", noteID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}

func (r *SQLiteRepository) deleteUserByID(id int64) error {
	res, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}

func (r *SQLiteRepository) updateNote(id int64, note doit.Note) (*doit.Note, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}

	res, err := r.db.Exec("UPDATE notes SET title = ?, description = ?, userID = ? WHERE id = ?",
		note.Title, note.Description, note.UserID, note.ID)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &note, nil
}

func (r *SQLiteRepository) updateUser(id int64, user doit.User) (*doit.User, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}

	res, err := r.db.Exec("UPDATE users SET username = ?, email = ?, name = ?, surname = ?, admin = ?, external = ?, active = ?, password = ? WHERE id = ?",
		user.Username, user.Email, user.Name, user.Surname, user.Admin, user.External, user.Active, user.Password, user.ID)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &user, nil
}
