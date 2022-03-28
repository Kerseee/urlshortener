package urlshortener

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
			r := httptest.NewRequest(test.method, test.shortURL, nil)
			w := httptest.NewRecorder()
			app.redirect(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != test.wantCode {
				t.Errorf("want status code %d, got %d", test.wantCode, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
				return
			}
			if !strings.Contains(string(body), test.wantBody) {
				t.Errorf("want body contains %s, got %s", test.wantBody, string(body))
			}
		})
	}
}
