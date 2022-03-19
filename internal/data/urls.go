package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound    = errors.New("record is not found")
	ErrDuplicateShortUrl = errors.New("duplicate unexpired shortened URL")
)

// URLModel is a wrapper of a db connection pool.
type URLModel struct {
	DB *sql.DB
}

// URL holds an entry of the table "urls" in the database.
type URL struct {
	ID        int64
	URL       string
	ExpireAt  time.Time
	ShortPath string
}
