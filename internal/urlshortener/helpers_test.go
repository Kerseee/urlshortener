package urlshortener

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWriteJson(t *testing.T) {
	tests := []struct {
		name       string
		code       int
		data       envelop
		headers    http.Header
		wantErr    error
		wantCode   int
		wantHeader http.Header
		wantBody   []string
	}{
		{
			name:       "short url response",
			code:       http.StatusOK,
			data:       envelop{"id": "abcd1234", "shortUrl": "http://localhost:8080/abcd1234"},
			headers:    nil,
			wantErr:    nil,
			wantCode:   http.StatusOK,
			wantHeader: http.Header{"Content-Type": []string{"application/json"}},
			wantBody:   []string{`"id": "abcd1234"`, `"shortUrl": "http://localhost:8080/abcd1234"`},
		},
		{
			name:    "error marshal type",
			data:    envelop{"id": make(chan int)},
			headers: nil,
			wantErr: errors.New(""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := writeJSON(w, test.code, test.data, test.headers)

			// Check the error.
			if err == nil && test.wantErr != nil {
				t.Error("want non nil error, got nil error")
			}
			if err != nil && test.wantErr == nil {
				t.Errorf(`want nil error, got "%v"`, err)
			}
			if err != nil && test.wantErr != nil {
				return
			}

			// Check the response status code.
			code, header, body := getResponse(t, w)
			validateCode(t, test.wantCode, code)
			validateHeader(t, test.wantHeader, header)
			for _, wantBody := range test.wantBody {
				validateBodyContains(t, []byte(wantBody), body)
			}
		})
	}

}

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
				t.Errorf(`want error message contains "%v", got nil error`, test.wantErrMsg)
			case err != nil && test.wantErrMsg == "":
				t.Errorf(`want nil error, got "%v"`, err)
			}
		})
	}
}

func TestReadJson(t *testing.T) {
	type urlBody struct {
		Url      string    `json:"url"`
		ExpireAt time.Time `json:"expireAt"`
	}

	tests := []struct {
		name       string
		body       string
		wantErrMsg string
		wantData   interface{}
	}{
		{
			name:       "valid request",
			body:       `{"url":"http://google.com","expireAt":"2023-12-22T12:00:00Z"}`,
			wantErrMsg: "",
			wantData: urlBody{
				Url:      "http://google.com",
				ExpireAt: time.Date(2023, 12, 22, 12, 0, 0, 0, time.UTC),
			},
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Send a request
			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/urls", bytes.NewReader([]byte(test.body)))
			r.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			var u urlBody
			err := readJSON(w, r, &u)

			// Check the response
			switch {
			case err == nil && test.wantErrMsg == "":
				wantData, ok := test.wantData.(urlBody)
				if !ok {
					t.Fatal("type assertion fail")
				}
				if u.Url != wantData.Url {
					t.Errorf(`want data has url "%s", got "%s"`, wantData.Url, u.Url)
				}
				if !u.ExpireAt.Equal(wantData.ExpireAt) {
					t.Errorf("want data has expire time %v, got %v", wantData.ExpireAt, u.ExpireAt)
				}
			case err == nil && test.wantErrMsg != "":
				t.Errorf(`want error message contains "%s", got nil`, test.wantErrMsg)
			case err != nil && test.wantErrMsg == "":
				t.Errorf(`want nil error, got "%v"`, err)
			case err != nil && !strings.Contains(err.Error(), test.wantErrMsg):
				t.Errorf(`want error message contains "%s", got "%v"`, test.wantErrMsg, err)
			}
			r.Body.Close()
		})

	}
}

func TestHashAndEncode(t *testing.T) {
	tests := []struct {
		name string
		str  string
	}{
		{"http url", "http://google.com"},
		{"https url", "https://github.com"},
		{"random string", `1qaz2wsx3edc4rfv5tgb6yhn7ujm8ik9ol0~!@#$%^&*()_+[]{}|\:;'",<.>/?`},
		{"empty string", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded := hashAndEncode(test.str)
			if len(encoded) != 43 {
				t.Fatalf(`want 43-byte-long encoded string, got "%s"`, encoded)
			}
		})
	}
}

func TestValidateExpireTime(t *testing.T) {
	tests := []struct {
		name       string
		time       time.Time
		wantErrMsg string
	}{
		{
			name:       "before now local",
			time:       time.Date(2021, 12, 22, 12, 0, 0, 0, time.Local),
			wantErrMsg: "expired time should after now",
		},
		{
			name:       "after now local",
			time:       time.Date(2023, 12, 22, 12, 0, 0, 0, time.Local),
			wantErrMsg: "",
		},
		{
			name:       "before now UTC",
			time:       time.Date(2021, 12, 22, 12, 0, 0, 0, time.UTC),
			wantErrMsg: "expired time should after now",
		},
		{
			name:       "after now local",
			time:       time.Date(2023, 12, 22, 12, 0, 0, 0, time.UTC),
			wantErrMsg: "",
		},
		{
			name:       "now local",
			time:       time.Now(),
			wantErrMsg: "expired time should after now",
		},
		{
			name:       "now UTC",
			time:       time.Now().UTC(),
			wantErrMsg: "expired time should after now",
		},
		{
			name:       "zero time value",
			time:       time.Time{},
			wantErrMsg: "expired time should after now",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateExpireTime(test.time)
			switch {
			case err == nil && test.wantErrMsg != "":
				t.Errorf(`want error message contains "%s", got nil error`, test.wantErrMsg)
			case err != nil && test.wantErrMsg == "":
				t.Errorf("want nil error, got %v", err)
			case err != nil && !strings.Contains(err.Error(), test.wantErrMsg):
				t.Errorf(`want error message contains "%s", got "%v"`, test.wantErrMsg, err)
			}
		})
	}
}

func TestShortenURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		lenShortURL int
		wantErrMsg  string
	}{
		{
			name:        "valid url",
			url:         "http://google.com",
			lenShortURL: 8,
			wantErrMsg:  "",
		},
		{
			name:        "empty url",
			url:         "",
			lenShortURL: 8,
			wantErrMsg:  "",
		},
		{
			name:        "negative length of short url",
			url:         "http://facebook.com",
			lenShortURL: -1,
			wantErrMsg:  "config.ShortURL.Len out of the range [1, 43]",
		},
		{
			name:        "invalid length of short url",
			url:         "http://github.com",
			lenShortURL: 44,
			wantErrMsg:  "config.ShortURL.Len out of the range [1, 43]",
		},
	}

	app, _ := newTestApp()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app.config.ShortURL.Len = test.lenShortURL
			shortURL, err := app.shortenURL(test.url)
			switch {
			case err == nil && test.wantErrMsg == "":
				if len(shortURL) != test.lenShortURL {
					t.Errorf("want short url %d-byte-long, got %d-byte-long", test.lenShortURL, len(shortURL))
				}
			case err == nil && test.wantErrMsg != "":
				t.Errorf(`want error message contains "%s", got nil error`, test.wantErrMsg)
			case err != nil && test.wantErrMsg == "":
				t.Errorf("want nil error, got %v", err)
			case err != nil && !strings.Contains(err.Error(), test.wantErrMsg):
				t.Errorf(`want error message contains "%s", got "%v"`, test.wantErrMsg, err)
			}
		})
	}
}
