package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const BaseURL = "https://api.spotify.com/v1"

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Doer
}

func New(doer Doer) *Client {
	return &Client{Doer: doer}
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("unexpected status code: %d %s", e.Status, e.Message)
}

func (c *Client) GetJSON(ctx context.Context, method, path string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, endpoint(path), nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorResponse := &ErrorResponse{
			Status:  resp.StatusCode,
			Message: "unexpected status code",
		}
		// Attempt to decode the error response, but don't worry if it fails.
		if err = json.NewDecoder(resp.Body).Decode(errorResponse); err != nil {
			log.Printf("failed to decode error response: %v", err)
		}
		return errorResponse
	}

	return json.NewDecoder(resp.Body).Decode(response)
}

func endpoint(path string) string {
	return BaseURL + path
}
