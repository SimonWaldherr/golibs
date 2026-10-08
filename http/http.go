// Package http provides helpers for HTTP requests with a shared transport.
package http

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Transporter is shared by clients created with Client. Configure it before
// making requests; transports reuse connections and support concurrent use.
var Transporter = &http.Transport{
	MaxIdleConns:    1024,
	MaxConnsPerHost: 1024,
	IdleConnTimeout: 10 * time.Second,
}

// NewRequest is the standard library request constructor.
var NewRequest = http.NewRequest

// Client returns a client with the shared transport and the given timeout.
func Client(timeout time.Duration) http.Client {
	return http.Client{
		Transport: Transporter,
		Timeout:   timeout,
	}
}

// GetString retrieves a response body with a 15-second timeout. HTTP error
// statuses are returned as response bodies, matching net/http.Client.Get.
func GetString(url string) (string, error) {
	return GetStringContext(context.Background(), url)
}

// GetStringContext retrieves a response body with the given context and a
// 15-second timeout. The context can cancel both the request and body reading.
func GetStringContext(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	client := Client(15 * time.Second)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
