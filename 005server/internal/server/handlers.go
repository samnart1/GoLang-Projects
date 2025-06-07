package server

import (
	"fmt"
	"net/http"
	"runtime"
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