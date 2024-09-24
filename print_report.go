package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
)

type Page struct {
	URL   string
	Count int
}

func printReport(pages map[string]int, baseURL string) {
	logStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("246")).Padding(0, 1)

	var title strings.Builder
	fmt.Fprintf(&title, "\nREPORT for %s", baseURL)
	fmt.Println(logStyle.Render(title.String()))

	sortedPages := sortPages(pages)

	columns := []string{"URL", "Count"}
	rows := make([][]string, 0, len(sortedPages))

	for _, page := range sortedPages {
		rows = append(rows, []string{
			page.URL,
			strconv.Itoa(page.Count),
		})
	}

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers(columns...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row%2 == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("246")).MarginRight(6)
			}
			return lipgloss.NewStyle()
		})

	fmt.Print(t)

	err := generateCSVReport(sortedPages, "report.csv")
	if err != nil {
		log.Error("failed to generate CSV report", "err", err)
	}

	fmt.Println("")

	log.Info("CSV Report generated", "url", baseURL, "filepath", "report.csv")
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

	return nil
}
