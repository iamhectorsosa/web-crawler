package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Helper function to normalize URL by stripping "www" prefix
func normalizeHost(u *url.URL) string {
	host := u.Hostname()
	if strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(host, "www.")
	}
	return host
}

func parseHTML(inputURL string) (string, error) {
	// Parse baseURL to compare later
	base, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %v", err)
	}

	// Initialize a custom HTTP client with redirect checking
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Parse the redirected URL
			redirectedURL := req.URL
			if normalizeHost(base) != normalizeHost(redirectedURL) {
				return fmt.Errorf("redirect to external URL: %s", redirectedURL.String())
			}
			// Allow the redirect to continue
			return nil
		},
	}

	// Make the GET request
	resp, err := client.Get(inputURL)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	// Only accept status 200
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Only accept text/html
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading from body: %v", err)
	}

	return string(body), nil
}
