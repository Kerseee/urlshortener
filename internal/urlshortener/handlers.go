package urlshortener

import (
	"errors"
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
	shortPath := app.shortenURL(input.URL)

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
			// Get the record that has the same shortUrl.
			record, recordErr := app.urlModel.Get(shortPath)
			if recordErr != nil {
				app.serverErrorResponse(w, r, recordErr)
				return
			}

			// If the origin URL does not equal record.URL, then reshorten the URL.
			if record.URL != u.URL {
				app.reShortenURL(w, r, &u)
				return
			}

			// Otherwise, check the expire time.
			// If the expire time is later than record's expire time, then update it.
			if record.ExpireAt.Before(u.ExpireAt) {
				u.ID = record.ID
				err := app.urlModel.Update(&u)
				if err != nil {
					app.serverErrorResponse(w, r, err)
					return
				}
			}
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// Write the short URL back.
	app.writeShortURL(w, r, u.ShortPath)
}

// redirectHandler extracts the shortened URL in the request and redirects to the corresponding origin URL.
// If the shortened URL is not found or is found but expired, then send 404 not found to the client.
func (app *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the method is allowed.
	if r.Method != http.MethodGet {
		app.methodNotAllowedResponse(w, r)
		return
	}

	// Extracts the URL instance.
	path := strings.TrimPrefix(r.URL.Path, "/")
	u, err := app.urlModel.Get(path)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.recordNotFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Check if the URL is expired.
	if u.ExpireAt.Before(time.Now()) {
		app.recordNotFoundResponse(w, r)
		return
	}

	// Redirect to the origin URL.
	http.Redirect(w, r, u.URL, http.StatusSeeOther)
}

// reShortenUrl re-shortens the URL in u.
//
// This method is called in registerURL in case of short URL conflict.
// reShortenURL keep adding 1 character to the short URL and trying to insert into the database.
// The range of the length of short URLs are from app.config.Short.Len + 1 to app.config.ShortURL.MaxReShortenLen.
func (app *App) reShortenURL(w http.ResponseWriter, r *http.Request, u *data.URL) {
	// Hash and encodes the origin URL.
	encodedURL := hashAndEncode(u.URL)

	// Try inserting the shortened URL by adding 1 charachter each time.
	for i := app.config.ShortURL.Len + 1; i <= app.config.ShortURL.MaxReShortenLen+1; i++ {
		u.ShortPath = encodedURL[:i]
		err := app.urlModel.Insert(u)
		if err == nil {
			app.writeShortURL(w, r, u.ShortPath)
			return
		}
		if !errors.Is(err, data.ErrDuplicateShortUrl) {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	app.serverErrorResponse(w, r, errors.New("server internal error: short URL conflict"))
}

// writeShortURL transform the shortPath into a valid short URL and writes the short URL to client.
func (app *App) writeShortURL(w http.ResponseWriter, r *http.Request, shortPath string) {
	shortURL := &url.URL{
		Scheme: "http",
		Host:   app.config.Addr,
		Path:   shortPath,
	}
	data := envelop{
		"id":       shortPath,
		"shortUrl": shortURL.String(),
	}
	err := writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logError(err)
	}
}
