package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DetectResult struct {
	MaskPath string `json:"mask_path"`
}

type InpaintResult struct {
	OutputPath string `json:"output_path"`
}

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

func NewWithHTTPClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL: baseURL,
		http:    httpClient,
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

func (c *Client) Detect(sourcePath string) (DetectResult, error) {
	payload := map[string]string{"image_path": sourcePath}
	var result DetectResult
	if err := c.post("/detect", payload, &result); err != nil {
		return DetectResult{}, err
	}
	return result, nil
}

func (c *Client) Inpaint(sourcePath, maskPath, mode string) (InpaintResult, error) {
	payload := map[string]string{
		"image_path": sourcePath,
		"mask_path":  maskPath,
		"mode":       mode,
	}
	var result InpaintResult
	if err := c.post("/inpaint", payload, &result); err != nil {
		return InpaintResult{}, err
	}
	return result, nil
}

func (c *Client) post(path string, payload any, target any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	response, err := c.http.Post(c.baseURL+path, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("processor request failed with status %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return err
	}

	return nil
}
