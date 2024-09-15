package server

import (
	"context"
	"net/http"

	"github.com/b0pof/avito-internship/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg config.Server, r http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.ServerAddr,
			Handler: r,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
