package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("couldn't parse URL: %s, %v\n", baseURL, err)
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("couldn't parse URL: %s, %v\n", currentURL, err)
		return
	}

	if baseURL.Host != currentURL.Host {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("couldn't normalize URL: %s, %v\n", normalizedURL, err)
		return
	}

	if _, visited := pages[normalizedURL]; visited {
		pages[normalizedURL]++
		return
	}

	pages[normalizedURL] = 1

	fmt.Printf("crawling: %s\n", normalizedURL)
	htmlBody, err := parseHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("couldn't parse HTML: %s, %v\n", normalizedURL, err)
		return
	}

	urls, err := getURLsFromHTML(string(htmlBody), baseURL.String())
	if err != nil {
		fmt.Printf("couldn't get URLs from URL: %s, %s\n", normalizedURL, err)
		return
	}

	for _, url := range urls {
		crawlPage(rawBaseURL, url, pages)
	}
}
