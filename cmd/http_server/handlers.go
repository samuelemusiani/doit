package http_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/samuelemusiani/doit/cmd/db"
	"github.com/samuelemusiani/doit/cmd/doit"
	"golang.org/x/crypto/bcrypt"
)

const SESSION_COOCKIE_NAME = "ST"

func staticHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path[1:]
	if p == "" || p == "static" || p == "static/" {
		p = "index.html"
	}

	f, err := fs.ReadFile(ui_fs, p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// If the file does not exists it could be a route that the SPA router
			// would catch. We serve the index.html instead

			f, err = fs.ReadFile(ui_fs, "index.html")
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					http.Error(w, "", http.StatusNotFound)
				} else {
					slog.With("err", err).Error("Reading index.html")
					http.Error(w, "", http.StatusInternalServerError)
				}
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(f)
			return
		}
		slog.With("path", p, "err", err).Error("Reading file")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	switch path.Ext(p) {
	case ".js":
		w.Header().Set("Content-Type", "text/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	}
	w.Write(f)
}

func rootAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Root APIs endpoint for DOIT"))
	return
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET OPTIONS POST")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
		return
	}

	c, err := r.Cookie(SESSION_COOCKIE_NAME)
	if err != nil {
		slog.With("err", err).Error("At this stage cookie should be present")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	s, b := getSession(c.Value)
	if !b || s.isExpired() {
		slog.With("err", err).Error("At this stage cookie should valid")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		notesHandlerGET(w, r, s.userID)
	case http.MethodPost:
		notesHandlerPOST(w, r, s.userID)
	default:
	}

	return
}

