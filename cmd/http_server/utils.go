package http_server

import (
	"errors"
	"net/http"

	"github.com/samuelemusiani/doit/cmd/db"
)

var (
	ErrInteral      = errors.New("Internal error")
	ErrUnauthorized = errors.New("Not authorized")
)

func sliceContains(s []string, e string) bool {
	for i := range s {
		if s[i] == e {
			return true
		}
	}
	return false
}

func isAdminFromRequest(r *http.Request) (bool, error) {
	c, err := r.Cookie(SESSION_COOCKIE_NAME)
	if err != nil {
		return false, errors.Join(ErrInteral, err)
	}

	s, ok := getSession(c.Value)
	if !ok || s.isExpired() {
		return false, ErrUnauthorized
	}

	user, err := db.GetUserByID(s.userID)
	if err != nil {
		return false, errors.Join(ErrInteral, err)
	}

	return user.Admin, nil
}
