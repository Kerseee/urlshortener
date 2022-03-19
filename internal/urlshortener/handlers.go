package urlshortener

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Kerseee/urlshortener/internal/data"
)

const maxRequestBody int64 = 1 << 20 // 1MB

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
		URL      string    `json:"url"`
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
	if err := validateURL(input.URL); err != nil {
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
	shortPath := shortenURL(input.URL)
	shortURL := &url.URL{
		Scheme: "http",
		Host:   app.config.Addr,
		Path:   shortPath,
	}

	// Insert the url.
	u := data.URL{
		URL:       input.URL,
		ExpireAt:  input.ExpireAt,
		ShortPath: shortPath,
	}
	err = app.urlModel.Insert(&u)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateShortUrl):
			msg := fmt.Errorf("this url has already been shorten, please use the following shortURL: %s", shortURL.String())
			app.badRequestResponse(w, r, msg)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Write the json back for now.
	data := envelop{
		"your-request": input,
		"id":           shortPath,
		"shortUrl":     shortURL.String(),
	}
	err = writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// redirectHandler extracts the shortened URL in the request and redirects to the corresponding origin URL.
// If the shortened URL is not found in the database, then response 404 not found to the client.
func (app *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the method is allowed.
	if r.Method != http.MethodGet {
		app.methodNotAllowedResponse(w, r)
		return
	}

	// Extracts the URL
	path := strings.TrimPrefix(r.URL.Path, "/")
	u, err := app.urlModel.Get(path)
	if err != nil {
		app.recordNotFoundResponse(w, r)
		return
	}

	// Redirect to the origin URL.
	http.Redirect(w, r, u.URL, http.StatusSeeOther)
}
