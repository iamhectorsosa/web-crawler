package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	inputURL := os.Args[1]
	fmt.Printf("starting crawl of: %s...\n", inputURL)

	pages := make(map[string]int)

	crawlPage(inputURL, inputURL, pages)

	for url, count := range pages {
		fmt.Printf("%s - %d\n", url, count)
	}
}
