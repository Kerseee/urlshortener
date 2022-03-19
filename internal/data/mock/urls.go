package mock

import (
	"time"

	"github.com/Kerseee/urlshortener/internal/data"
)

// URLModel mocks the data.URLModel.
type URLModel struct{}

// mockURL is a mocked data.URL instance.
var mockURL = data.URL{
	ID:        1,
	URL:       "http://google.com",
	ExpireAt:  time.Date(2022, time.December, 22, 12, 0, 0, 0, time.Local),
	ShortPath: "qiI5wXYJ",
}

// Get mocks the data.URLModel.Get method.
func (m *URLModel) Get(s string) (*data.URL, error) {
	if s != mockURL.ShortPath {
		return nil, data.ErrRecordNotFound
	}
	return &mockURL, nil
}

// Insert mocks the data.URLModel.Insert method.
func (m *URLModel) Insert(u *data.URL) error {
	if u.ShortPath == mockURL.ShortPath && time.Now().Before(mockURL.ExpireAt) {
		return data.ErrDuplicateShortUrl
	}
	u.ID = 2
	return nil
}
