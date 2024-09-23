package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pageCount          int
	maxPages           int
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func newCrawler(rawBaseURL string, maxPages, maxConcurrency int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &config{
		pageCount:          0,
		maxPages:           maxPages,
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}, nil
}

func (cfg *config) addPageCount() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.pageCount++
}

func (cfg *config) addPageVisit(url string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[url]; visited {
		cfg.pages[url]++
		return false
	}

	cfg.pages[url] = 1
	return true
}

func (cfg *config) acquire() {
	cfg.concurrencyControl <- struct{}{}
}

func (cfg *config) release() {
	<-cfg.concurrencyControl
	cfg.wg.Done()
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.acquire()
	defer cfg.release()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("couldn't parse URL: %s, %v\n", currentURL, err)
		return
	}

	if cfg.baseURL.Host != currentURL.Host {
		return
	}

	if cfg.pageCount >= cfg.maxPages {
		return
	}

	cfg.addPageCount()

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("couldn't normalize URL: %s, %v\n", normalizedURL, err)
		return
	}

	if isFirst := cfg.addPageVisit(normalizedURL); !isFirst {
		return
	}

	fmt.Printf("crawling: %s\n", normalizedURL)
	htmlBody, err := parseHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("couldn't parse HTML: %s, %v\n", normalizedURL, err)
		return
	}

	urls, err := getURLsFromHTML(htmlBody, cfg.baseURL.String())
	if err != nil {
		fmt.Printf("couldn't get URLs from HTML at: %s, %s\n", normalizedURL, err)
		return
	}

	for _, url := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}
