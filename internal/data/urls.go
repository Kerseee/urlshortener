// Package data provides the url model for executing query from and to the database.
package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	ErrRecordNotFound    = errors.New("record is not found")
	ErrDuplicateShortUrl = errors.New("duplicate unexpired shortened URL")
)

const (
	errMsgViolateUniquePQ string = "pq: duplicate key value violates unique constraint"
)

// URLModel is a wrapper of a db connection pool.
type URLModel struct {
	DB           *sql.DB
	QueryTimeOut time.Duration
}

// URL holds an entry of the table "urls" in the database.
type URL struct {
	ID        int64
	URL       string
	ExpireAt  time.Time
	ShortPath string
}

// Get return a URL instance based on given shortPath.
func (m *URLModel) Get(s string) (*URL, error) {
	// Prepare the query and arguments
	query := `
		SELECT id, url, short_url, expire_at
		FROM urls
		WHERE short_url = $1`
	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeOut)
	defer cancel()

	// Execute the query
	var u URL
	err := m.DB.QueryRowContext(ctx, query, s).Scan(
		&u.ID,
		&u.URL,
		&u.ShortPath,
		&u.ExpireAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &u, nil
}

// Insert inserts a URL into urls table in the database.
func (m *URLModel) Insert(u *URL) error {
	// Prepare the query and arguments.
	query := `
		INSERT INTO urls(url, short_url, expire_at)
		VALUES ($1, $2, $3)
		RETURNING id`
	args := []interface{}{u.URL, u.ShortPath, u.ExpireAt.UTC()}
	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeOut)
	defer cancel()

	// Execute the query.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&u.ID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), errMsgViolateUniquePQ):
			return ErrDuplicateShortUrl
		default:
			return err
		}
	}

	return nil
}

// Update updates a URL in the urls table in the database.
func (m *URLModel) Update(u *URL) error {
	// Prepare the query
	query := `
		UPDATE urls SET url = $1, short_url = $2, expire_at = $3
		WHERE id = $4`
	args := []interface{}{u.URL, u.ShortPath, u.ExpireAt, u.ID}
	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeOut)
	defer cancel()

	// Execute the query
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}
