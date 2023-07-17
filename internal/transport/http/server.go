package http

import (
	"fmt"
	"hezzl/config"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(cfg config.ServerConfig, h http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Handler:      h,
			WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
			ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
			Addr:         cfg.Addr,
		},
	}
}

func (s Server) Run() error {
	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("error while trying to run server: %w", err)
	}
	return nil
}
