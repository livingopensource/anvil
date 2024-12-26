package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/1.0").Subrouter()

	api.HandleFunc("/", s.HelloWorldHandler)
	api.HandleFunc("/health", s.healthHandler)

	instances := api.PathPrefix("/virtual-machines").Subrouter()
	instances.HandleFunc("", s.ListVMInstancesHandler).Methods(http.MethodGet)
	instances.HandleFunc("", s.CreateVMInstanceHandler).Methods(http.MethodPost)
	instances.HandleFunc("/{name}", s.GetVMInstanceHandler).Methods(http.MethodGet)
	instances.HandleFunc("/{name}", s.DeleteVMInstanceHandler).Methods(http.MethodDelete)
	instances.HandleFunc("/{name}", s.UpdateVMInstanceHandler).Methods(http.MethodPut)

	return r
}
