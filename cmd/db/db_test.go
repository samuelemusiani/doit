package db

import (
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/samuelemusiani/doit/cmd/config"
	"github.com/samuelemusiani/doit/cmd/doit"
	"golang.org/x/crypto/bcrypt"
	"gotest.tools/v3/assert"
)

// Utils functions

func setup() error {
	conf := config.GetConfig()
	conf.Databse.Path = ":memory:"
	return Init()
}

func cleanup() error {
	return Close()
}

func randBool() bool {
	if rand.Intn(2) == 0 {
		return false
	} else {
		return true
	}
}

func newPassword() (string, error) {
	l := rand.Intn(24-8) + 8
	passwd := randString(l)

	h, e := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(h), e
}

func newTodo() doit.Todo {
	return doit.Todo{
		Title:       randString(10),
		Description: randString(250),
		StateID:     1,
		PriorityID:  1,
		ColorID:     1,
		Expiration: doit.Expiration{
			DoesExpire: true,
			Date:       time.Now().Add(time.Hour).Round(time.Second),
		},
	}
}

func newUser() (doit.User, error) {
	passwd, err := newPassword()
	return doit.User{
		Username: randString(8),
		Email:    randString(10) + "@" + randString(4) + "." + randString(3),
		Name:     randString(10),
		Surname:  randString(12),
		Admin:    randBool(),
		Active:   randBool(),
		Password: passwd,
	}, err
}

func createAndInsertUser() (*doit.User, error) {
	user, err := newUser()
	if err != nil {
		return &doit.User{}, err
	}

	return CreateUser(user)
}

func createAndInsertTodo(userID int64) (*doit.Todo, error) {
	todo := newTodo()
	todo.UserID = userID
	return CreateTodo(todo)
}

// Acctual testing

func TestInit(t *testing.T) {
	err := setup()
	assert.NilError(t, err)
	err = cleanup()
	assert.NilError(t, err)
}

func TestCreateUser(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := newUser()
	assert.NilError(t, err)

	newUser, err := CreateUser(user)
	user.ID = newUser.ID
	assert.NilError(t, err)
	assert.DeepEqual(t, &user, newUser)

	err = cleanup()
	assert.NilError(t, err)
}

func TestCreateUserDuplicate(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	_, err = CreateUser(*user)
	assert.ErrorIs(t, err, ErrDuplicate)

	err = cleanup()
	assert.NilError(t, err)
}

func TestAllUsers(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	n := rand.Intn(10) + 10
	users := make([]doit.User, n)

	for i := range n {
		user, err := createAndInsertUser()
		assert.NilError(t, err)

		users[i] = *user
	}

	dbUsers, err := AllUsers()
	assert.NilError(t, err)

	sortFunc := func(a, b doit.User) int { return int(a.ID - b.ID) }
	slices.SortFunc(users, sortFunc)
	slices.SortFunc(dbUsers, sortFunc)

	// Skip the first user as not generated here, it's the default admin
	assert.DeepEqual(t, users, dbUsers[1:])

	err = cleanup()
	assert.NilError(t, err)
}

func TestGetUserByID(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	getUser, err := GetUserByID(user.ID)
	assert.NilError(t, err)
	assert.DeepEqual(t, user, getUser)

	err = cleanup()
	assert.NilError(t, err)
}

