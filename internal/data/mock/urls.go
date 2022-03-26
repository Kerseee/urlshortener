package mock

import (
	"time"

	"github.com/Kerseee/urlshortener/internal/data"
)

// URLModel mocks the data.URLModel.
type URLModel struct{}

// mockURL is a mocked data.URL instance.
var mockURLs = map[string]data.URL{
	"abcd1234": {
		ID:        1,
		URL:       "http://google.com",
		ExpireAt:  time.Date(2022, time.December, 22, 12, 0, 0, 0, time.UTC),
		ShortPath: "abcd1234",
	},
	"abcd1235": {
		ID:        2,
		URL:       "http://github.com",
		ExpireAt:  time.Date(2021, time.December, 22, 12, 0, 0, 0, time.UTC),
		ShortPath: "abcd1235",
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
	return nil
}

func (m *URLModel) Update(u *data.URL) error {
	return nil
}
