package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	url := parsedURL.Host + parsedURL.Path
	url = strings.ToLower(url)
	url = strings.TrimSuffix(url, "/")

	return url, nil
}
