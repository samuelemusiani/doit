package http_server

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/samuelemusiani/doit/cmd/config"
)

var router *mux.Router = nil

var NO_AUTH_PATHS = [...]string{
	"/",
	"/api",
	"/api/login",
	"/api/options/states",
	"/api/options/priorities",
	"/api/options/colors",
}

func Init(fs fs.FS) {
	slog.Debug("Init http server")

	router = mux.NewRouter()
	router.HandleFunc("/api", rootAPIHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/notes", notesHandler).Methods("GET", "OPTIONS", "POST")
	router.HandleFunc("/api/notes/{id}", singleNoteHandler).Methods("GET", "OPTIONS", "PUT", "DELETE")
	router.HandleFunc("/api/login", loginHandler).Methods("GET", "OPTIONS", "POST", "DELETE")
	router.HandleFunc("/api/users", usersHandler).Methods("GET", "POST", "OPTIONS")
	router.HandleFunc("/api/users/{id}", singleUserHandler).Methods("GET", "OPTIONS", "PUT", "DELETE")

	router.HandleFunc("/api/options/states", noteStatesHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/options/priorities", notePrioritiesHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/options/colors", noteColorsHandler).Methods("GET", "OPTIONS")

	router.PathPrefix("/").Handler(http.FileServerFS(fs))

	router.Use(logginMiddleware)
	router.Use(authMiddleware)
}

func ListenAndServe() error {
	config := config.GetConfig()
	addr := config.Server.Listen

	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	errc := make(chan error, 1)

	go func() {
		slog.With("addr", addr).Info("Listening and serving")
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- err
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	select {
	case <-c:
		slog.Info("Received SIGINT")
	case err := <-errc:
		return err
	}

	// Create a deadline to wait for.
	wait := 5 * time.Second
	slog.With("wait", wait).Debug("Waiting for http server to shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	slog.Info("Shutting http server")

	return nil
}
