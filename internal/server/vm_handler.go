package server

import (
	"cloud/internal/vm"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
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

func (s *Server) VNCVMInstanceHandler(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	project := r.URL.Query().Get("project")
	if project == "" {
		crw.response(http.StatusBadRequest, "project is required", nil, nil)
		return
	}
	req := newRequest("vm", r)
	resource := req.useProject(project)
	virtualMachine := vm.NewCluster(resource)
	vmInstance, err := virtualMachine.VNC()
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
	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	vmiConn := vmInstance.AsConn()
	// Start copying data between WebSocket and VMI console
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("Error reading from websocket: " + err.Error())
				break
			}
			_, err = vmiConn.Write(message)
			if err != nil {
				slog.Error("Error writing to VMI console: n" + err.Error())
				break
			}
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := vmiConn.Read(buf)
		if err != nil {
			slog.Error("Error reading from VMI console: " + err.Error())
			break
		}
		err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
		if err != nil {
			slog.Error("Error writing to websocket: " + err.Error())
			break
		}
	}
}

func (s *Server) WatchVMInstanceHandler(w http.ResponseWriter, r *http.Request) {
	crw := customResponseWriter{w: w}
	project := r.URL.Query().Get("project")
	if project == "" {
		crw.response(http.StatusBadRequest, "project is required", nil, nil)
		return
	}
	req := newRequest("vm", r)
	resource := req.useProject(project)
	virtualMachine := vm.NewCluster(resource)
	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		crw.response(http.StatusInternalServerError, err.Error(), nil, nil)
		return
	}

	watcher, err := virtualMachine.Watch()
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
	defer watcher.Stop()

	// Stream events to the websocket
	for event := range watcher.ResultChan() {
		jsonData, err := json.Marshal(event)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			slog.Error(err.Error())
			break
		}
	}
}