func notesHandlerGET(w http.ResponseWriter, r *http.Request, userID int64) {
	notes, err := db.AllNotes(userID)
	if err != nil {
		slog.With("err", err).Error("While getting notes from DB")
		http.Error(w, "Could not get notes", http.StatusInternalServerError)
		return
	}

	var response []byte
	if len(notes) == 0 {
		response = []byte("[]")
	} else {
		response, err = json.Marshal(notes)
		if err != nil {
			slog.With("err", err).Error("While parsing notes for json")
			http.Error(w, "Could not get notes", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func notesHandlerPOST(w http.ResponseWriter, r *http.Request, userID int64) {
	decoder := json.NewDecoder(r.Body)
	var note doit.Note
	err := decoder.Decode(&note)
	if err != nil {
		slog.With("err", err).Error("Decoding request body")
		// Should we return another type of error?
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	// Basic check, we could improve it in the future
	if note.Title == "" {
		http.Error(w, "Title is empty or not present", http.StatusBadRequest)
		return
	}

	slog.With("note", note).Debug("Adding note to db")
	note.UserID = userID
	noteCreated, err := db.CreateNote(note)
	if err != nil {
		slog.With("note", note, "err", err).Error("Adding note to db")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	jnote, err := json.Marshal(doit.NoteToResponse(noteCreated))
	if err != nil {
		slog.With("note", note, "err", err).Error("Could not parse note to json")
		http.Error(w, "Note was added but we could not send the note back", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jnote)
	return
}

func singleNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET OPTIONS PUT DELETE")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
		return
	}

	id_string, ok := mux.Vars(r)["id"]
	if !ok {
		slog.With("vars", mux.Vars(r)).Error("Could not get id from router vars in singleNoteHandler")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(id_string, 10, 64)
	if err != nil {
		slog.With("err", err).Error("Parsing int")
		http.Error(w, "Id is not valid", http.StatusBadRequest)
		return
	}

	c, err := r.Cookie(SESSION_COOCKIE_NAME)
	if err != nil {
		slog.With("err", err).Error("At this stage cookie should be present")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	s, b := getSession(c.Value)
	if !b || s.isExpired() {
		slog.With("err", err).Error("At this stage cookie should valid")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		singleNoteHandlerGET(w, r, id, s.userID)
	case http.MethodDelete:
		singleNoteHandlerDELETE(w, r, id, s.userID)
	case http.MethodPut:
		singleNoteHandlerPUT(w, r, id, s.userID)
	default:
		slog.With("method", r.Method).Error("Method not valid. How did we get here?")
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
	}
	return
}

func singleNoteHandlerGET(w http.ResponseWriter, r *http.Request, noteID int64, userId int64) {
	note, err := db.GetNoteByID(noteID)
	if err != nil {
		slog.With("err", err, "id", noteID).Error("Getting notes")
		if errors.Is(err, db.ErrNotExists) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte("Could not get note"))
		return
	}

	if note.UserID != userId {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Could not get note"))
		return
	}

	jnote, err := json.Marshal(note)
	if err != nil {
		slog.With("err", err).Error("While parsing note for json")
		http.Error(w, "Could not get note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jnote)
}

func singleNoteHandlerPUT(w http.ResponseWriter, r *http.Request, noteID int64, userID int64) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.With("err", err).Error("Reading body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var note doit.Note
	err = json.Unmarshal(body, &note)
	if err != nil {
		http.Error(w, "Could not unmarshal body", http.StatusBadRequest)
		return
	}

	note.UserID = userID

	newNote, err := db.UpdateNote(noteID, note, userID)
	if err != nil {
		slog.With("err", err).Error("Updating note")
		http.Error(w, "Could not update note", http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(*newNote)
	if err != nil {
		slog.With("err", err).Error("Marshaling note update")
		w.Write([]byte("Note updated, but can't be returned"))
		return
	}

	w.Write(b)
}

func singleNoteHandlerDELETE(w http.ResponseWriter, r *http.Request, noteID int64, userID int64) {
	err := db.DeleteNoteByID(noteID, userID)
	if err != nil {
		if errors.Is(err, db.ErrDeleteFailed) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte("Error deleting note"))
		return
	}
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET OPTIONS POST DELETE")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
		return
	}

	switch r.Method {
	case http.MethodGet:
		loginHandlerGET(w, r)
	case http.MethodPost:
		loginHandlerPOST(w, r)
	case http.MethodDelete:
		loginHandlerDELETE(w, r)
	}
	return
}

func loginHandlerGET(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(SESSION_COOCKIE_NAME)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}
		slog.With("err", err).Error("While getting cookies")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	s, ok := getSession(c.Value)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := db.GetUserByID(s.userID)
	if err != nil {
		http.Error(w, "Could not get user", http.StatusInternalServerError)
		return
	}

	userResp := doit.UserToResponse(user)
	b, err := json.Marshal(*userResp)
	if err != nil {
		http.Error(w, "Could not marshal user response", http.StatusInternalServerError)
		return
	}
	w.Write(b)
	return
}

func loginHandlerPOST(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(SESSION_COOCKIE_NAME)
	if !errors.Is(err, http.ErrNoCookie) {
		if err != nil {
			slog.With("err", err).Error("While getting cookies")
		} else {
			s, p := getSession(c.Value)
			if p && !s.isExpired() {
				deleteSession(c.Value)
			}
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.With("err", err).Error("Reading body")
		http.Error(w, "Reading body", http.StatusInternalServerError)
		return
	}

	type UserPasswd struct {
		Username string
		Password string
	}

	var u UserPasswd
	err = json.Unmarshal(body, &u)
	if err != nil {
		slog.With("err", err, "body", body).Error("Unmarshaling body")
		http.Error(w, "Reading body", http.StatusBadRequest)
		return
	}

	if len(u.Username) == 0 || len(u.Password) == 0 {
		http.Error(w, "Username or password are empty", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByUsername(u.Username)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			http.Error(w, "User does not exists or password is not correct", http.StatusNotFound)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		slog.With("err", err, "username", u.Username).Error("During user lookup on db")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			http.Error(w, "User does not exists or password is not correct", http.StatusNotFound)
			return
		}
		slog.With("err", err, "user", user).Error("Comparing hash with password hashed")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if !user.Active {
		http.Error(w, "Username and password are correct, but user in not active", http.StatusForbidden)
		return
	}

	expire := time.Now().Add(2 * 24 * time.Hour)
	sToken := newSession(session{userID: user.ID, expire: expire})
	http.SetCookie(w, &http.Cookie{
		Name:  SESSION_COOCKIE_NAME,
		Value: sToken,
		// Domain ??
		Path:     "/",
		Expires:  expire,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	slog.With("user", u.Username).Info("Logged in")
	w.Write([]byte(fmt.Sprintf("Logged in as user %s with id %d", user.Username, user.ID)))
	return
}

func loginHandlerDELETE(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(SESSION_COOCKIE_NAME)
	if err != nil {
		w.WriteHeader(http.StatusResetContent)
		return
	}
	deleteSession(c.Value)
	w.WriteHeader(http.StatusResetContent)
	return
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET POST OPTIONS")
		w.WriteHeader(http.StatusOK)
		return
	}

	isAdmin, err := isAdminFromRequest(r)
	if errors.Is(err, ErrInteral) {
		slog.With("err", err).Error("Checking if user is admin")
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else if errors.Is(err, ErrUnauthorized) {
		// This should never happen
		slog.Error("Passing middleware of authentication, but not authenticated")
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	if !isAdmin {
		http.Error(w, "Not authenticated", http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodGet:
		usersHandlerGET(w, r)
	case http.MethodPost:
		usersHandlerPOST(w, r)
	}

	return
}

func usersHandlerGET(w http.ResponseWriter, r *http.Request) {
	users, err := db.AllUsers()
	if err != nil {
		slog.With("err", err).Error("Gettin users from DB")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	usersResponse := make([]doit.UserResponse, len(users))
	for i := range users {
		usersResponse[i] = *doit.UserToResponse(&users[i])
	}

	var res []byte
	if len(usersResponse) == 0 {
		res = []byte("[]")
	} else {
		res, err = json.Marshal(usersResponse)
		if err != nil {
			slog.With("err", err).Error("Marshaling users for response")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	w.Write(res)
	return
}

func usersHandlerPOST(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.With("err", err).Error("Could not read body of a request")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var u_user doit.UserUnmarshaling
	err = json.Unmarshal(body, &u_user)
	if err != nil {
		slog.With("err", err).Error("Unmarshaling body")
		http.Error(w, "Could not unmarshal body", http.StatusBadRequest)
		return
	}

	if u_user.Username == nil || u_user.Password == nil || u_user.Email == nil {
		http.Error(w, "Username, password or email are not present", http.StatusBadRequest)
		return
	}

	user := doit.UserUnmarshalingToUser(&u_user)

	pw := user.Password
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			http.Error(w, "Password too long (> 72 bytes)", http.StatusBadRequest)
		} else {
			slog.With("err", err).Error("Generating hash from password")
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	user.Password = string(h)

	new_user, err := db.CreateUser(*user)
	if err != nil {
		if errors.Is(err, db.ErrDuplicate) {
			http.Error(w, "User already present", http.StatusBadRequest)
		} else {
			slog.With("err", err).Error("Inserting new user into db")
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	user_res := doit.UserToResponse(new_user)
	res, err := json.Marshal(user_res)
	if err != nil {
		slog.With("err", err).Error("Marshaling update user")
		w.Write([]byte("User created, but can't be returned"))
		return
	}

	w.Write(res)
}

func singleUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET OPTIONS PUT DELETE")
		w.WriteHeader(http.StatusOK)
		return
	}

	id_string, ok := mux.Vars(r)["id"]
	if !ok {
		slog.With("vars", mux.Vars(r)).Error("Could not get id from router vars in singleNoteHandler")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(id_string, 10, 64)
	if err != nil {
		slog.With("err", err).Error("Parsing int")
		http.Error(w, "Id is not valid", http.StatusBadRequest)
		return
	}

	c, err := r.Cookie(SESSION_COOCKIE_NAME)
	if err != nil {
		slog.With("err", err).Error("At this stage cookie should be present")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	s, b := getSession(c.Value)
	if !b || s.isExpired() {
		slog.With("err", err).Error("At this stage cookie should valid")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	user, err := db.GetUserByID(s.userID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if !user.Admin && user.ID != id {
		http.Error(w, "Your not an admin or this is not your account", http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodGet:
		singleUserHandlerGET(w, r, id, user)
	case http.MethodPut:
		singleUserHandlerPUT(w, r, id, user)
	case http.MethodDelete:
		singleUserHandlerDELETE(w, r, id, user)
	default:
		slog.With("method", r.Method).Error("Method not valid. How did we get here?")
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
	}
	return
}

func singleUserHandlerGET(w http.ResponseWriter, r *http.Request, userID int64, author *doit.User) {
	author, err := db.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			http.Error(w, "User does not exists", http.StatusNotFound)
			return
		}
		slog.With("err", err, "userId", userID).Error("Getting user in db")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	userResponse := doit.UserToResponse(author)
	res, err := json.Marshal(userResponse)
	if err != nil {
		slog.With("err", err, "user", userResponse).Error("Marshaling user to JSON")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func singleUserHandlerPUT(w http.ResponseWriter, r *http.Request, userID int64, author *doit.User) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.With("err", err).Error("Reading body")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	originalUser, err := db.GetUserByID(userID)
	if err != nil {
		slog.With("err", err, "userId", userID).Error("Getting user from db")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var updateRequested doit.UserUnmarshaling
	err = json.Unmarshal(body, &updateRequested)
	if err != nil {
		slog.With("err", err).Error("Unmarshaling body")
		http.Error(w, "Could not unmarshal body", http.StatusBadRequest)
		return
	}

	if updateRequested.Username != nil &&
		*updateRequested.Username != originalUser.Username {
		http.Error(w, "Username is not updatable", http.StatusBadRequest)
		return
	}

	if updateRequested.Email != nil {
		originalUser.Email = *updateRequested.Email
	}

	if updateRequested.Name != nil {
		originalUser.Name = *updateRequested.Name
	}

	if updateRequested.Surname != nil {
		originalUser.Surname = *updateRequested.Surname
	}

	if updateRequested.Admin != nil {
		if author.Admin {
			slog.With("author", *author, "updateUser", updateRequested).Info("Admin modification")
			originalUser.Admin = *updateRequested.Admin
		} else {
			http.Error(w, "Not an admin, cannot become one", http.StatusForbidden)
			return
		}
	}

	if updateRequested.Active != nil {
		originalUser.Active = *updateRequested.Active
	}

	if updateRequested.Password != nil {
		h, err := bcrypt.GenerateFromPassword([]byte(*updateRequested.Password), bcrypt.DefaultCost)
		if err != nil {
			if errors.Is(err, bcrypt.ErrPasswordTooLong) {
				http.Error(w, "Password too long (> 72 bytes)", http.StatusBadRequest)
			} else {
				slog.With("err", err).Error("Generating hash from password")
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		originalUser.Password = string(h)
	}

	updatedUser, err := db.UpdateUser(userID, *originalUser)
	if err != nil {
		slog.With("err", err).Error("Updating user")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	updateResponse := doit.UserToResponse(updatedUser)

	res, err := json.Marshal(updateResponse)
	if err != nil {
		slog.With("err", err).Error("Marshaling update user")
		w.Write([]byte("User updated, but can't be returned"))
		return
	}

	w.Write(res)
}

func singleUserHandlerDELETE(w http.ResponseWriter, r *http.Request, userID int64, author *doit.User) {
	err := db.DeleteNotesByUserID(userID)
	if err != nil && !errors.Is(err, db.ErrDeleteFailed) {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = db.DeleteUserByID(userID)
	if err != nil {
		if errors.Is(err, db.ErrDeleteFailed) {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		slog.With("err", err, "userID", userID).Error("Deleting user from DB")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User deleted successfuly"))
}

func noteStatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET OPTIONS")
		w.WriteHeader(http.StatusOK)
		return
	}

	b, err := json.Marshal(doit.States)
	if err != nil {
		slog.With("err", err).Error("Marshaling notes states")
		http.Error(w, "Could not get states", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func notePrioritiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET OPTIONS")
		w.WriteHeader(http.StatusOK)
		return
	}

	b, err := json.Marshal(doit.Priorities)
	if err != nil {
		slog.With("err", err).Error("Marshaling notes priorities")
		http.Error(w, "Could not get states", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func noteColorsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET OPTIONS")
		w.WriteHeader(http.StatusOK)
		return
	}

	b, err := json.Marshal(doit.Colors)
	if err != nil {
		slog.With("err", err).Error("Marshaling notes colors")
		http.Error(w, "Could not get states", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
