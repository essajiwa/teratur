package http

import (
	"net/http"

	"github.com/essajiwa/teratur/pkg/grace"
)

type UserHandler interface {
	GetUserHandler(w http.ResponseWriter, r *http.Request)
}

// Server holds http server object and services that will be injected to start the http handler
type Server struct {
	server *http.Server
	User   UserHandler
}

// Serve will inject the service used at handler and then run an HTTP server on provided port
func (s *Server) Serve(port string) error {
	return grace.Serve(port, s.Handler())
}
