package http_server

import (
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

type session struct {
	userID int64
	expire time.Time
}

var activeSessions sync.Map

func (s session) isExpired() bool {
	return s.expire.Before(time.Now())
}

func newSession(s session) string {
	t := uuid.NewString()

	activeSessions.Store(t, s)
	slog.With("token", t, "session", s).Debug("New session")

	activeSessions.Range(func(key, value any) bool {
		slog.With("key", key, "value", value).Debug("")
		return true
	})

	return t
}

// Return the session and true if is present, false otherwise
func getSession(token string) (session, bool) {
	s, present := activeSessions.Load(token)
	if !present {
		s = session{}
	}
	return s.(session), present
}

func deleteSession(token string) {
	activeSessions.Delete(token)
}
