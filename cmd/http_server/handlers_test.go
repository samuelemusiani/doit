package http_server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func request(method string, endpoint string, funcToTest http.HandlerFunc) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(funcToTest)

	handler.ServeHTTP(rr, req)
	return rr, nil
}

func TestRootAPIHandler(t *testing.T) {
	rr, err := request("GET", "/api", rootAPIHandler)
	assert.NilError(t, err)
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "Root APIs endpoint for DOIT")
}

func TestNotesHandlerOPTIONS(t *testing.T) {
	rr, err := request("OPTIONS", "/api/notes", notesHandler)
	assert.NilError(t, err)
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Result().Header.Get("Allow"), "GET OPTIONS POST")
}
