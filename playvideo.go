// Package playvideo provides an official Go SDK for the PlayVideo API.
//
// PlayVideo is a video hosting platform for developers. This SDK provides
// a simple interface to interact with the PlayVideo API.
//
// Basic usage:
//
//	client := playvideo.NewClient("play_live_xxx")
//
//	// List collections
//	collections, err := client.Collections.List(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Upload a video
//	upload, err := client.Videos.UploadFile(context.Background(), "./video.mp4", "my-collection", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Watch processing progress
//	for event := range client.Videos.WatchProgress(context.Background(), upload.Video.ID) {
//	    fmt.Printf("%s: %s\n", event.Stage, event.Message)
//	}
package playvideo

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default API base URL
	DefaultBaseURL = "https://api.playvideo.dev/api/v1"
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second
)

// Client is the PlayVideo API client
type Client struct {
	// API key for authentication
	apiKey string

	// Base URL for API requests
	baseURL string

	// HTTP client used for requests
	httpClient *http.Client

	// Internal HTTP client wrapper
	http *httpClient

	// API Resources
	Collections *CollectionsService
	Videos      *VideosService
	Webhooks    *WebhooksService
	Embed       *EmbedService
	APIKeys     *APIKeysService
	Account     *AccountService
	Usage       *UsageService
}

// ClientOption is a function that configures the client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new PlayVideo API client
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Initialize internal HTTP client
	c.http = newHTTPClient(c.apiKey, c.baseURL, c.httpClient)

	// Initialize services
	c.Collections = &CollectionsService{client: c.http}
	c.Videos = &VideosService{client: c.http}
	c.Webhooks = &WebhooksService{client: c.http}
	c.Embed = &EmbedService{client: c.http}
	c.APIKeys = &APIKeysService{client: c.http}
	c.Account = &AccountService{client: c.http}
	c.Usage = &UsageService{client: c.http}

	return c
}
