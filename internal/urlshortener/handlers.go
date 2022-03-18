package urlshortener

import (
	"errors"
	"net/http"
	"time"
)

const maxRequestBody int64 = 1 << 20 // 1MB

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// registerURL extracts the to-shorten url from the request, shortens the url,
// and writes the shortened url into response.
func (app *App) registerURL(w http.ResponseWriter, r *http.Request) {
	// Check if the method is allowed.
	if r.Method != http.MethodPost {
		app.methodNotAllowedResponse(w, r)
		return
	}

	// Read the request body.
	var input struct {
		Url      string    `json:"url"`
		ExpireAt time.Time `json:"expireAt"`
	}
	err := readJSON(w, r, &input)
	if err != nil {
		var internalErr *InternalError
		switch {
		case errors.As(err, &internalErr):
			app.serverErrorResponse(w, r, err)
		default:
			app.badRequestResponse(w, r, err)
		}
		return
	}

	// Validate input.
	var errs []string
	if err := validateURL(input.Url); err != nil {
		errs = append(errs, err.Error())
	}
	if err := validateExpireTime(input.ExpireAt); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		writeJSON(w, http.StatusBadRequest, envelop{"error": errs}, nil)
		return
	}

	// Shorten the url.
	shortUrl := shortenURL(input.Url)

	// Write the json back for now.
	data := envelop{
		"your-request": input,
		"shorten-url":  shortUrl,
		"todo":         "shorten the url",
	}
	err = writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
