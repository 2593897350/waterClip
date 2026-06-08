package processor

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{},
	}
}

func (c *Client) Health() error {
	response, err := c.http.Get(c.baseURL + "/health")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var body map[string]string
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return err
	}

	if body["status"] != "ok" {
		return fmt.Errorf("unexpected status %q", body["status"])
	}

	return nil
}
