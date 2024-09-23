package main

import (
	"fmt"
	"io"
	"net/http"
)

func parseHTML(inputURL string) (string, error) {
	resp, err := http.Get(inputURL)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "text/html" {
		return "", fmt.Errorf("invalid content type: %s", resp.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading from body: %v", err)
	}

	return string(body), nil
}