func TestGetUserByUsername(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	getUser, err := GetUserByUsername(user.Username)
	assert.NilError(t, err)
	assert.DeepEqual(t, user, getUser)

	err = cleanup()
	assert.NilError(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	getUser, err := GetUserByEmail(user.Email)
	assert.NilError(t, err)
	assert.DeepEqual(t, user, getUser)

	err = cleanup()
	assert.NilError(t, err)
}

func TestDeleteUserByID(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	err = DeleteUserByID(user.ID)
	assert.NilError(t, err)

	users, err := AllUsers()
	assert.NilError(t, err)

	for i := range users {
		if users[i].ID == user.ID {
			t.Fatalf("User present in db after deletion")
		}
	}

	err = cleanup()
	assert.NilError(t, err)
}

func TestUpdateUser(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	nUser, err := createAndInsertUser()
	assert.NilError(t, err)

	user, err := newUser()
	assert.NilError(t, err)
	user.ID = nUser.ID

	uUser, err := UpdateUser(nUser.ID, user)
	assert.NilError(t, err)
	assert.DeepEqual(t, &user, uUser)

	err = cleanup()
	assert.NilError(t, err)
}

func TestCreateTodoWithoutUser(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	todo := newTodo()

	_, err = CreateTodo(todo)
	assert.ErrorContains(t, err, "FOREIGN KEY constraint failed")

	err = cleanup()
	assert.NilError(t, err)
}

func TestCreateTodo(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	u, err := createAndInsertUser()
	assert.NilError(t, err)

	n := newTodo()
	n.UserID = u.ID
	nn, err := CreateTodo(n)
	assert.NilError(t, err)

	n.ID = nn.ID
	assert.DeepEqual(t, &n, nn)

	err = cleanup()
	assert.NilError(t, err)
}

func TestAllTodos(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	numberOfUsers := 3
	users := make([]doit.User, numberOfUsers)
	for i := range numberOfUsers {
		user, err := createAndInsertUser()
		assert.NilError(t, err)

		users[i] = *user
	}

	numberOfTodos := rand.Intn(10) + 10
	todos := make([]doit.Todo, numberOfTodos)
	for i := range numberOfTodos {
		todo, err := createAndInsertTodo(int64(rand.Intn(numberOfUsers) + 1))
		assert.NilError(t, err)

		todos[i] = *todo
	}

	todosFromUserID := func(todos []doit.Todo, userID int64) []doit.Todo {
		var newTodos []doit.Todo
		for i := range todos {
			if todos[i].UserID == userID {
				newTodos = append(newTodos, todos[i])
			}
		}
		return newTodos
	}

	sortFunc := func(a, b doit.Todo) int { return int(a.ID - b.ID) }

	for userID := int64(1); userID <= int64(numberOfUsers); userID++ {
		dbTodos, err := AllTodos(userID)
		assert.NilError(t, err)

		todosFiltered := todosFromUserID(todos, userID)
		slices.SortFunc(todosFiltered, sortFunc)
		assert.DeepEqual(t, todosFiltered, dbTodos)
	}

	err = cleanup()
	assert.NilError(t, err)
}

func TestGetTodoById(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	todo, err := createAndInsertTodo(user.ID)
	assert.NilError(t, err)

	getTodo, err := GetTodoByID(todo.ID)
	assert.NilError(t, err)

	assert.DeepEqual(t, todo, getTodo)

	err = cleanup()
	assert.NilError(t, err)
}

func TestUpdateTodoByID(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	todo, err := createAndInsertTodo(user.ID)
	assert.NilError(t, err)

	newTodo := newTodo()
	newTodo.ID = todo.ID
	newTodo.UserID = todo.UserID

	modTodo, err := UpdateTodo(todo.ID, newTodo, todo.UserID)
	assert.NilError(t, err)
	assert.Equal(t, *modTodo, newTodo)
}

func TestDeleteTodoByID(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	todo, err := createAndInsertTodo(user.ID)
	assert.NilError(t, err)

	err = DeleteTodoByID(todo.ID, user.ID)
	assert.NilError(t, err)

	todos, err := AllTodos(user.ID)
	assert.NilError(t, err)

	for i := range todos {
		if todos[i].ID == todo.ID {
			t.Fatalf("Todo present in db after deletion")
		}
	}

	err = cleanup()
	assert.NilError(t, err)
}

func TestDeleteTodoByIDWrongUserID(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	todo, err := createAndInsertTodo(user.ID)
	assert.NilError(t, err)

	err = DeleteTodoByID(todo.ID, 123982)
	assert.ErrorIs(t, err, ErrDeleteFailed)

	err = cleanup()
	assert.NilError(t, err)
}

func TestDeleteTodosByUserID(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	numberOfUsers := 3
	users := make([]doit.User, numberOfUsers)
	for i := range numberOfUsers {
		user, err := createAndInsertUser()
		assert.NilError(t, err)

		users[i] = *user
	}

	numberOfTodos := rand.Intn(10) + 10
	todos := make([]doit.Todo, numberOfTodos)
	for i := range numberOfTodos {
		todo, err := createAndInsertTodo(int64(rand.Intn(numberOfUsers) + 1))
		assert.NilError(t, err)

		todos[i] = *todo
	}

	userID := int64(rand.Intn(numberOfUsers) + 1)
	err = DeleteTodosByUserID(userID)
	assert.NilError(t, err)

	dbTodos, err := AllTodos(userID)
	assert.NilError(t, err)
	assert.Check(t, len(dbTodos) == 0)

	err = cleanup()
	assert.NilError(t, err)
}

func TestInsertTodoStates(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	err = cleanup()
	assert.NilError(t, err)
}
