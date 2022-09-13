package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONRequest Sends an HTTP Request with some default headers, and writes the
// result into v.
func JSONRequest(req *http.Request, c *http.Client, v interface{}) error {
	// default http client
	if c == nil {
		c = http.DefaultClient
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	// wrap the HTTP error
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("error, status code: %d", res.StatusCode)
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
			return err
		}
	}

	return nil
}
