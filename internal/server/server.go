package server

import (
	"context"
	"fmt"
	"github.com/shamank/warehouse-service/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(httpConfig config.HTTPConfig) *Server {

	fmt.Println(httpConfig)
	return &Server{
		httpServer: &http.Server{
			Addr:           httpConfig.Host + ":" + httpConfig.Port,
			Handler:        nil,
			WriteTimeout:   httpConfig.WriteTimeOut,
			ReadTimeout:    httpConfig.ReadTimeOut,
			MaxHeaderBytes: httpConfig.MaxHeaderBytes << 20,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) SetHandler(handler http.Handler) {
	s.httpServer.Handler = handler
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
