package urlshortener

import (
	"log"
	"net/http"

	"github.com/Kerseee/urlshortener/config"
	"github.com/Kerseee/urlshortener/internal/data"
	"github.com/Kerseee/urlshortener/internal/data/mock"
)

type App struct {
	config   config.Config
	logger   *log.Logger
	urlModel interface {
		Get(s string) (*data.URL, error)
		Insert(u *data.URL) error
	}
}

// New creates and returns an application instance.
func New(cfg config.Config) *App {
	return &App{
		config:   cfg,
		logger:   log.Default(),
		urlModel: &mock.URLModel{},
	}
}

// Serve opens a http server and serves http requests.
func (app *App) Serve() error {
	app.logger.Printf("Start server at %s\n", app.config.Addr)
	server := &http.Server{
		Addr:    app.config.Addr,
		Handler: app.routes(),
	}
	return server.ListenAndServe()
}
