package http

import (
	"github.com/gorilla/mux"
)

// Handler will initialize mux router and register handler
func (s *Server) Handler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/{id:[0-9]+}", s.User.GetUserHandler).Methods("GET")
	return r
}
