package urlshortener

import (
	"encoding/json"
	"net/http"
)

type envelop map[string]interface{} // wrap the data to be parsed into JSON

// writeJson encodes data into JSON, and writes status, encoded data and headers into a response.
func writeJSON(w http.ResponseWriter, status int, data envelop, headers http.Header) error {
	// Encode the data into JSON.
	body, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Write headers.
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Add("Content-Type", "application/json")

	// Write http status.
	w.WriteHeader(status)

	// Write response body.
	_, err = w.Write(body)

	return err
}
