package search

import (
	"context"
	"fkmcps/constants"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func NewDuckDuckGoSearch(ctx context.Context) (Search, error) {
	// 1. Get custom proxy environment variable
	proxyStr := os.Getenv(constants.MCP_PROXY_URL)

	var proxyFunc func(*http.Request) (*url.URL, error)

	if proxyStr != "" {
		// If custom variable exists, parse and use it exclusively
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			return nil, fmt.Errorf("invalid %s: %w", constants.MCP_PROXY_URL, err)
		}
		proxyFunc = http.ProxyURL(proxyURL)
	} else {
		// 2. If not present, fall back to system environment variables (HTTP_PROXY, HTTPS_PROXY, NO_PROXY)
		proxyFunc = http.ProxyFromEnvironment
	}

	transport := &http.Transport{
		Proxy:                 proxyFunc,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30,
	}

	return NewSearch(ctx, &Config{
		Region:     RegionWT,
		MaxResults: 10,
		HTTPClient: httpClient,
	})
}
