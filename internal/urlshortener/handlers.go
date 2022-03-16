package urlshortener

import (
	"encoding/json"
	"errors"
	"io"
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

	// Define the input.
	var input struct {
		Url      string    `json:"url"`
		ExpireAt time.Time `json:"expireAt"`
	}

	// Prepare a JSON decoder.
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode the body.
	err := decoder.Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Prevent the case that the request contains more than 1 JSON.
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		app.badRequestResponse(w, r, errors.New("more than 1 JSON in the request"))
		return
	}

	// Write the json back for now.
	data := envelop{
		"your-request": input,
		"todo":         "shorten the url",
	}
	err = writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
