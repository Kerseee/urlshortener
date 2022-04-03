/*
Package urlshortener provides the url shortener application.

The url shortener application has two APIs. One is to shorten a url,
and the other is to redirect the shortened URL to the origin url.

End point "/api/v1/urls" handles json-encoded POST requests and shorten urls.
End point "/:shortenedURL" handles GET requests and redirect to the origin url.

To create a url shortener application:

	import (
		"log"

		"github.com/Kerseee/urlshortener/config"
		"github.com/Kerseee/urlshortener/internal/urlshortener"
	)

	func main() {
		conf := config.New()
		app, err := urlshortener.New(conf)
		if err != nil {
			// handle error
		}
		log.Fatal(app.Serve())
	}
*/

package urlshortener

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kerseee/urlshortener/config"
	"github.com/Kerseee/urlshortener/internal/data"
)

// An App is a url shortener application.
type App struct {
	config config.Config // see package config
	logger *log.Logger

	// A urlModel is a model for executing queries to the urls table in the DB.
	urlModel interface {
		Get(s string) (*data.URL, error)
		Insert(u *data.URL) error
		Update(u *data.URL) error
	}
}

// New creates and returns an application instance including opened database connection pool.
func New(conf config.Config) (*App, error) {
	db, err := OpenDB(conf)
	if err != nil {
		return nil, err
	}
	app := &App{
		config:   conf,
		logger:   log.Default(),
		urlModel: &data.URLModel{DB: db, QueryTimeOut: conf.DB.QueryTimeout},
	}
	app.logInfo("Database connection established!")
	return app, nil
}

// Serve opens a http server and serves http requests.
func (app *App) Serve() error {
	app.logInfo(fmt.Sprintf("Start server at %s\n", app.config.Addr))
	server := &http.Server{
		Addr:    app.config.Addr,
		Handler: app.routes(),
	}
	return server.ListenAndServe()
}

// OpenDB creates a database connection pool and executes first ping for checking connections.
func OpenDB(conf config.Config) (*sql.DB, error) {
	// Create a database connection pool.
	db, err := sql.Open("postgres", conf.DB.DSN)
	if err != nil {
		return nil, err
	}

	// Configure the database connection pool.
	db.SetMaxOpenConns(conf.DB.MaxOpenConns)
	db.SetMaxIdleConns(conf.DB.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Minute * time.Duration(conf.DB.MaxIdleTime))

	// Check connections.
	ctx, cancel := context.WithTimeout(context.Background(), conf.DB.QueryTimeout)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
