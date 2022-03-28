package urlshortener

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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

// stringSet creates a set of string containing values in slice s.
func stringSet(s []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

// validateCode check if want equals got.
func validateCode(t *testing.T, want, got int) {
	if want != got {
		t.Errorf("want status code %d, got %d", want, got)
	}
}

// validateHeader check if headers in want are all present in got.
func validateHeader(t *testing.T, want, got http.Header) {
	for header, vals := range want {
		gotVals, ok := got[header]
		if !ok {
			t.Errorf(`miss header %q with values "%v"`, header, gotVals)
			continue
		}
		gotValSet := stringSet(gotVals)
		var miss []string
		for _, v := range vals {
			if _, ok := gotValSet[v]; !ok {
				miss = append(miss, v)
			}
		}
		if len(miss) > 0 {
			t.Errorf(`miss values "%v" in the header %q`, miss, header)
		}
	}
}

// validateBodyContains check if got contains want.
func validateBodyContains(t *testing.T, want, got []byte) {
	if !bytes.Contains(got, want) {
		t.Errorf("want body contains %q, got %q", want, got)
	}
}

// getResponse extracts the status code, header and body from the response recorder w.
// It close the resp.Body after reading the response body.
func getResponse(t *testing.T, w *httptest.ResponseRecorder) (code int, header http.Header, body []byte) {
	resp := w.Result()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return resp.StatusCode, resp.Header, body
}
