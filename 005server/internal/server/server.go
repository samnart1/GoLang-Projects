package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/samnart1/GoLang-Projects/005server/internal/config"
	"github.com/samnart1/GoLang-Projects/005server/internal/log"
)

type Server struct {
	httpServer *http.Server
	config *config.Config
	logger *log.Logger
}

func New(cfg *config.Config, log *log.Logger) *Server {
	s := &Server{
		config: cfg,
		logger: log,
	}

	s.httpServer = &http.Server{
		Addr: ":" + cfg.Port,
		Handler: s.setupRoutes(),
		ReadTimeout: cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout: cfg.IdleTimeout,
	}

	return s
}

func (s *Server) Start() error {
	s.logger.Info("Server starting", "address", s.httpServer.Addr)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Server failed to start: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Server shutting down...")

	//shutdown the http server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.logger.Info("Server shutdown complete")
	return nil
}

func (s *Server) setupRoutes() http.Handler {
	mux := http.NewServeMux()

	// apply middleware chain
	handler := s.applyMiddleware(mux)

	//register routes
	s.registerRoutes(mux)
	
	return handler
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	//main
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/hello", s.handleHello)
	mux.HandleFunc("/about", s.handleAbout)

	//api routes
	mux.HandleFunc("/api/info", s.handleAPIInfo)
	mux.HandleFunc("/api/echo", s.handleAPIEcho)
	mux.HandleFunc("/api/time", s.handleAPITime)

	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ready", s.handleReady)
	
	fileServer := http.FileServer(http.Dir(s.config.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
}

func (s *Server) applyMiddleware(handler http.Handler) http.Handler {
	handler = s.recoveryMiddleware(handler)
	handler = s.loggingMiddleware(handler)

	if s.config.EnableCORS {
		handler = s.corsMiddleware(handler)
	}

	return handler
}