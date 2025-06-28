package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samnart1/golang/009l9gger/internal/config"
	"github.com/samnart1/golang/009l9gger/internal/logger"
)



type Server struct {
	config *config.Config
	logger *logger.Logger
	router *mux.Router
}

func New(cfg *config.Config, l *logger.Logger) *Server {
	s := &Server{
		config: cfg,
		logger: l,
		router: mux.NewRouter(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) Router() http.Handler {
	return s.router
}

func (s *Server) setupRoutes() {
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)

	s.router.HandleFunc("/", s.handleRoot).Methods("GET")
	s.router.HandleFunc("log", s.handleLog).Methods("POST")
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := logger.Entry{
			Message: fmt.Sprintf("%s %s", r.Method, r.URL.Path),
			Level:	"info",
			Source:	"http-server",
			Data:	map[string]interface{}{
				"method":		r.Method,
				"path":			r.URL.Path,
				"remote_addr": 	r.RemoteAddr,
				"user_agent":	r.UserAgent(),
			},
		}

		s.logger.Log(entry)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}