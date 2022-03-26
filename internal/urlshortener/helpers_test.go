package urlshortener

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantErrMsg string
	}{
		{"valid http URL", "http://google.com", ""},
		{"valid https URL", "https://google.com", ""},
		{"unvalid URL", "http/", "invalid url"},
		{"non http URL", "ftp://", "invalid url"},
		{"empty URL", "", "invalid url"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateURL(test.url)
			switch {
			case err == nil && test.wantErrMsg != "":
				t.Errorf("want error message contains %v, got %v", test.wantErrMsg, "")
			case err != nil && test.wantErrMsg == "":
				t.Errorf("want nil error, got %v", err)
			}
		})
	}
}

func TestReadJson(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantErrMsg string
	}{
		{
			name:       "valid request",
			body:       `{"url":"http://google.com","expireAt":"2023-12-22T12:00:00Z"}`,
			wantErrMsg: "",
		},
		{
			name:       "invalid json syntax",
			body:       `{"url":"http://facebook.com"sd,"expireAt":"2023-12-22T12:00:00Z"s}`,
			wantErrMsg: "has syntax error at character",
		},
		{
			name:       "invalid json type",
			body:       `{"url":123, "expireAt":"2023-12-22T12:00:00Z"}`,
			wantErrMsg: "has incorrect type",
		},
		{
			name:       "invalid time format",
			body:       `{"url":"http://github.com", "expireAt":"2023/12/22"}`,
			wantErrMsg: "time format error",
		},
		{
			name:       "empty json",
			body:       "",
			wantErrMsg: "JSON should not be empty",
		},
		{
			name:       "oversize body",
			body:       fmt.Sprintf(`{"url":"http://google.com/%s","expireAt":"2023-12-22T12:00:00Z"}`, strings.Repeat("a", 1<<21)),
			wantErrMsg: "body size should not exceed 1 MB",
		},
		{
			name:       "two json",
			body:       `{"url":"https://youtube.com","expireAt":"2023-12-22T12:00:00Z"}{"url":"https://youtube.com","expireAt":"2023-12-22T12:00:00Z"}`,
			wantErrMsg: "more than 1 JSON in the request",
		},
	}

	// Create a json body.
	type Body struct {
		Url      string    `json:"url"`
		ExpireAt time.Time `json:"expireAt"`
	}

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/urls", bytes.NewReader([]byte(test.body)))
		r.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()

		t.Run(test.name, func(t *testing.T) {
			var v Body
			err := readJSON(w, r, &v)
			switch {
			case err == nil && test.wantErrMsg != "":
				t.Errorf("want error message contains %v, got nil", test.wantErrMsg)
			case err != nil && test.wantErrMsg == "":
				t.Errorf("want nil error, got %v", err)
			case err != nil && !strings.Contains(err.Error(), test.wantErrMsg):
				t.Errorf("want error message contains %v, got %v", test.wantErrMsg, err)
			}
		})
		r.Body.Close()
	}
}
