package urlshortener

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMethodNotAllowedResponse(t *testing.T) {
	want := struct {
		code int
		body string
	}{
		code: http.StatusMethodNotAllowed,
		body: "this method is not allowed",
	}

	tests := []struct {
		name   string
		method string
	}{
		{"Get", http.MethodGet},
		{"Post", http.MethodPost},
		{"Delete", http.MethodDelete},
		{"Patch", http.MethodPatch},
		{"Options", http.MethodOptions},
		{"Put", http.MethodPut},
		{"Head", http.MethodHead},
		{"Trace", http.MethodTrace},
		{"Connect", http.MethodConnect},
	}

	app, _ := newTestApp()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Send a request.
			r := httptest.NewRequest(test.method, "http://localhost:8080/invalid-end-point", nil)
			w := httptest.NewRecorder()
			app.methodNotAllowedResponse(w, r)

			// Check the response.
			validateHTTPJsonResponse(t, w, r, want.code, want.body)
		})
	}
}

func TestBadRequestResponse(t *testing.T) {
	want := struct {
		err  error
		code int
		body string
	}{
		err:  errors.New("this is a bad request"),
		code: http.StatusBadRequest,
		body: "this is a bad request",
	}

	// Send a request.
	app, _ := newTestApp()
	r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewBuffer([]byte("some request")))
	w := httptest.NewRecorder()
	app.badRequestResponse(w, r, want.err)

	// Check the response.
	validateHTTPJsonResponse(t, w, r, want.code, want.body)
}

func TestServerErrorResponse(t *testing.T) {
	want := struct {
		err  error
		log  string
		code int
		body string
	}{
		err:  errors.New("some internal server error"),
		log:  "some internal server error",
		code: http.StatusInternalServerError,
		body: "server cannot process your request now",
	}

	// Send a request.
	app, logger := newTestApp()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	w := httptest.NewRecorder()
	app.serverErrorResponse(w, r, want.err)

	// Check the response.
	validateHTTPJsonResponse(t, w, r, want.code, want.body)
	logMsg, err := io.ReadAll(logger)
	if err != nil {
		t.Fatal(err)
	}
	if msg := string(logMsg); !strings.Contains(msg, want.log) {
		t.Errorf(`want log contains "%s", got "%s"`, want.log, msg)
	}
}

func TestRecordNotFoundResponse(t *testing.T) {
	want := struct {
		code int
		body string
	}{
		code: http.StatusNotFound,
		body: "record not found",
	}

	// Send a request.
	app, _ := newTestApp()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/some-end-point?query=something", nil)
	w := httptest.NewRecorder()
	app.recordNotFoundResponse(w, r)

	// Validate the response
	validateHTTPJsonResponse(t, w, r, want.code, want.body)
}
