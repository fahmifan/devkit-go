package httpsvc

import (
	"github.com/fahmifan/devkit/pkg/service"
)

type Option func(*Server)

func WithService(s *service.Service) Option {
	return func(srv *Server) {
		srv.service = s
	}
}

func WithJWTKey(key string) Option {
	return func(s *Server) {
		s.jwtKey = key
	}
}
