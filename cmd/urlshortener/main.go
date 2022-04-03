package main

import (
	"log"

	"github.com/Kerseee/urlshortener/config"
	"github.com/Kerseee/urlshortener/internal/urlshortener"
)

func main() {
	cfg := config.New()
	app, err := urlshortener.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(app.Serve())
}
