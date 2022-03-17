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
	switch {
	case input.Url == "":
		err := errors.New("url should not be empty")
		app.badRequestResponse(w, r, err)
		return
	case input.ExpireAt.Before(time.Now()):
		err := errors.New("expire time should be provided and be after now")
		app.badRequestResponse(w, r, err)
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
