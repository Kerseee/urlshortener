package urlshortener

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

type envelop map[string]interface{} // wrap the data to be parsed into JSON

var validURLExp = regexp.MustCompile(`^https?:\/\/`)

// writeJson encodes data into JSON, and writes status, encoded data and headers into a response.
func writeJSON(w http.ResponseWriter, status int, data envelop, headers http.Header) error {
	// Encode the data into JSON.
	body, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Write headers.
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Add("Content-Type", "application/json")

	// Write http status.
	w.WriteHeader(status)

	// Write response body.
	_, err = w.Write(body)

	return err
}

// readJSON reads the JSON-encoded request body with the body size limited to 1MB,
// decodes the JSON and stores it into a instance pointed to by v.
func readJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// Prepare a JSON decoder.
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode the body.
	err := decoder.Decode(&v)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var unMarshalTypeErr *json.UnmarshalTypeError
		var invalidUnmarshalErr *json.InvalidUnmarshalError
		var timeParsingErr *time.ParseError
		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("has syntax error at character %d in JSON", syntaxErr.Offset)
		case errors.As(err, &unMarshalTypeErr):
			return fmt.Errorf("has incorrect type at character %d in JSON", unMarshalTypeErr.Offset)
		case errors.As(err, &timeParsingErr):
			return fmt.Errorf("time format error")
		case errors.Is(err, io.EOF):
			return errors.New("JSON should not be empty")
		case err.Error() == ErrRequestBodyTooLarge.Error():
			return errors.New("body size should not exceed 1 MB")
		case errors.As(err, &invalidUnmarshalErr):
			return &InternalError{Msg: "JSON decoding error", Err: err}
		default:
			return err
		}
	}

	// Prevent the case that the request contains more than 1 JSON.
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("more than 1 JSON in the request")
	}
	return nil
}

// hashAndEncode uses sha256 hash the string s and then encodes it with base64.
// It returns a 44-byte-long string.
func hashAndEncode(s string) string {
	hash := sha256.Sum256([]byte(s))
	encoded := base64.RawURLEncoding.EncodeToString(hash[:])
	return encoded
}

// validateURL returns error if s is not an URL.
func validateURL(s string) error {
	match := validURLExp.MatchString(s)
	if !match {
		return errors.New("invalid url")
	}
	return nil
}

// validateExpireTime returns error if t is before now.
func validateExpireTime(t time.Time) error {
	if t.Before(time.Now()) {
		return errors.New("expired time should after now")
	}
	return nil
}

// shortenURL shortens s into 8 bytes long string.
func (app *App) shortenURL(s string) string {
	return hashAndEncode(s)[:app.config.ShortURL.Len]
}
