// Package mock provides mocked data for testing.
package mock

import (
	"time"

	"github.com/Kerseee/urlshortener/internal/data"
)

// URLModel mocks the data.URLModel.
type URLModel struct{}

// mockURL is a mocked data.URL instance.
var mockURLs = map[string]data.URL{
	"BQRvJsg-": {
		ID:        1,
		URL:       "https://google.com",
		ExpireAt:  time.Date(2022, time.December, 22, 12, 0, 0, 0, time.UTC),
		ShortPath: "BQRvJsg-",
	},
	"FGeTGg6M": {
		ID:        2,
		URL:       "https://youtube.com",
		ExpireAt:  time.Date(2021, time.December, 22, 12, 0, 0, 0, time.UTC),
		ShortPath: "FGeTGg6M",
	},
	"zXWCjacZ": {
		ID:        3,
		URL:       "https://not-netflix.com",
		ExpireAt:  time.Date(2024, time.December, 22, 12, 0, 0, 0, time.UTC),
		ShortPath: "zXWCjacZ",
	},
	"zXWCjacZn": {
		ID:        4,
		URL:       "https://not-netflix.com/1",
		ExpireAt:  time.Date(2024, time.December, 22, 12, 0, 0, 0, time.UTC),
		ShortPath: "zXWCjacZn",
	},
	"zXWCjacZns": {
		ID:        5,
		URL:       "https://not-netflix.com/2",
		ExpireAt:  time.Date(2024, time.December, 22, 12, 0, 0, 0, time.UTC),
		ShortPath: "zXWCjacZns",
	},
	"zXWCjacZnsJ": {
		ID:        6,
		URL:       "https://not-netflix.com/3",
		ExpireAt:  time.Date(2024, time.December, 22, 12, 0, 0, 0, time.UTC),
		ShortPath: "zXWCjacZnsJ",
	},
	"zXWCjacZnsJ4": {
		ID:        7,
		URL:       "https://not-netflix.com/4",
		ExpireAt:  time.Date(2024, time.December, 22, 12, 0, 0, 0, time.UTC),
		ShortPath: "zXWCjacZnsJ4",
	},
}

// Get mocks the data.URLModel.Get method.
func (m *URLModel) Get(s string) (*data.URL, error) {
	u, ok := mockURLs[s]
	if !ok {
		return nil, data.ErrRecordNotFound
	}
	return &u, nil
}

// Insert mocks the data.URLModel.Insert method.
func (m *URLModel) Insert(u *data.URL) error {
	if _, ok := mockURLs[u.ShortPath]; ok {
		return data.ErrDuplicateShortUrl
	}
	return nil
}

func (m *URLModel) Update(u *data.URL) error {
	return nil
}
