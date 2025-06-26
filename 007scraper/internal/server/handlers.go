package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/samnart1/golang/007scrapper/internal/config"
	"github.com/samnart1/golang/007scrapper/internal/scraper"
	"github.com/samnart1/golang/007scrapper/pkg/types"
)

type Server struct {
	scraper *scraper.Scraper
	config  *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{
		scraper: scraper.New(cfg),
		config: cfg,
	}
}

type ScrapeRequest struct {
	URL				string				`json:"url"`
	IncludeLinks	bool				`json:"include_links,omitempty"`
	IncludeImages	bool				`json:"include_images,omitempty"`
	Timeout			string				`json:"timeout,omitempty"`
	UserAgent		string				`json:"user_agent,omitempty"`
	CustomHeaders	map[string]string	`json:"custom_headers,omitempty"`
}

type BatchScrapedResult struct {
	URLs 			[]string			`json:"urls"`
	IncludeLinks	bool				`json:"include_links,omitempty"`
	IncludeImages	bool				`json:"include_images,omitempty"`
	Timeout			string				`json:"timeout,omitempty"`
	UserAgent		string				`json:"user_agent,omitempty"`
	CustomHeaders	map[string]string	`json:"custom_headers,omitempty"`
	Concurrent		int					`json:"concurrent,omitempty"`
	RateLimit		string				`json:"rate_limit,omitempty"`
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status": 		"health",
		"timestamp":	time.Now().UTC(),
		"version":		"1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ScrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	options := types.DefaultScrapeOptions()
	options.IncludeLinks = req.IncludeLinks
	options.IncludeImages = req.IncludeImages

	if req.Timeout != "" {
		if timeout, err := time.ParseDuration(req.Timeout); err == nil {
			options.Timeout = timeout
		}
	}

	if req.UserAgent != "" {
		options.UserAgent = req.UserAgent
	}

	if req.CustomHeaders != nil {
		options.CustomeHeaders = req.CustomHeaders
	}

	result := s.scraper.ScrapeURL(req.URL, options)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) BatchScrapeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BatchScrapedResult
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if len(req.URLs) == 0 {
		http.Error(w, "URLs are required", http.StatusBadRequest)
		return
	}

	if len(req.URLs) > 100 {
		http.Error(w, "Maximum 100 URLs allowed per batch", http.StatusBadRequest)
		return
	}

	options := types.DefaultScrapeOptions()
	options.IncludeLinks = req.IncludeLinks
	options.IncludeImages = req.IncludeImages

	if req.Timeout != "" {
		if timeout, err := time.ParseDuration(req.Timeout); err == nil {
			options.Timeout = timeout
		}
	}

	if req.UserAgent != "" {
		options.UserAgent = req.UserAgent
	}

	if req.CustomHeaders != nil {
		options.CustomeHeaders = req.CustomHeaders
	}

	workers := s.config.MaxConcurrent
	if req.Concurrent > 0 && req.Concurrent <= 20 {
		workers = req.Concurrent
	}

	delay := s.config.RateLimit
	if req.RateLimit != "" {
		if d, err := time.ParseDuration(req.RateLimit); err == nil {
			delay = d
		}
	}

	batchResult := s.runBatchScraping(req.URLs, options, workers, delay)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batchResult)
} 

func (s *Server) runBatchScraping(urls []string, options types.ScrapeOptions, workers int, delay time.Duration) types.BatchScrapedResult {
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
		go func() {
			defer wg.Done()
			for url := range urlChan {
				result := s.scraper.ScrapeURL(url, options)
				resultChan <- result

				if delay > 0 {
					time.Sleep(delay)
				}
			}
		}()
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