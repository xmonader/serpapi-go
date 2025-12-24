package serpapi

import (
	"net/http"
	"time"
)

const (
	DefaultBaseURL = "https://serpapi.com"
	DefaultTimeout = 60 * time.Second
	Version        = "1.0.0"
	UserAgent      = "serpapi-go/" + Version
)

// Client is the main entry point for the SerpApi client.
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient creates a new SerpApi client with the given API key and options.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL: DefaultBaseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
