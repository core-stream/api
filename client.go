package corestream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.core.stream"
	userAgent      = "corestream-go/1.0"
)

// Client is the core.stream API client.
type Client struct {
	baseURL    *url.URL
	token      string
	httpClient HTTPClient
}

// Option is a functional option for configuring the client.
type Option func(*Client) error

// NewClient creates a new core.stream API client.
// The token is required for authentication.
func NewClient(token string, opts ...Option) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("corestream: token is required")
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("corestream: invalid base URL: %w", err)
		}
		c.baseURL = u
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("corestream: HTTP client cannot be nil")
		}
		c.httpClient = httpClient
		return nil
	}
}

// request performs an HTTP request to the API.
func (c *Client) request(ctx context.Context, method, path string, query url.Values, body, result interface{}) error {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return fmt.Errorf("corestream: invalid path %q: %w", path, err)
	}

	log.Println("request", method, u.String())

	if query != nil {
		u.RawQuery = query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("corestream: failed to encode request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("corestream: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("corestream: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("corestream: failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if len(respBody) > 0 {
			var errResp struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			if json.Unmarshal(respBody, &errResp) == nil {
				apiErr.Code = errResp.Error.Code
				apiErr.Message = errResp.Error.Message
			}
		}
		return apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("corestream: failed to decode response: %w", err)
		}
	}

	return nil
}
