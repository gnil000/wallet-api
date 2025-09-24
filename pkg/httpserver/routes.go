package httpserver

import (
	"net/http"
)

type Router interface {
	GET(relativePath string, handler http.HandlerFunc) Router
	POST(relativePath string, handler http.HandlerFunc) Router
}

func (s *server) GET(relativePath string, handler http.HandlerFunc) Router {
	pattern := http.MethodGet + " /api/v1" + relativePath
	s.mux.Handle(pattern, handler)
	return s
}

func (s *server) POST(relativePath string, handler http.HandlerFunc) Router {
	pattern := http.MethodPost + " /api/v1" + relativePath
	s.mux.Handle(pattern, handler)
	return s
}
