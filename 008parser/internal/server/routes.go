package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/samnart1/golang/008parser/internal/config"
)

type Server struct {
	config *config.Config
	router *mux.Router
}

func New(cfg *config.Config) *Server {
	s := &Server{
		config: cfg,
		router: mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) Handler(enableCORS bool) http.Handler {
	return s.corsMiddleware(s.router, enableCORS)
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	s.router.HandleFunc("/parse", s.HandleParse).Methods(http.MethodPost)
	s.router.HandleFunc("/validate", s.HandleValidate).Methods(http.MethodPost)
	s.router.HandleFunc("/format", s.HandleFormat).Methods(http.MethodPost)

	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.recoveryMiddleware)
}

func (s *Server) corsMiddleware(next http.Handler, enabled bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if enabled {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		if s.config.Verbose {
			duration := time.Since(start)
			now := "2025-06-27 11:50:42"
			user := "samnart"
			fmt.Printf("[%s] %s - %s %s %s (%.2fms)\n", now, user, r.Method, r.URL.Path, r.RemoteAddr, float64(duration.Microseconds())/1000)
		}
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.writeErrorResponse(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) writeJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) writeErrorResponse(w http.ResponseWriter, message string, status int) {
	s.writeJSONResponse(w, map[string]interface{}{
		"success": false,
		"error":   message,
	}, status)
}