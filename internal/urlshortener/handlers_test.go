package urlshortener

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	app, _ := newTestApp()
	tests := []struct {
		name     string
		method   string
		shortURL string
		wantCode int
		wantBody string
	}{
		{
			name:     "valid path",
			method:   http.MethodGet,
			shortURL: "http://localhost:8080/abcd1234",
			wantCode: http.StatusSeeOther,
			wantBody: "http://google.com",
		},
		{
			name:     "not exist path",
			method:   http.MethodGet,
			shortURL: "http://localhost:8080/abcd1236",
			wantCode: http.StatusNotFound,
			wantBody: "record not found or expired",
		},
		{
			name:     "invalid method",
			method:   http.MethodPost,
			shortURL: "http://localhost:8080/abcd1234",
			wantCode: http.StatusMethodNotAllowed,
			wantBody: "this method is not allowed",
		},
		{
			name:     "record expired",
			method:   http.MethodGet,
			shortURL: "http://localhost:8080/abcd1235",
			wantCode: http.StatusNotFound,
			wantBody: "record not found or expired",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Send a request.
			r := httptest.NewRequest(test.method, test.shortURL, nil)
			w := httptest.NewRecorder()
			app.redirect(w, r)

			// Extract the response.
			code, _, body := getResponse(t, w)

			// Validate the response.
			validateCode(t, test.wantCode, code)
			validateBodyContains(t, []byte(test.wantBody), body)
		})
	}
}
