package gptj

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// sendRequest Sends an HTTP Request with some default headers, and writes the
// result into v.
func (c *gptjClient) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
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

// gptjClient is a client implementation of GPT-J meant to implement
// the AI Client interface.
type gptjClient struct {
	baseURL        string
	httpClient     *http.Client
	generateParams *GenerateParams
}

// choice Represents a list element which is returned from the GPT-J API.
type choice struct {
	GeneratedText string `json:"generated_text"`
}
