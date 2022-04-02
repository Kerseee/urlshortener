/*
Package config provides a configuration struct used in the url shortener application.

Before creating the url shortener application, please create a Config first:

	conf := Config.New()

This will parse all flags and store the configuration into a Config struct.
*/

package config

import (
	"flag"
	"os"
	"time"
)

// Config contains all configuration used in the application,
// including the address of the server, configurations of database connection pool
// and other customized configurations.
type Config struct {
	// Addr is the address or domain name and the port of this server.
	//
	// Example: localhost:8080
	Addr string

	// DB holds the settings of the database connection pool.
	DB struct {
		DSN          string        // dsn of the database, like "postgres://username:password@localhost/database"
		MaxOpenConns int           // number of maximum open connections
		MaxIdleConns int           // number of maximum idle connections
		MaxIdleTime  int           // maximum idle time of a connection (minutes)
		QueryTimeout time.Duration // maximum time for executing query to the database (seconds)
	}
	ShortURL struct {
		Len int // length of shortened URL

		// MaxReShortenLen is the maximum length of shortened URL
		// for trying re-shorten URL in case of short URL conflicts.
		//
		// It should be strictly larger than Len.
		MaxReShortenLen int
	}
}

// New parses the flags, store all config into a config.Config and returns.
func New() Config {
	var conf Config
	flag.StringVar(&conf.Addr, "addr", "localhost:8080", "Server address (hostname:port)")

	flag.StringVar(&conf.DB.DSN, "db", os.Getenv("URLSHORTENER_DB_DSN"), "Database dsn")
	flag.IntVar(&conf.DB.MaxOpenConns, "db-max-open-conns", 25, "Database maximum open connections")
	flag.IntVar(&conf.DB.MaxIdleConns, "db-max-idle-conns", 25, "Database maximum idle connections")
	flag.IntVar(&conf.DB.MaxIdleTime, "db-max-idle-time", 15, "Database maximum idle time (minutes)")
	queryTimeOut := flag.Int("db-query-timeout", 3, "Database maximum query time (seconds)")
	conf.DB.QueryTimeout = time.Second * time.Duration(*queryTimeOut)

	flag.IntVar(&conf.ShortURL.Len, "len-short-url", 8, "Length of shortened URL (should be greater than 4 and less than 17)")
	flag.IntVar(&conf.ShortURL.MaxReShortenLen, "max-len-reshort-url", 12, "Maximum length of shortened URL for reshortening URL in case of short URL conflicts, should be greater or equal than len-short-url and less than 44")

	flag.Parse()

	conf.Validate()
	return conf
}

// Validate validates the config and automatically adjusts the config to default setting
// if the config is not valid.
func (conf *Config) Validate() {
	if conf.ShortURL.Len <= 4 || conf.ShortURL.Len >= 17 {
		conf.ShortURL.Len = 8
	}
	if conf.ShortURL.MaxReShortenLen < conf.ShortURL.Len || conf.ShortURL.MaxReShortenLen >= 44 {
		conf.ShortURL.MaxReShortenLen = conf.ShortURL.Len + 4
	}
}
