package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const cacheAPIRequestTimeout = 30 * time.Second

// cacheKeyExists reports whether a Bitrise key-value cache entry exists for
// key, without downloading the cached archive.
//
// It calls the Advanced Build Cache Service's restore-lookup endpoint, the
// same one restore-cache steps use to resolve a cache key to a signed
// download URL. A 200 response means an entry exists; the archive itself is
// never fetched, so this is a pure existence check.
func cacheKeyExists(apiBaseURL, accessToken, key string) (bool, error) {
	reqURL := fmt.Sprintf("%s/restore?cache_keys=%s", apiBaseURL, url.QueryEscape(key))

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return false, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{Timeout: cacheAPIRequestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("call cache service: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected response from cache service (HTTP %d): %s", resp.StatusCode, body)
	}
}
