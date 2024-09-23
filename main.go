package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	inputURL := os.Args[1]
	fmt.Printf("starting crawl of: %s...\n", inputURL)

	body, err := parseHTML(inputURL)
	if err != nil {
		fmt.Printf("error parsing HTML: %v\n", err)
		os.Exit(1)
	}

	urls, err := getURLsFromHTML(string(body), inputURL)
	if err != nil {
		fmt.Printf("error getting URLs from body: %v\n", err)
		os.Exit(1)
	}

	for _, url := range urls {
		fmt.Println(url)
	}
}
