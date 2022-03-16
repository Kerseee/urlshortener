package urlshortener

import (
	"log"
	"net/http"

	"github.com/Kerseee/urlshortener/config"
)

type App struct {
	config config.Config
	logger *log.Logger
}

// New creates and returns an application instance.
func New(cfg config.Config) *App {
	return &App{config: cfg, logger: log.Default()}
}

// Serve opens a http server and serves http requests.
func (app *App) Serve() error {
	server := &http.Server{
		Addr:    app.config.Addr,
		Handler: app.routes(),
	}
	return server.ListenAndServe()
}
