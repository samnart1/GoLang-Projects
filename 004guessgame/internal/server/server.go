package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	Host		string
	Port		int
	Debug 		bool
	Timeout		time.Time
}

type Server struct {
	config	*Config
	router	*mux.Router
	server	*http.Server
}