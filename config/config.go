package config

import "flag"

type Config struct {
	Addr string
}

// New parses the flags, store all config into a config.Config and returns.
func New() Config {
	addr := flag.String("addr", "localhost:8080", "Server address (hostname:port)")
	flag.Parse()
	return Config{
		Addr: *addr,
	}
}
