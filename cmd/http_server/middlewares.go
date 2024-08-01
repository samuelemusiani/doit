package http_server

import (
	"errors"
	"log/slog"
	"net/http"
)

func logginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.With("method", r.Method, "URL", r.URL, "client", r.RemoteAddr, "agent", r.UserAgent()).Debug("")
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sliceContains(NO_AUTH_PATHS[:], r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		c, err := r.Cookie(SESSION_COOCKIE_NAME)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				slog.With("err", err).Error("Getting cookie")
			}
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		s, ok := getSession(c.Value)
		if !ok || s.isExpired() {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
