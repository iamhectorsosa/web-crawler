package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mygo [command]",
	Short: "Mygo Crawl is a Go-based Web Crawler tool for HTML analysis.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var mvp = &cobra.Command{
	Use:   "mvp [url]",
	Short: "Get the most linked-to URLs (MVPs) from a website.",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("Need to specify website's base URL: %v", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL := args[0]
		concurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			return err
		}
		maxPages, err := cmd.Flags().GetInt("max-pages")
		if err != nil {
			return err
		}

		cfg, err := crawl(baseURL, concurrency, maxPages)

		if err != nil {
			return err
		}

		printReport(cfg.pages, baseURL)
		return nil
	},
}

func init() {
	mvp.Flags().IntP("concurrency", "c", 12, "concurrency capacity for analysis")
	mvp.Flags().IntP("max-pages", "m", 120, "maximum number of pages for analysis")
	rootCmd.AddCommand(mvp)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
