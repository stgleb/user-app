package server

import "net/http"

type Server struct {
	srv *http.Server
}

func NewServer(addr string) *Server {
	return nil
}
