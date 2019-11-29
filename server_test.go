package nozzle

import (
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test that OptionsHandler defined in server.go is being invoked for any OPTIONS request
func TestOptionsHandling(t *testing.T) {

	handler := http.HandlerFunc(OptionsHandler)
	reader := strings.NewReader("")

	requiredHeaders := map[string]string{
		headers.Origin:					     "localhost",
		headers.AccessControlRequestMethod:  http.MethodPost,
		headers.AccessControlRequestHeaders: headers.ContentType}

	noHeaders := make(map[string]string)

	cases := map[string]struct {
		headers map[string]string
		expectedStatus int
	}{
		"OK:":{requiredHeaders, http.StatusOK},
		"Missing origin":{ noHeaders, http.StatusBadRequest},
	}

	for name, testCase := range cases {
		t.Run(name, func(t *testing.T){
			request, err := http.NewRequest(http.MethodOptions, `/`, reader)
			if err != nil {
				t.Fatal(err)
			}

			for key, value := range testCase.headers {
				request.Header.Add(key, value)
			}
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			actualStatus := recorder.Code
			assert.Equal(t, testCase.expectedStatus, actualStatus)
		})
	}

}
