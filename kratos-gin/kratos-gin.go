package gin

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kHttp "github.com/go-kratos/kratos/v2/transport/http"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

type Server struct {
	*gin.Engine
	server *http.Server

	err error

	tlsConf  *tls.Config
	endpoint *url.URL
	timeout  time.Duration
	addr     string

	filters []kHttp.FilterFunc
	ms      []middleware.Middleware
	dec     kHttp.DecodeRequestFunc  // req
	enc     kHttp.EncodeResponseFunc // rsp
	ene     kHttp.EncodeErrorFunc
}

// 初始化配置
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		timeout: 1 * time.Second,
		// DefaultRequestDecoder 将请求正文解码为对象。
		dec: kHttp.DefaultRequestDecoder,
		// DefaultResponseEncoder 将对象编码为 HTTP response
		enc: kHttp.DefaultResponseEncoder,
		// DefaultErrorEncoder 将错误编码为 HTTP response
		ene: kHttp.DefaultErrorEncoder,
	}

	srv.init(opts...)
	return srv
}

func (s *Server) init(opts ...ServerOption) {
	// gin engine
	s.Engine = gin.New()

	for _, o := range opts {
		o(s) // with 配置
	}

	s.server = &http.Server{
		Addr:      s.addr,
		Handler:   s.Engine,
		TLSConfig: s.tlsConf,
	}
}

// kratos transport server 规范 start
func (s *Server) Start(ctx context.Context) error {
	// 启动 server
	log.Infof("[GIN] server listening on: %s", s.addr)

	var err error
	// 辨别为http/https  Transport Layer Security(tls)
	if s.tlsConf != nil {
		err = s.server.ListenAndServeTLS("", "")
	} else {
		err = s.server.ListenAndServe()
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return err
}

// kratos transport server 规范 stop
func (s *Server) Stop(ctx context.Context) error {
	// 关闭 server
	log.Infof("[GIN] server stopping")
	return s.server.Shutdown(ctx)
}

// kratos transport endpointer 规范 endpoint
func (s *Server) Endpoint() (*url.URL, error) {
	// 解析 ip 地址 http/https
	addr := s.addr // 域名/ip

	prefix := ""
	if s.tlsConf != nil {
		if !strings.Contains(addr, "https://") {
			prefix = "https://"
		}
	} else {
		if !strings.Contains(addr, "http://") {
			prefix = "http://"
		}
	}
	addr = prefix + addr
	var endpoint *url.URL
	// 解析 addr
	endpoint, s.err = url.Parse(addr)

	return endpoint, s.err
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.Engine.ServeHTTP(res, req)
}
