package search

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	// Timeout specifies the maximum duration for a single request.
	// Default: 30 seconds
	Timeout time.Duration

	// HTTPClient specifies the client to send HTTP requests.
	// If HTTPClient is set, Timeout will not be used.
	// Optional. Default &http.client{Timeout: Timeout}
	HTTPClient *http.Client `json:"http_client"`

	// MaxResults limits the number of results returned
	// Default: 10
	MaxResults int `json:"max_results"`

	// Region is the geographical region for results
	// Default: RegionWT, means all regions
	// Reference: https://duckduckgo.com/duckduckgo-help-pages/settings/params
	Region Region `json:"region"`
}

func NewSearch(ctx context.Context, config *Config) (Search, error) {
	return buildClient(ctx, config)
}

func buildClient(_ context.Context, config *Config) (Search, error) {
	if config == nil {
		config = &Config{}
	}

	region := config.Region
	if region == "" {
		region = RegionWT
	}

	maxResults := config.MaxResults
	if maxResults <= 0 {
		maxResults = 10
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	var httpCli *http.Client
	if config.HTTPClient != nil {
		httpCli = config.HTTPClient
	} else {
		httpCli = &http.Client{
			Timeout: timeout,
		}
	}

	return &client{
		httpCli:    httpCli,
		maxResults: maxResults,
		region:     region,
	}, nil
}
