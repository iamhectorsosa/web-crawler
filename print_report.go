package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Page struct {
	URL   string
	Count int
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Printf(`
=============================
  REPORT for %s
=============================
`, baseURL)

	sortedPages := sortPages(pages)
	for _, page := range sortedPages {
		url := page.URL
		count := page.Count
		fmt.Printf("Found %d internal links to %s\n", count, url)
	}

	err := generateCSVReport(sortedPages, "report.csv")
	if err != nil {
		fmt.Println("Error generating CSV report:", err)
	}
}

func sortPages(pages map[string]int) []Page {
	pagesSlice := []Page{}
	for url, count := range pages {
		pagesSlice = append(pagesSlice, Page{URL: url, Count: count})
	}
	sort.Slice(pagesSlice, func(i, j int) bool {
		if pagesSlice[i].Count == pagesSlice[j].Count {
			return pagesSlice[i].URL < pagesSlice[j].URL
		}
		return pagesSlice[i].Count > pagesSlice[j].Count
	})
	return pagesSlice
}

func generateCSVReport(pages []Page, filename string) error {
	// Create a new CSV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header to the CSV file
	header := []string{"URL", "Count"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write the page data to the CSV file
	for _, page := range pages {
		row := []string{page.URL, strconv.Itoa(page.Count)}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	fmt.Printf("CSV report generated successfully: %s\n", filename)
	return nil
}
