package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Router() *mux.Router {
	r := mux.NewRouter()

	r.Use(corsMiddleware)

	r.Use(loggingMiddleware)

	r.HandleFunc("/health", s.HealthHandler).Methods("GET")

	r.HandleFunc("/scrape", s.ScrapeHandler).Methods("POST")
	r.HandleFunc("/scrape/batch", s.BatchScrapeHandler).Methods("POST")

	r.HandleFunc("/", s.IndexHandler).Methods("GET")

	return r
}

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	
		<!DOCTYPE html>
		<html>
			<head>
				<title>Web Scraper API</title>
				<style>
					body { font-family: Arial, sans-serif; margin: 40px; }
					.endpoint { background: #f5f5f5; padding: 10px; margin: 10px; border-radius: 5px; }
					code { background: #eee; padding: 2px 4px; border-radius: 3px }
				</style>
			</head>
			<body>
				<h1>Web Scraper API</h1>
				<p>Welcome to the Web Scraper API. Available enpoints:</p>

				<div class="endpoint">
					<h3>GET /health</h3>
					<p>Health check endpoint</p>
				</div>
					
				<div class="endpoint">
					<h3>POST /scrape</h3>
					<p>Scrape a single URL</p>
					<p>Example payload:</p>
					<pre>
						<code>
							{
								"url": "https://example.com",
								"include_links": true,
								"include_images": true
							}
						</code>
					</pre>
				</div>

				<div class="endpoint">
					<h3>POST /scrape/batch</h3>
					<p>Scrape multiple URLs</p>
					<p>Example payload:</p>
					<pre>
						<code>
							{
								"urls": ["https://example.com", "https://google.com"],
								"include_links": true,
								"include_images": true,
								"concurrent": 3
							}
						</code>
					</pre>
				</div>

				<h2>Test the API</h2>
				<p>You can test the API using curl:</p>
				<pre>
					<code>
						curl -X POST http://localhost:8080/scrape \
						-H "Content-Type: application/json"	\
						-d '{"url": "https://example.com"}'
					</code>
				</pre>

			</body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}