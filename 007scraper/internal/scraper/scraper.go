package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/samnart1/golang/007scrapper/internal/config"
	"github.com/samnart1/golang/007scrapper/pkg/types"
)

type Scraper struct {
	client *http.Client
	config *config.Config
}

func New(cfg *config.Config) *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		config: cfg,
	}
}

func (s *Scraper) ScrapeURL(url string, options types.ScrapeOptions) types.ScrapeResult {
	startTime := time.Now()
	result := types.ScrapeResult{
		URL: url,
		ScrapedAt: startTime,
		Success: false,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		result.Duration = time.Since(startTime)
		return result
	}

	req.Header.Set("User-Agent", options.UserAgent)

	for key, value := range options.CustomeHeaders {
		req.Header.Set(key, value)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		result.Error = fmt.Sprintf("failed to fetch URL: %v", err)
		result.Duration = time.Since(startTime)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		result.Error = fmt.Sprintf("HTTP error: %d %s", resp.StatusCode, resp.Status)
		result.Duration = time.Since(startTime)
		return result
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("failed to parse HTML: %v", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.Title = strings.TrimSpace(doc.Find("title").Text())
	if result.Title == "" {
		result.Title = "No title found!"
	}

	if desc, exsists := doc.Find("meta[name='description']").Attr("content"); exsists {
		result.Description = strings.TrimSpace(desc)
	}

	if keywords, exists := doc.Find("meta[name='keywords']").Attr("content"); exists {
		keywordList := strings.Split(keywords, ",")
		for i, keyword := range keywordList {
			keywordList[i] = strings.TrimSpace(keyword)
		}
		result.Keywords = keywordList
	}

	if options.IncludeLinks {
		result.Links = s.extractLinks(doc, url)
	}

	if options.IncludeImages {
		result.Images = s.extractImages(doc, url)
	}

	result.Headers = s.extractHeaders(doc)

	result.Success = true
	result.Duration = time.Since(startTime)
	return result
}

func (s *Scraper) extractLinks(doc *goquery.Document, baseURL string) []types.Link {
	var links []types.Link

	doc.Find("a[href]").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		absoluteURL := s.resolveURL(href, baseURL)

		link := types.Link{
			URL: absoluteURL,
			Text: strings.TrimSpace(sel.Text()),
		}

		if rel, exists := sel.Attr("rel"); exists {
			link.Rel = rel
		}

		links = append(links, link)
	})

	return links
}

func (s *Scraper) extractImages(doc *goquery.Document, baseURL string) []types.Image {
	var images []types.Image

	doc.Find("img").Each(func(i int, sel *goquery.Selection) {
		src, exists := sel.Attr("src")
		if !exists {
			return
		}

		absoluteURL := s.resolveURL(src, baseURL)

		image := types.Image{
			URL: absoluteURL,
			Src: src,
		}

		if alt, exists := sel.Attr("alt"); exists {
			image.Alt = alt
		}

		images = append(images, image)
	})

	return images
}

func (s *Scraper) extractHeaders(doc *goquery.Document) map[string]string {
	headers := make(map[string]string)

	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		if text := strings.TrimSpace(doc.Find(tag).First().Text()); text != "" {
			headers[tag] = text
		}
	}

	return headers
}

func (s *Scraper) resolveURL(href, baseURL string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	if strings.HasPrefix(href, "//") {
		return "https:" + href
	}

	if strings.HasPrefix(href, "/") {
		if idx := strings.Index(baseURL[8:], "/"); idx != -1 {
			return baseURL[:8+idx] + href
		}
		return baseURL + href
	}

	lastSlash := strings.LastIndex(baseURL, "/")
	if lastSlash > 7 {
		return baseURL[:lastSlash+1] + href
	}

	return baseURL + "/" + href
}