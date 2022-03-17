package urlshortener

import (
	"errors"
	"net/http"
)

var (
	// ErrRequestBodyTooLarge describe the error in http.MaxBytesReader
	ErrRequestBodyTooLarge = errors.New("http: request body too large")
)

// InternalError wrap an error with customized error message Msg and origin error Err.
type InternalError struct {
	Msg string // customized message
	Err error  // origin error
}

func (e *InternalError) Error() string {
	return e.Err.Error()
}

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
	app.logger.Println(err.Error())

	msg := envelop{"error": errors.New("server cannot process your request now")}
	writeJSON(w, http.StatusInternalServerError, msg, nil)
}
