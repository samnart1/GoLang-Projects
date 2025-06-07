package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"github.com/samnart1/GoLang-Projects/005server/pkg/response"
)

type PageData struct {
	Title		string
	Message		string
	CurrentTime	string
	Version		string
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		s.handleNotFound(w, r)
		return
	}

	data := PageData{
		Title: 			"Go HTTP Server",
		Message: 		"Welcome to the Go HTTP Server!",
		CurrentTime: 	time.Now().Format(time.RFC3339),
		Version: 		"1.0.0",
	}

	s.renderTemplate(w, "index.html", data)
}

func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Go"
	}

	message := fmt.Sprintf("Hello, %s!", name)

	if r.Header.Get("Accept") == "application/json" {
		response.JSON(w, http.StatusOK, map[string]string{
			"message": 	message,
			"time":		time.Now().Format(time.RFC3339),
		})
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, message)
}

func (s *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:	"About - Go Http Server",
		Message: "This is a simple http server built with Go to demonstrate web development concepts",
		Version: "1.0.0",
	}

	s.renderTemplate(w, "about.html", data)
}

func (s *Server) handleAPIInfo(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"server": map[string]string{
			"name":			"Go HTTP Server",
			"version":		"1.0.0",
			"environment":	s.config.Environment,
			"go_version":	runtime.Version(),
		},
		"request": map[string]string{
			"method":		r.Method,
			"url":			r.URL.String(),
			"user_agent": 	r.UserAgent(),
			"remote_ip":	r.RemoteAddr,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	response.JSON(w, http.StatusOK, info)
}

func (s *Server) handleAPIEcho(w http.ResponseWriter, r *http.Request) {
	var body interface{}
	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&body); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}
	}

	// collect query paramms
	queryParams := make(map[string][]string)
	for key, values := range r.URL.Query() {
		queryParams[key] = values
	}

	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = values[0]
	}

	echo := map[string]interface{}{
		"method":		r.Method,
		"url":			r.URL.String(),
		"headers":		headers,
		"query":		queryParams,
		"remote_addr":	r.RemoteAddr,
		"timestamp":	time.Now().Format(time.RFC3339),
	}

	if body != nil {
		echo["body"] = body
	}

	response.JSON(w, http.StatusOK, echo)
}

func (s *Server) handleAPITime(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	timeInfo := map[string]interface{}{
		"timestamp": 	now.Format(time.RFC3339),
		"unix":			now.Unix(),
		"timezone":		now.Location().String(),
		"formats":		map[string]string{
			"rfc3339":	now.Format(time.RFC3339),
			"iso8601":	now.Format("2006-01-02T15:04:05Z07:00"),
			"human":	now.Format("January 2, 2006 at 3:04 PM MST"),
			"date":		now.Format("2006-01-02"),
			"time":		now.Format("15:04:05"),
		},
	}

	response.JSON(w, http.StatusOK, timeInfo)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":		"healthy",
		"timestamp":	time.Now().Format(time.RFC3339),
		"uptime":		time.Since(time.Now()).String(),
		"version":		"1.0.0",
	}

	response.JSON(w, http.StatusOK, health)
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	ready := map[string]interface{}{
		"ready":	true,
		"timestamp":	time.Now().Format(time.RFC3339),
	}

	response.JSON(w, http.StatusOK, ready)
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") == "application/json" {
		response.Error(w, http.StatusNotFound, "Endpoint Not Found!")
		return
	}

	data := PageData{
		Title: "404 - Page Not Found",
		Message: fmt.Sprintf("The page '%s' was not found!", r.URL.Path),
	}

	w.WriteHeader(http.StatusNotFound)
	s.renderTemplate(w, "error.html", data)
}

func (s *Server) renderTemplate(w http.ResponseWriter, templateName string, data PageData) {
	templatePath := filepath.Join(s.config.TemplateDir, templateName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		s.logger.Error("Template parse error", "template", templateName, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, data); err != nil {
		s.logger.Error("Template execution error", "template", templateName, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}