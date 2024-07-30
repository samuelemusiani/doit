package http_server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/samuelemusiani/doit/cmd/db"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the root of DOIT"))
	return
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET OPTIONS")
		w.Write([]byte{})
		return
	}

	if !equalMethods(r.Method, http.MethodGet, w) {
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

	w.Write(b_notes)
	return
}
