package api

import (
	"api-gateway/api/handlers"
	"context"
	"net/http"
)

type ServerBuilder struct {
	server http.Server
	router *http.ServeMux
	
}

func NewServerBuilder(ctx context.Context) *ServerBuilder {
	router := http.NewServeMux()

	return &ServerBuilder {
		server: http.Server{},
		router: router,
	}
}

func (s *ServerBuilder) AddHandlers() {
	s.router.HandleFunc("/", handlers.PingHandler)
	s.router.HandleFunc("/health", handlers.HealthCheck)
}

func (s *ServerBuilder) Builder(addr string) error {
	s.server.Addr = addr
	s.server.Handler = s.router
	return s.server.ListenAndServe()
}

func (s *ServerBuilder) GetRouter() *http.ServeMux {
	return s.router
}
