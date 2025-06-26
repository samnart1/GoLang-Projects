package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/samnart1/golang/007scrapper/internal/scraper"
	"github.com/samnart1/golang/007scrapper/pkg/types"
	"github.com/spf13/cobra"
)

var (
	inputFile 	string
	outputFile	string
	concurrent 	int
	rateLimit	time.Duration
)

var batchCmd = &cobra.Command{
	Use: "batch",
	Short: "scrape multiple URLs from a file",
	Long: `scrape multiple URLs from a file with concurrent processing
	
		The input file should contain one URL per line. Empty lines and lines starting with # are ignored.
		
		Examples:
			webscraper batch --input urls.txt
			webscraper batch --input urls.txt --output results.json
			webscraper batch --input urls.txt --concurrent 5 --rate-limit 200ms`,
		
	RunE: runBatch,
}

func init() {
	rootCmd.AddCommand(batchCmd)

	batchCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file containing URLs (required)")
	batchCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for results (default: stdout)")
	batchCmd.Flags().IntVarP(&concurrent, "concurrent", "c", 0, "Number of concurrent workers (default: from config)")
	batchCmd.Flags().DurationVar(&rateLimit, "rate-limit", 0, "Rate limit between requests (default: from config)")
	batchCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format: json, pretty")
	batchCmd.Flags().BoolVar(&includeLinks, "links", true, "Include links in output")
	batchCmd.Flags().BoolVar(&includeImages, "images", true, "Include images in output")

	batchCmd.MarkFlagRequired("input")
}

func runBatch(cmd *cobra.Command, args []string) error {
	urls, err := readURLsFromFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read URLs from file: %v", err)
	}

	if len(urls) == 0 {
		return fmt.Errorf("no valid URLs found in input files")
	}

	if verbose {
		fmt.Printf("Found %d URLs to scrape\n", len(urls))
	}

	workers := cfg.MaxConcurrent
	if concurrent > 0 {
		workers = concurrent
	}

	delay := cfg.RateLimit 
	if rateLimit > 0 {
		delay = rateLimit
	}

	s := scraper.New(cfg)

	options := types.DefaultScrapeOptions()
	options.IncludeLinks = includeLinks
	options.IncludeImages = includeImages

	batchResult := runBatchScraping(s, urls, options, workers, delay)

	if outputFile != "" {
		return writeResultToFile(batchResult, outputFile)
	}

	return outputBatchResults(batchResult)
}

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			urls = append(urls, line)
		} else if verbose {
			fmt.Printf("Skipping invalid URL: %s\n", line)
		}
	}

	return urls, scanner.Err()
}

func runBatchScraping(s *scraper.Scraper, urls []string, options types.ScrapeOptions, workers int, delay time.Duration) types.BatchScrapedResult {
	startTime := time.Now()

	urlChan := make(chan string, len(urls))
	resultChan := make(chan types.ScrapeResult, len(urls))

	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(s, urlChan, resultChan, options, delay, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []types.ScrapeResult
	var successCount, failedCount int

	for result := range resultChan {
		results = append(results, result)
		if result.Success {
			successCount++
		} else {
			failedCount++
		}

		if verbose {
			status := "✔"
			if !result.Success {
				status = "✗"
			}
			fmt.Printf("%s %s (%v)\n", status, result.URL, result.Duration)
		}
	}

	endTime := time.Now()

	return types.BatchScrapedResult{
		Results: results,
		Total: len(urls),
		Success: successCount,
		Failed: failedCount,
		StartTime: startTime,
		EndTime: endTime,
		Duration: endTime.Sub(startTime),
	}
}

func worker(s *scraper.Scraper, urlChan <-chan string, resultChan chan<- types.ScrapeResult, options types.ScrapeOptions, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range urlChan {
		result := s.ScrapeURL(url, options)
		resultChan <- result

		if delay > 0 {
			time.Sleep(delay)
		}
	}
}

func writeResultToFile(batchResult types.BatchScrapedResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(batchResult)
}

func outputBatchResults(batchResult types.BatchScrapedResult) error {
	switch outputFormat {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", " ")
		return encoder.Encode(batchResult)
	case "pretty":
		fmt.Printf("Batch Scraping Results\n")
		fmt.Printf("======================\n")
		fmt.Printf("Total URLs: %d\n", batchResult.Total)
		fmt.Printf("Successful: %d\n", batchResult.Success)
		fmt.Printf("Failed: %d\n", batchResult.Failed)
		fmt.Printf("Duration: %v\n", batchResult.Duration)
		fmt.Printf("Average: %v per URL\n\n", batchResult.Duration/time.Duration(batchResult.Total))

		for _, result := range batchResult.Results {
			status := "🗸"
			if !result.Success {
				status = "✗"
			}
			fmt.Printf("%s %s - %s (%v)\n", status, result.URL, result.Title, result.Duration)
		}
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}