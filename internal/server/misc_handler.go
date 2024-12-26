package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// responseBody is the response body that is returned to the client.
type responseBody struct {
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
}

// customResponseWriter is a custom response writer that implements the http.ResponseWriter interface.
type customResponseWriter struct {
	w http.ResponseWriter
}

// response writes a response to the client in a json marshalled format.
func (rw customResponseWriter) response(status int, message string, payload interface{}, pagination interface{}) {
	resp := responseBody{
			Data:    payload,
			Meta:    pagination,
			Status:  status,
			Message: message,
	}
	data, err := json.Marshal(resp)
	if err != nil {
			slog.Error(err.Error(), "status_code", http.StatusInternalServerError)
			rw.w.WriteHeader(http.StatusInternalServerError)
			rw.w.Write(data)
			return
	}
	rw.w.WriteHeader(status)
	rw.w.Write(data)
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("error handling JSON marshal.",  "error", err.Error())
	}
	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp := []byte(`{"alive": true}`)
	_, _ = w.Write(jsonResp)
}