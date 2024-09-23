package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("too few arguments provided")
		os.Exit(1)
	}

	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	inputURL := os.Args[1]
	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error - maxConcurrency: %v", err)
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Error - maxPages: %v", err)
		os.Exit(1)
	}

	cfg, err := newCrawler(inputURL, maxConcurrency, maxPages)

	if err != nil {
		fmt.Printf("Error - configure: %v", err)
		os.Exit(1)
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(inputURL)
	cfg.wg.Wait()

	for url, count := range cfg.pages {
		fmt.Printf("%s - %d\n", url, count)
	}
}
