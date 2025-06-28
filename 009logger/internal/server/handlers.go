package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/samnart1/golang/009l9gger/internal/logger"
)

type LogRequest struct {
	Message string					`json:"message"`
	Level	string					`json:"level"`
	Source	string					`json:"source,omitempty"`
	Data	map[string]interface{}	`json:"data,omitempty"`
}

type LogResponse struct {
	Success		bool	`json:"success"`
	Timestamp	string	`json:"timestamp"`
	Error		string	`json:"error,omitempty"`
}

type HealthResponse struct {
	Status		string	`json:"status"`
	Timestamp	string	`json:"timestamp"`
	Version		string	`json:"version"`
}

func (s *Server) handleLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendErrorResponse(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		s.sendErrorResponse(w, "Message is required", http.StatusBadRequest)
		return
	}

	if req.Level == "" {
		req.Level = "info"
	}

	entry := logger.Entry{
		Message: req.Message,
		Level: req.Level,
		Source: req.Source,
		Timestamp: time.Now(),
		Data: req.Data,
	}

	if err := s.logger.Log(entry); err != nil {
		s.sendErrorResponse(w, fmt.Sprintf("failed to log message: %v", err), http.StatusInternalServerError)
		return
	}

	response := LogResponse{
		Success: true,
		Timestamp: entry.Timestamp.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status: "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Version: "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Go Simple Logger</title>
			<style>
				body { font-fammily: Arial, sans-serif; margin: 40px; }
				.container { max-width: 800px; margin: 0 auto; }
				.endpoint { background: #f5f5f5; padding: 10px, margin: 10px 0; border-radius: 5px; }
				.method { font-weight: bold; color: #2196F3; }
				.url { font-family: monospace; background: #fff; padding: 2px 5px; border-radius: 3px }
				.example { background: #f0f0f0; padding: 10px; border-radius: 5px; overflow-x: auto }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Go Simple Logger API</h1>
				<p>A simple and powerful logging utility with HTTP API</p>

				<h2>Available Endpoints</h2>

				<div class="endpoint">
					<span class="method">POST</span> <span class="url">/log</span>
					<p>Log a message with optional level and metadata</p>
				</div>

				<div class="endpoint">
					<span class="method">GET</span> <span class="url">/health</span>
					<p>Health check endpoint</p>
				</div>

				<h2>Usage Examples</h2>

				<div class="example">
					<h3>Basic Log Message</h3>
					<pre>curl -X POST http://localhost:8080/log \
						-H "Content-Type: application/json" \
						-d '{"message": "Hello, World!", "level": "info"}'
					</pre>
				</div>

				<div class="example">
					<h3>Log with Source and Data</h3>
					<pre>curl -X POST http://localhost:8080/log \
						-H "Content-Type: application/json" \
						-d '{
								"message": "User login successful",
								"level": "info",
								"source": "auth-service",
								"data": {"user_id": 123, "ip": "192.168.1.1"}
							}'
					</pre>
				</div>

				<div class="example">
					<h3>Health check</h3>
					<pre>curl http://localhost:8080/health</pre>
				</div>
			</div>
		<body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (s *Server) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := LogResponse{
		Success: false,
		Timestamp: time.Now().Format(time.RFC3339),
		Error: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}