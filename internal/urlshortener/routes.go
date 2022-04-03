package urlshortener

import "net/http"

// routes creates and returns a http servemux.
func (app *App) routes() http.Handler {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", app.redirect)
	mux.HandleFunc("/api/v1/urls", app.registerURL)
	return mux
}
