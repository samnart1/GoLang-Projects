package types

import "time"

type ScrapeResult struct {
	URL			string				`json:"url"`
	Title		string				`json:"title"`
	Description	string				`json:"description,omitempty"`
	Keywords	[]string			`json:"keywords,omitempty"`
	Links		[]Link				`json:"links,omitempty"`
	Images		[]Image				`json:"images,omitempty"`
	Headers		map[string]string	`json:"headers,omitempty"`
	StatusCode	int					`json:"status_code"`
	Success		bool				`json:"success"`
	Error		string				`json:"error,omitempty"`
	ScrapedAt	time.Time			`json:"scraped_at"`
	Duration	time.Duration		`json:"duration"`
}

type Link struct {
	URL	 string	`json:"url"`
	Text string	`json:"text"`
	Rel	 string	`json:"rel,omitempty"`
}

type Image struct {
	URL	string	`json:"url"`
	Alt	string	`json:"alt,omitempty"`
	Src	string	`json:"src"`
}

type ScrapeOptions struct {
	Timeout			time.Duration		`json:"timeout"`
	UserAgent		string				`json:"user_agent"`
	FollowLinks		bool				`json:"follow_links"`
	MaxDepth		int					`json:"max_depth"`
	IncludeImages	bool				`json:"include_images"`
	IncludeLinks	bool				`json:"include_links"`
	CustomeHeaders	map[string]string	`json:"custom_headers,omitempty"`
}

type BatchScrapedResult struct {
	Results		[]ScrapeResult	`json:"results"`
	Total		int				`json:"total"`
	Success		int			`json:"success"`
	Failed		int			`json:"failed"`
	StartTime	time.Time		`json:"start_time"`
	EndTime		time.Time		`json:"end_time"`
	Duration	time.Duration	`json:"duration"`
}

func DefaultScrapeOptions() ScrapeOptions {
	return ScrapeOptions{
		Timeout: 30 * time.Second,
		UserAgent: "Go-web-scraper/1.0",
		FollowLinks: false,
		MaxDepth: 1,
		IncludeImages: true,
		IncludeLinks: true,
		CustomeHeaders: make(map[string]string),
	}
}