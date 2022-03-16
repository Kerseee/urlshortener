package main

import (
	"log"

	"github.com/Kerseee/urlshortener/config"
	"github.com/Kerseee/urlshortener/internal/urlshortener"
)

func main() {
	cfg := config.Config{Addr: ":8080"}
	app := urlshortener.New(cfg)
	log.Fatal(app.Serve())
}
