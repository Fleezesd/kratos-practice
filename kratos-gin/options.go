package gin

import (
	"crypto/tls"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	kHttp "github.com/go-kratos/kratos/v2/transport/http"
)

type ServerOption func(*Server)

func WithTLSConfig(c *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = c
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithFilter(filters ...kHttp.FilterFunc) ServerOption {
	return func(s *Server) {
		s.filters = filters
	}
}

func WithMiddleware(m ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.ms = m
	}
}

func WithRequestDecoder(dec kHttp.DecodeRequestFunc) ServerOption {
	return func(s *Server) {
		s.dec = dec
	}
}

func WithResponseEncoder(en kHttp.EncodeResponseFunc) ServerOption {
	return func(s *Server) {
		s.enc = en
	}
}

func WithErrorEncoder(en kHttp.EncodeErrorFunc) ServerOption {
	return func(s *Server) {
		s.ene = en
	}
}
