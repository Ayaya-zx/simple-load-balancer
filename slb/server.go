package slb

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type option func(*Server)

type Server struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	ready  string
}

func NewServer(target *url.URL, opts ...option) *Server {
	server := &Server{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func WithReadyCheck(check string) option {
	return func(s *Server) {
		s.ready = check
	}
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}

func (s *Server) isReady() bool {
	if s.ready == "" {
		return true
	}
	r, err := http.Get(fmt.Sprintf("%s://%s%s",
		s.target.Scheme, s.target.Host, s.ready))
	return err == nil && r.StatusCode >= 200 && r.StatusCode < 400
}
