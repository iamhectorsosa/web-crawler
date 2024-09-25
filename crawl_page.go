package main

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/charmbracelet/log"
)

type config struct {
	maxPages           int
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func crawl(rawBaseURL string, concurrency, maxPages int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	config := &config{
		maxPages:           maxPages,
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, concurrency),
		wg:                 &sync.WaitGroup{},
	}

	logger := log.WithPrefix("STARTING CRAWLER")
	logger.Print("parameters", "baseUrl", rawBaseURL, "concurrency", concurrency, "maxPages", maxPages)

	config.wg.Add(1)
	go config.crawlPage(rawBaseURL)
	config.wg.Wait()

	return config, nil
}

func (cfg *config) acquire() {
	cfg.concurrencyControl <- struct{}{}
}

func (cfg *config) release() {
	<-cfg.concurrencyControl
	cfg.wg.Done()
}

func (cfg *config) checkMaxPages() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages) < cfg.maxPages
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

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.acquire()
	defer cfg.release()

	if ok := cfg.checkMaxPages(); !ok {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Error("failed to parse URL", "url", rawCurrentURL, "err", err)
		return
	}

	if cfg.baseURL.Host != currentURL.Host {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Error("failed to normalize URL", "url", rawCurrentURL, "err", err)
		return
	}

	if isFirst := cfg.addPageVisit(normalizedURL); !isFirst {
		return
	}

	log.Info("crawling", "url", normalizedURL, "limit", fmt.Sprintf("%d/%d", len(cfg.pages), cfg.maxPages))
	htmlBody, err := parseHTML(rawCurrentURL)
	if err != nil {
		log.Error("failed to crawl", "url", normalizedURL, "err", err)
		return
	}

	urls, err := getURLsFromHTML(htmlBody, cfg.baseURL.String())
	if err != nil {
		log.Error("failed to collect internal URLs", "url", normalizedURL, "err", err)
		return
	}

	for _, url := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}
