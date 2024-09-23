package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	inputURL := args[0]
	fmt.Printf("starting crawl of: %s...\n", inputURL)

	// resp, err := http.Get(inputURL)
	// if err != nil {
	// 	fmt.Printf("error fetching URL: %v\n", err)
	// 	os.Exit(1)
	// }

	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Printf("error fetching URL. Status code: %d\n", resp.StatusCode)
	// 	os.Exit(1)
	// }

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Printf("error reading body: %v\n", err)
	// 	os.Exit(1)
	// }

	// normalizedURL, err := normalizeURL(inputURL)
	// if err != nil {
	// 	fmt.Printf("error normalizing URL: %v\n", err)
	// 	os.Exit(1)
	// }

	// urls, err := getURLsFromHTML(string(body), inputURL)
	// if err != nil {
	// 	fmt.Printf("error getting URLs from body: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println("Printing URLs:")
	// for _, url := range urls {
	// 	fmt.Println(url)
	// }
}
