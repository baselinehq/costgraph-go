package costgraph

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/baselinehq/costgraph-go/internal/api"
)

type Client struct {
	api *api.Client
}

type Option func(*options)

type options struct {
	httpClient *http.Client
}

func WithHTTPClient(client *http.Client) Option {
	return func(o *options) {
		o.httpClient = client
	}
}

func New(apiURL, apiKey string, opts ...Option) (*Client, error) {
	endpoint, err := url.Parse(apiURL)
	if err != nil || endpoint == nil || endpoint.Host == "" || endpoint.Opaque != "" || (endpoint.Scheme != "https" && endpoint.Scheme != "http") || endpoint.User != nil || endpoint.RawQuery != "" || endpoint.ForceQuery || endpoint.Fragment != "" {
		return nil, errors.New("api_url must be an absolute HTTP or HTTPS deployment URL without credentials, query, or fragment")
	}
	if strings.TrimSpace(apiKey) == "" || strings.ContainsAny(apiKey, "\r\n") {
		return nil, errors.New("api_key must be nonempty and contain no line breaks")
	}

	config := options{httpClient: &http.Client{Timeout: 30 * time.Second}}
	for _, option := range opts {
		if option != nil {
			option(&config)
		}
	}
	if config.httpClient == nil {
		return nil, errors.New("HTTP client must not be nil")
	}
	transport := *config.httpClient
	transport.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}

	baseURL := strings.TrimRight(endpoint.String(), "/") + "/"
	generated, err := api.NewClient(baseURL, api.WithHTTPClient(&transport), api.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-API-Key", apiKey)
		req.Header.Set("User-Agent", "costgraph-go")
		return nil
	}))
	if err != nil {
		return nil, errors.New("create management client: invalid deployment URL")
	}
	return &Client{api: generated}, nil
}
