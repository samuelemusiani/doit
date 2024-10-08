package db

import (
	"database/sql"
	"errors"
	"log/slog"
	"math/rand"

	_ "github.com/mattn/go-sqlite3"
	"github.com/samuelemusiani/doit/cmd/config"
	"github.com/samuelemusiani/doit/cmd/doit"
	"golang.org/x/crypto/bcrypt"
)

var global_db *SQLiteRepository = nil

func Init() error {
	slog.Debug("Init db connection")

	config := config.GetConfig()

	rawDB, err := sql.Open("sqlite3", config.Databse.Path+"?_foreign_keys=on")
	if err != nil {
		return errors.Join(err, errors.New("Can't open db"))
	}

	global_db = newSQLiteRepository(rawDB)
	err = global_db.migrate()
	if err != nil {
		return errors.Join(err, errors.New("Can't generate tables on DB"))
	}

	err = fillDB()
	if err != nil {
		return errors.Join(err, errors.New("Can't fill DB with defaults tables"))
	}

	err = generateDefaultAdmin(config.Users.First_User)
	if err != nil {
		return errors.Join(err, errors.New("Can't create default user"))
	}

	return nil
}

func Close() error {
	slog.Debug("Closing DB connection")
	err := global_db.db.Close()
	return err
}

func CreateTodo(note doit.Todo) (*doit.Todo, error) {
	return global_db.createTodo(note)
}

func AllTodos(userID int64) ([]doit.Todo, error) {
	return global_db.allTodos(userID)
}

func GetTodoByID(id int64) (*doit.Todo, error) {
	return global_db.getTodoByID(id)
}

// Delte note with id noteID only if userID match
func DeleteTodoByID(noteID int64, userID int64) error {
	return global_db.deleteTodoByID(noteID, userID)
}

func DeleteTodosByUserID(userID int64) error {
	return global_db.deleteTodosByUserID(userID)
}

func UpdateTodo(noteID int64, note doit.Todo, userID int64) (*doit.Todo, error) {
	return global_db.updateTodo(noteID, note, userID)
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
	err := global_db.insertTodoStates(doit.States)
	if err != nil {
		return errors.Join(err, errors.New("Inserting states into db"))
	}

	err = global_db.insertTodoPriorities(doit.Priorities)
	if err != nil {
		return errors.Join(err, errors.New("Inserting states into db"))
	}

	err = global_db.insertTodoColors(doit.Colors)
	if err != nil {
		return errors.Join(err, errors.New("Inserting states into db"))
	}

	return nil
}

func generateDefaultAdmin(user config.FirstUser) error {
	u, err := global_db.getInternal("first_user")
	// User found
	if err == nil {
		// We check if config has changed, just to notify :)
		if string(u) != user.Username {
			slog.With("current_user", string(u), "new_user", user).Warn("Default user is changed in config but already present in DB. Can't do anything")
		}
		return nil
	}

	// If error is not related to user existance
	if !errors.Is(err, ErrNotExists) {
		return err
	}

	// User not found, it's the first time
	plain_passwd := randString(24)
	passwd, err := bcrypt.GenerateFromPassword([]byte(plain_passwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	new_user := doit.User{
		Username: user.Username,
		Email:    user.Email,
		Password: string(passwd),
		Admin:    true,
		Active:   true,
	}

	_, err = global_db.createUser(new_user)
	if err != nil {
		return err
	}

	err = global_db.addInternal("first_user", []byte(user.Username))
	if err != nil {
		return err
	}

	slog.With("password", plain_passwd).Info("First user generated")

	return nil
}

// This is only here for generating random password for the first user and for tests
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range n {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
