package urlshortener

import (
	"net/http"
)

// methodNotAllowedResponse informs the client that the method of this request is not allowed.
func (app *App) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := envelop{"error": "this method is not allowed"}
	writeJSON(w, http.StatusMethodNotAllowed, msg, nil)
}

// badRequestResponse informs the client of bad request.
func (app *App) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg := envelop{"error": err.Error()}
	writeJSON(w, http.StatusBadRequest, msg, nil)
}

// serverErrorResponse informs the client of server internal error.
func (app *App) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg := envelop{"error": err.Error()}
	writeJSON(w, http.StatusInternalServerError, msg, nil)
}
