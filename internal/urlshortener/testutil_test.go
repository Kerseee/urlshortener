package urlshortener

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kerseee/urlshortener/config"
	"github.com/Kerseee/urlshortener/internal/data/mock"
)

// newTestServer returns a pointer point to an App instance and a bytes.Buffer as logger.
func newTestApp() (*App, *bytes.Buffer) {
	conf := config.Config{
		Addr: "http://localhost:8080",
		ShortURL: struct {
			Len             int
			MaxReShortenLen int
		}{
			Len:             8,
			MaxReShortenLen: 12,
		},
	}

	logger := bytes.Buffer{}
	return &App{
		config:   conf,
		logger:   log.New(&logger, "", 0),
		urlModel: &mock.URLModel{},
	}, &logger
}

// validateHTTPJsonResponse validates the populated responseRecorder w and Request r with wantCode and wantBody.
//
// It checks if the response has "Content-Type" header with value as "application/json".
//
// It will close w.Result().Body before return.
func validateHTTPJsonResponse(t *testing.T, w *httptest.ResponseRecorder, r *http.Request, wantCode int, wantBody string) {
	wantContentType := "application/json"

	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != wantCode {
		t.Errorf("want status code %d, got %d", wantCode, resp.StatusCode)
	}
	if contentType := resp.Header.Get("Content-Type"); contentType != wantContentType {
		t.Errorf(`want header "Content-Type" has value "%s", got "%s"`, wantContentType, contentType)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if b := string(body); !strings.Contains(b, wantBody) {
		t.Errorf(`want body contains "%s", got "%s"`, wantBody, b)
	}
}
