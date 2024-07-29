package http_server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/samuelemusiani/doit/cmd/config"
)

var router *mux.Router = nil

func rootHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World from http :)"))
	return
}

func Init() {
	slog.Debug("Init http server")

	router = mux.NewRouter()
	router.HandleFunc("/", rootHandler)
	router.Use(logginMiddleware)
}

func ListenAndServe() {
	config := config.GetConfig()
	addr := config.ListeningAddress + ":" + strconv.Itoa(int(config.ListeningPort))

	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		slog.With("addr", addr).Info("Listening and serving")
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.With("err", err).Error("Listening and serving")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	wait := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	slog.Info("Shutting http server")
}
