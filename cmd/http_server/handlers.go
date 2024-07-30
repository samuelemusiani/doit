package http_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/samuelemusiani/doit/cmd/db"
)

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
		w.Header().Set("Allow", "GET OPTIONS")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
		return
	}

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
