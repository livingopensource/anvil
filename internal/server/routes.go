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
	api.HandleFunc("/hash/{email}", s.hashGen).Methods(http.MethodGet)

	instances := api.PathPrefix("/virtual-machines").Subrouter()
	instances.HandleFunc("", s.ListVMInstancesHandler).Methods(http.MethodGet)
	instances.HandleFunc("", s.CreateVMInstanceHandler).Methods(http.MethodPost)
	instances.HandleFunc("/{name}", s.GetVMInstanceHandler).Methods(http.MethodGet)
	instances.HandleFunc("/{name}", s.DeleteVMInstanceHandler).Methods(http.MethodDelete)
	instances.HandleFunc("/{name}", s.UpdateVMInstanceHandler).Methods(http.MethodPut)
	instances.HandleFunc("/{name}/vnc", s.VNCInstanceHandler).Methods(http.MethodGet)

	return r
}
