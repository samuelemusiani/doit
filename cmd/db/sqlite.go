package db

import (
	"database/sql"
	"errors"
	"time"

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
    active BOOL NOT NULL,
    password TINYTEXT NOT NULL
  );
  CREATE TABLE IF NOT EXISTS todo_states(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    state TINYTEXT NOT NULL
  );
  CREATE TABLE IF NOT EXISTS todo_priority(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    priority INTEGER NOT NULL
  );
  CREATE TABLE IF NOT EXISTS todo_colors(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    color TINYTEXT NOT NULL
  );
  CREATE TABLE IF NOT EXISTS todos(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    stateID INTEGER,
    priorityID INTEGER,
    colorID INTEGER,
    does_expire BOOL,
    expiration_date INTEGER,
    userID INTEGER,
    FOREIGN KEY(stateID) REFERENCES todo_states(id),
    FOREIGN KEY(priorityID) REFERENCES todo_priority(id),
    FOREIGN KEY(colorID) REFERENCES todo_colors(id),
    FOREIGN KEY(userID) REFERENCES users(id)
  );
  CREATE TABLE IF NOT EXISTS internals(
    key TEXT NOT NULL UNIQUE PRIMARY KEY,
    data BLOB NOT NULL
  );
  `
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) createTodo(todo doit.Todo) (*doit.Todo, error) {
	res, err := r.db.Exec("INSERT INTO todos(title, description, stateID, priorityID, colorID, does_expire, expiration_date, userID) values(?, ?, ?, ?, ?, ?, ?, ?)", todo.Title, todo.Description, todo.StateID, todo.PriorityID, todo.ColorID, todo.Expiration.DoesExpire, todo.Expiration.Date.Unix(), todo.UserID)

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
	todo.ID = id
	return &todo, nil
}

func (r *SQLiteRepository) createUser(user doit.User) (*doit.User, error) {
	res, err := r.db.Exec("INSERT INTO users(username, email, name, surname, admin, active, password) values(?, ?, ?, ?, ?, ?, ?)",
		user.Username, user.Email, user.Name, user.Surname,
		user.Admin, user.Active, user.Password)

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

func (r *SQLiteRepository) allTodos(userId int64) ([]doit.Todo, error) {
	rows, err := r.db.Query("SELECT * FROM todos WHERE userID = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []doit.Todo
	for rows.Next() {
		var todo doit.Todo
		var t int64
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.StateID, &todo.PriorityID, &todo.ColorID, &todo.Expiration.DoesExpire, &t, &todo.UserID)
		if err != nil {
			return nil, err
		}

		if todo.Expiration.DoesExpire {
			todo.Expiration.Date = time.Unix(t, 0)
		}

		all = append(all, todo)
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
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.Admin, &user.Active, &user.Password)
		if err != nil {
			return nil, err
		}

		all = append(all, user)
	}

	return all, nil
}

func (r *SQLiteRepository) getTodoByID(id int64) (*doit.Todo, error) {
	row := r.db.QueryRow("SELECT * FROM todos WHERE id = ?", id)

	var todo doit.Todo
	var t int64
	err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.StateID, &todo.PriorityID, &todo.ColorID, &todo.Expiration.DoesExpire, &t, &todo.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	if todo.Expiration.DoesExpire {
		todo.Expiration.Date = time.Unix(t, 0)
	}
	return &todo, nil
}

func scanUser(row *sql.Row) (*doit.User, error) {
	var user doit.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.Admin, &user.Active, &user.Password)
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

// Delete todo with id todoID only if userID match
func (r *SQLiteRepository) deleteTodoByID(todoID int64, userID int64) error {
	res, err := r.db.Exec("DELETE FROM todos WHERE id = ? AND userID = ?", todoID, userID)
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

func (r *SQLiteRepository) deleteTodosByUserID(userID int64) error {
	res, err := r.db.Exec("DELETE FROM todos WHERE userID = ?", userID)
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

func (r *SQLiteRepository) updateTodo(id int64, todo doit.Todo, userID int64) (*doit.Todo, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}

	res, err := r.db.Exec("UPDATE todos SET title = ?, description = ?, stateID = ?, priorityID = ?, colorID = ?, does_expire = ?, expiration_date = ?, userID = ? WHERE id = ? AND userID = ?",
		todo.Title, todo.Description, todo.StateID, todo.PriorityID, todo.ColorID, todo.Expiration.DoesExpire, todo.Expiration.Date.Unix(), todo.UserID, todo.ID, userID)

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

	return &todo, nil
}

func (r *SQLiteRepository) updateUser(id int64, user doit.User) (*doit.User, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}

	res, err := r.db.Exec("UPDATE users SET username = ?, email = ?, name = ?, surname = ?, admin = ?, active = ?, password = ? WHERE id = ?",
		user.Username, user.Email, user.Name, user.Surname, user.Admin, user.Active, user.Password, user.ID)

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

func (r *SQLiteRepository) insertTodoStates(s []*doit.TodoState) error {
	for i := range s {
		row := r.db.QueryRow("SELECT * FROM todo_states WHERE state = ?", s[i].State)

		if err := row.Scan(&s[i].ID, &s[i].State); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				res, err := r.db.Exec("INSERT INTO todo_states(state) values(?)", s[i].State)
				if err != nil {
					return err
				}

				id, err := res.LastInsertId()
				if err != nil {
					return err
				}

				s[i].ID = id
			} else {
				return err
			}
		}
	}
	return nil
}

func (r *SQLiteRepository) insertTodoPriorities(s []*doit.TodoPriority) error {
	for i := range s {
		row := r.db.QueryRow("SELECT * FROM todo_priority WHERE priority = ?", s[i].Priority)

		if err := row.Scan(&s[i].ID, &s[i].Priority); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				res, err := r.db.Exec("INSERT INTO todo_priority(priority) values(?)", s[i].Priority)
				if err != nil {
					return err
				}

				id, err := res.LastInsertId()
				if err != nil {
					return err
				}

				s[i].ID = id
			} else {
				return err
			}
		}
	}
	return nil
}

func (r *SQLiteRepository) insertTodoColors(s []*doit.Color) error {
	for i := range s {
		row := r.db.QueryRow("SELECT * FROM todo_colors WHERE color = ?", s[i].Hex)

		if err := row.Scan(&s[i].ID, &s[i].Hex); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				res, err := r.db.Exec("INSERT INTO todo_colors(color) values(?)", s[i].Hex)
				if err != nil {
					return err
				}

				id, err := res.LastInsertId()
				if err != nil {
					return err
				}

				s[i].ID = id
			} else {
				return err
			}
		}
	}
	return nil
}

func (r *SQLiteRepository) getInternal(key string) ([]byte, error) {
	row := r.db.QueryRow("SELECT data FROM internals WHERE key = ?", key)

	var blob []byte
	if err := row.Scan(&blob); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		} else {
			return nil, err
		}
	}

	return blob, nil
}

func (r *SQLiteRepository) addInternal(key string, data []byte) error {
	_, err := r.db.Exec("INSERT INTO internals(key, data) values(?, ?)", key, data)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return ErrDuplicate
			}
		}
		return err
	}

	return nil
}
