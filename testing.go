package nozzle

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// GetTestRequest creates a request for execution
func GetTestRequest(t *testing.T, verb string, path string, body string) (*http.Request, *httptest.ResponseRecorder) {

	reader := strings.NewReader(body)

	request, err := http.NewRequest(verb, path, reader)
	if err != nil {
		t.Fatal(err)
	}

	return request, httptest.NewRecorder()
}

// ExecuteTestRequest against the test server
func ExecuteTestRequest(t *testing.T, verb string, path string, body string, funcHandler Handler) (*http.Request, *httptest.ResponseRecorder) {

	request, recorder := GetTestRequest(t, verb, path, body)

	handler := http.HandlerFunc(funcHandler)
	handler.ServeHTTP(recorder, request)

	return request, recorder
}
