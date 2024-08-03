package db

import (
	"math/rand"
	"slices"
	"testing"

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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range n {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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

func newNote() doit.Note {
	return doit.Note{
		Title:       randString(10),
		Description: randString(250),
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

func createAndInsertNote(userID int64) (*doit.Note, error) {
	note := newNote()
	note.UserID = userID
	return CreateNote(note)
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

	assert.DeepEqual(t, users, dbUsers)

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

func TestCreateNoteWithoutUser(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	note := newNote()

	_, err = CreateNote(note)
	assert.ErrorContains(t, err, "FOREIGN KEY constraint failed")

	err = cleanup()
	assert.NilError(t, err)
}

func TestCreateNote(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	u, err := createAndInsertUser()
	assert.NilError(t, err)

	n := newNote()
	n.UserID = u.ID
	nn, err := CreateNote(n)
	assert.NilError(t, err)

	n.ID = nn.ID
	assert.DeepEqual(t, &n, nn)

	err = cleanup()
	assert.NilError(t, err)
}

func TestAllNotes(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	numberOfUsers := 3
	users := make([]doit.User, numberOfUsers)
	for i := range numberOfUsers {
		user, err := createAndInsertUser()
		assert.NilError(t, err)

		users[i] = *user
	}

	numberOfNotes := rand.Intn(10) + 10
	notes := make([]doit.Note, numberOfNotes)
	for i := range numberOfNotes {
		note, err := createAndInsertNote(int64(rand.Intn(numberOfUsers) + 1))
		assert.NilError(t, err)

		notes[i] = *note
	}

	notesFromUserID := func(notes []doit.Note, userID int64) []doit.Note {
		var newNotes []doit.Note
		for i := range notes {
			if notes[i].UserID == userID {
				newNotes = append(newNotes, notes[i])
			}
		}
		return newNotes
	}

	sortFunc := func(a, b doit.Note) int { return int(a.ID - b.ID) }

	for userID := int64(1); userID <= int64(numberOfUsers); userID++ {
		dbNotes, err := AllNotes(userID)
		assert.NilError(t, err)

		notesFiltered := notesFromUserID(notes, userID)
		slices.SortFunc(notesFiltered, sortFunc)
		assert.DeepEqual(t, notesFiltered, dbNotes)
	}

	err = cleanup()
	assert.NilError(t, err)
}

func TestGetNoteById(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	note, err := createAndInsertNote(user.ID)
	assert.NilError(t, err)

	getNote, err := GetNoteByID(note.ID)
	assert.NilError(t, err)

	assert.DeepEqual(t, note, getNote)

	err = cleanup()
	assert.NilError(t, err)
}

func TestDeleteNoteByID(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	note, err := createAndInsertNote(user.ID)
	assert.NilError(t, err)

	err = DeleteNoteByID(note.ID, user.ID)
	assert.NilError(t, err)

	notes, err := AllNotes(user.ID)
	assert.NilError(t, err)

	for i := range notes {
		if notes[i].ID == note.ID {
			t.Fatalf("Note present in db after deletion")
		}
	}

	err = cleanup()
	assert.NilError(t, err)
}

func TestDeleteNoteByIDWrongUserID(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	user, err := createAndInsertUser()
	assert.NilError(t, err)

	note, err := createAndInsertNote(user.ID)
	assert.NilError(t, err)

	err = DeleteNoteByID(note.ID, 123982)
	assert.ErrorIs(t, err, ErrDeleteFailed)

	err = cleanup()
	assert.NilError(t, err)
}

func TestDeleteNotesByUserID(t *testing.T) {
	err := setup()
	assert.NilError(t, err)

	numberOfUsers := 3
	users := make([]doit.User, numberOfUsers)
	for i := range numberOfUsers {
		user, err := createAndInsertUser()
		assert.NilError(t, err)

		users[i] = *user
	}

	numberOfNotes := rand.Intn(10) + 10
	notes := make([]doit.Note, numberOfNotes)
	for i := range numberOfNotes {
		note, err := createAndInsertNote(int64(rand.Intn(numberOfUsers) + 1))
		assert.NilError(t, err)

		notes[i] = *note
	}

	userID := int64(rand.Intn(numberOfUsers) + 1)
	err = DeleteNotesByUserID(userID)
	assert.NilError(t, err)

	dbNotes, err := AllNotes(userID)
	assert.NilError(t, err)
	assert.Check(t, len(dbNotes) == 0)

	err = cleanup()
	assert.NilError(t, err)
}
