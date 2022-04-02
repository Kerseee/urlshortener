package urlshortener

import (
	"bytes"
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
			shortURL: "http://localhost:8080/BQRvJsg-",
			wantCode: http.StatusSeeOther,
			wantBody: "https://google.com",
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
			shortURL: "http://localhost:8080/FGeTGg6M",
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
			validateBodyContains(t, test.wantBody, string(body))
		})
	}
}

func TestRegisterURL(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantCode   int
		wantHeader http.Header
		wantBody   []string
	}{
		{
			name:       "valid request",
			method:     http.MethodPost,
			body:       `{"url":"https://facebook.com", "expireAt":"2023-12-22T12:00:00Z"}`,
			wantCode:   http.StatusOK,
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []string{"id", "shortUrl", "localhost:8080/"},
		},
		{
			name:       "invalid method",
			method:     http.MethodGet,
			body:       "",
			wantCode:   http.StatusMethodNotAllowed,
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []string{"error"},
		},
		{
			name:       "body syntax error",
			method:     http.MethodPost,
			body:       `some non json request`,
			wantCode:   http.StatusBadRequest,
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []string{"error"},
		},
		{
			name:       "request lack of expire time",
			method:     http.MethodPost,
			body:       `{"url":"https://facebook.com"}`,
			wantCode:   http.StatusBadRequest,
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []string{"error"},
		},
		{
			name:       "invalid URL",
			method:     http.MethodPost,
			body:       `{"url":"httpp/foo", "expireAt":"2023-12-22T12:00:00Z"}`,
			wantCode:   http.StatusBadRequest,
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []string{"error"},
		},
		{
			name:       "unknown field",
			method:     http.MethodPost,
			body:       `{"url":"https://facebook.com", "expireAt":"2023-12-22T12:00:00Z", "user":"userA"}`,
			wantCode:   http.StatusBadRequest,
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []string{"error"},
		},
	}

	app, _ := newTestApp()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Send a request.
			r := httptest.NewRequest(test.method, "http://localhost:8080/api/v1/urls", bytes.NewBuffer([]byte(test.body)))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			app.registerURL(w, r)

			// Extract the response.
			code, header, body := getResponse(t, w)

			// Validate the response.
			validateCode(t, test.wantCode, code)
			validateHeader(t, test.wantHeader, header)
			for _, wantBody := range test.wantBody {
				validateBodyContains(t, wantBody, string(body))
			}
		})
	}
}
