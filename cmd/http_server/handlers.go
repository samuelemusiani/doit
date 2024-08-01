package http_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/samuelemusiani/doit/cmd/db"
	"github.com/samuelemusiani/doit/cmd/doit"
	"golang.org/x/crypto/bcrypt"
)

const SESSION_COOCKIE_NAME = "ST"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Not implemented yet"))
	return
}

func rootAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello there :)"))
	return
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET OPTIONS POST")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
		return
	}

	switch r.Method {
	case http.MethodGet:
		notesHandlerGET(w, r)
	case http.MethodPost:
		notesHandlerPOST(w, r)
	default:
	}

	return
}

func notesHandlerGET(w http.ResponseWriter, r *http.Request) {
	notes, err := db.AllNotes()
	if err != nil {
		slog.With("err", err).Error("While getting notes from DB")
		http.Error(w, "Could not get notes", http.StatusInternalServerError)
		return
	}

	b_notes, err := json.Marshal(notes)
	if err != nil {
		slog.With("err", err).Error("While parsing notes for json")
		http.Error(w, "Could not get notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b_notes)
}

func notesHandlerPOST(w http.ResponseWriter, r *http.Request) {
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
	noteCreated, err := db.CreateNote(note)
	if err != nil {
		slog.With("note", note, "err", err).Error("Adding note to db")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	jnote, err := json.Marshal(noteCreated)
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
		w.Header().Set("Allow", "GET OPTIONS DELETE")
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

	switch r.Method {
	case http.MethodGet:
		singleNoteHandlerGET(w, r, id)
	case http.MethodDelete:
		singleNoteHandlerDELETE(w, r, id)
	default:
		// Should never be here
		fmt.Println("Why here?")
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
	}
	return
}

func singleNoteHandlerGET(w http.ResponseWriter, r *http.Request, id int64) {
	note, err := db.GetNoteByID(id)
	if err != nil {
		slog.With("err", err, "id", id).Error("Getting notes")
		// We could send 404 here
		http.Error(w, "Could not get note", http.StatusInternalServerError)
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

func singleNoteHandlerDELETE(w http.ResponseWriter, r *http.Request, id int64) {
	slog.With("id", id).Debug("Deleting note")
	err := db.DeleteNoteByID(id)
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
	if r.Method == "OPTIONS" {
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

	user, err := db.GetUserByUsername(s.username)
	w.Write([]byte(fmt.Sprintf("Logged in as user %s with id %d", user.Username, user.ID)))
	return
}

func loginHandlerPOST(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(SESSION_COOCKIE_NAME)
	if !errors.Is(err, http.ErrNoCookie) {
		if err != nil {
			slog.With("err", err).Error("While getting cookies")
		} else {
			slog.Debug("hey")
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

	expire := time.Now().Add(2 * 24 * time.Hour)
	sToken := newSession(session{username: user.Username, expire: expire})
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
