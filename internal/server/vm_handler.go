package server

import (
	"cloud/internal/vm"
	"log/slog"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"
)

func (s *Server) ListVMInstancesHandler(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	project := r.URL.Query().Get("project")
	if project == "" {
		crw.response(http.StatusBadRequest, "project is required", nil, nil)
		return
	}
	req := newRequest("vm", r)
	resource := req.useProject(project)
	virtualMachine := vm.NewCluster(resource)
	vms, err := virtualMachine.FindAll()
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	crw.response(http.StatusOK, "success", vms, nil)
}

func (s *Server) GetVMInstanceHandler(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	project := r.URL.Query().Get("project")
	if project == "" {
		crw.response(http.StatusBadRequest, "project is required", nil, nil)
		return
	}
	req := newRequest("vm", r)
	resource := req.useProject(project)
	virtualMachine := vm.NewCluster(resource)
	vm, err := virtualMachine.Find()
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	crw.response(http.StatusOK, "success", vm, nil)
}

func (s *Server) DeleteVMInstanceHandler(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	project := r.URL.Query().Get("project")
	if project == "" {
		crw.response(http.StatusBadRequest, "project is required", nil, nil)
		return
	}
	req := newRequest("vm", r)
	resource := req.useProject(project)
	virtualMachine := vm.NewCluster(resource)
	err := virtualMachine.Delete()
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	crw.response(http.StatusOK, "success", nil, nil)
}

func (s *Server) UpdateVMInstanceHandler(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	project := r.URL.Query().Get("project")
	if project == "" {
		crw.response(http.StatusBadRequest, "project is required", nil, nil)
		return
	}
	req := newRequest("vm", r)
	resource := req.useProject(project)
	virtualMachine := vm.NewCluster(resource)
	vm, err := virtualMachine.Patch()
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	crw.response(http.StatusOK, "success", vm, nil)
}

func (s *Server) CreateVMInstanceHandler(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	project := r.URL.Query().Get("project")
	if project == "" {
		crw.response(http.StatusBadRequest, "project is required", nil, nil)
		return
	}
	req := newRequest("vm", r)
	resource := req.useProject(project)
	virtualMachine := vm.NewCluster(resource)
	err := virtualMachine.Create()
	if err != nil {
		statusError, isStatus := err.(*errors.StatusError)
		if isStatus {
			errCode := statusError.Status().Code
			slog.Error("Kubernetes error", "code", errCode, "message", err.Error())
			crw.response(int(errCode), err.Error(), nil, nil)
		} else {
			slog.Error("Unknown error", "message", err.Error())
			crw.response(http.StatusUnprocessableEntity, err.Error(), nil, nil)
		}
		return
	}
	crw.response(http.StatusOK, "success", nil, nil)
}
