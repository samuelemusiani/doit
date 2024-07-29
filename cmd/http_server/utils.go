package http_server

import "net/http"

// Return false if the methods are equal. Respond with an error
func equalMethods(m1 string, m2 string, w http.ResponseWriter) bool {
	if m1 != m2 {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
