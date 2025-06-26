package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/samnart1/golang/007scrapper/internal/scraper"
	"github.com/samnart1/golang/007scrapper/pkg/types"
	"github.com/spf13/cobra"
)

var (
	outputFormat	string
	includeLinks	bool
	includeImages	bool
	customTimeout	time.Duration
	customAgent		string
)

var scrapeCmd = &cobra.Command{
	Use: "scrape [URL]",
	Short: "scrape a single URL",
	Long: `scrape a single URL and extract information like title, description, links and images
		
		Examples:
			webscraper scrape https://example.com
			webscraper scrape https://example.com --format json
			webscraper scrape https://example.com --no-links --no-images
			webscraper scrape https://example.com --timeout 60s --agent "Custom Bot"`,
	
	Args: cobra.ExactArgs(1),
	RunE: runScrape,
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().StringVarP(&outputFormat, "format", "f", "pretty", "Output format: pretty, json")
	scrapeCmd.Flags().BoolVar(&includeLinks, "links", true, "Include links in output")
	scrapeCmd.Flags().BoolVar(&includeImages, "images", true, "Include images in output")
	scrapeCmd.Flags().DurationVar(&customTimeout, "timeout", 0, "Custom timeout (e.g., 30s, 1m)")
	scrapeCmd.Flags().StringVar(&customAgent, "agent", "", "Custom user agent")
}

func runScrape(cmd *cobra.Command, args []string) error {
	url := args[0]

	if verbose {
		fmt.Printf("Scraping URL: %s\n", url)
	}

	s := scraper.New(cfg)

	options := types.DefaultScrapeOptions()
	options.IncludeLinks = includeLinks
	options.IncludeImages = includeImages

	if customTimeout > 0 {
		options.Timeout = customTimeout
	}

	if customAgent != "" {
		options.UserAgent = customAgent
	}

	result := s.ScrapeURL(url, options)

	switch outputFormat {
	case "json":
		return outputJSON(result)
	case "pretty":
		return outputPretty(result)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

func outputJSON(result types.ScrapeResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	return encoder.Encode(result)
}

func outputPretty(result types.ScrapeResult) error {
	fmt.Printf("URL: %s\n", result.URL)
	fmt.Printf("Status: ")
	if result.Success {
		fmt.Printf("✓ Succdss (HTTP %d)\n", result.StatusCode)
	} else {
		fmt.Printf("✗ Failed (HTTP %d)\n", result.StatusCode)
		if result.Error != "" {
			fmt.Printf("Error: %s\n", result.Error)
		}
		return nil
	}

	fmt.Printf("Duration: %v\n", result.Duration)
	fmt.Printf("Scraped: %s\n\n", result.ScrapedAt.Format(time.RFC3339))

	fmt.Printf("Title: %s\n", result.Title)

	if result.Description != "" {
		fmt.Printf("Description: %s\n", result.Description)
	}

	if len(result.Keywords) > 0 {
		fmt.Printf("Keywords: %v\n", result.Keywords)
	}

	if len(result.Headers) > 0 {
		fmt.Printf("\nHeaders:\n")
		for tag, text := range result.Headers {
			fmt.Printf("	%s: %s\n", tag, text)
		}
	}

	if len(result.Links) > 0 {
		fmt.Printf("\nLnks (%d):\n", len(result.Links))
		for i, link := range result.Links {
			if i >= 10 {
				fmt.Printf("	... and %d more links\n", len(result.Links)-10)
				break
			}
			fmt.Printf("	%s", link.URL)
			if link.Text != "" {
				fmt.Printf("	(%s)", link.Text)
			}
			fmt.Printf("\n")
		}
	}

	if len(result.Images) > 0 {
		fmt.Printf("\nImages (%d):\n", len(result.Images))	
		for i, img := range result.Images {
			if i >= 10 {
				fmt.Printf("	... and %d more images\n", len(result.Images) - 10)
				break
			}
			fmt.Printf("	%s", img.URL)
			if img.Alt != "" {
				fmt.Printf("	(alt: %s)", img.Alt)
			}
			fmt.Printf("\n")
		}
	}

	return nil
}