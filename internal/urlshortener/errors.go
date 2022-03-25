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
	err := writeJSON(w, http.StatusMethodNotAllowed, msg, nil)
	if err != nil {
		app.logError(err)
	}
}

// badRequestResponse informs the client of bad request.
func (app *App) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	msg := envelop{"error": err.Error()}
	err = writeJSON(w, http.StatusBadRequest, msg, nil)
	if err != nil {
		app.logError(err)
	}
}

// serverErrorResponse informs the client of server internal error.
func (app *App) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Println(err.Error())

	msg := envelop{"error": "server cannot process your request now"}
	err = writeJSON(w, http.StatusInternalServerError, msg, nil)
	if err != nil {
		app.logError(err)
	}
}

// recordNotFoundResponse informs the client that the requested record is not found
func (app *App) recordNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := envelop{"error": "record not found or expired"}
	err := writeJSON(w, http.StatusNotFound, msg, nil)
	if err != nil {
		app.logError(err)
	}
}
