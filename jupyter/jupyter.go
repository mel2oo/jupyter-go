package jupyter

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Client is the main client for interacting with the Jupyter Server API.
type Client struct {
	client  *http.Client
	baseURL *url.URL
	token   string

	Kernels  *KernelService
	Sessions *SessionService
	Contents *ContentService
}

type service struct {
	client *Client
}

// NewClient creates a new Jupyter Server API client.
// host should be the base URL of the Jupyter server, e.g., http://localhost:8888
// token is the authentication token. If empty, no authentication will be used.
func NewClient(opts ...Option) (*Client, error) {
	client := &Client{
		client: &http.Client{},
	}

	client.baseURL, _ = url.Parse("http://localhost:8888")
	client.Contents = &ContentService{client}

	for _, o := range opts {
		o(client)
	}

	return client, nil
}

type Option func(*Client)

func WithBaseURL(baseurl string) Option {
	return func(c *Client) {
		url, err := url.Parse(baseurl)
		if err == nil {
			c.baseURL = url
		}
	}
}

func WithAuthToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

func (c *Client) NewRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+c.token)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		return resp, json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, nil
}

// VersionResponse defines the structure for the /api/ version endpoint.
type VersionResponse struct {
	Version string `json:"version"`
}

// GetVersion fetches the Jupyter Server version.
func (c *Client) GetVersion(ctx context.Context) (*VersionResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "/api", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req, &VersionResponse{})
	if err != nil {
		return nil, err
	}
	var versionResp VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		return nil, err
	}
	return &versionResp, nil
}
