package server

import (
	"cloud/internal/clusters"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type CloudRequest struct {
	ctx        context.Context
	project    string
	request    *http.Request
	resource   string
	kubeconfig string
}

func newRequest(resource string, request *http.Request) *CloudRequest {
	rid := uuid.New().String()
	ctx := context.WithValue(request.Context(), "request_id", rid)
	kubeconfig := ""
	switch resource {
	case "vm":
		kubeconfig = viper.GetString("cluster.vm")
	}

	return &CloudRequest{
		ctx:        ctx,
		request:    request,
		resource:   resource,
		kubeconfig: kubeconfig,
	}
}

func (c *CloudRequest) useProject(name string) clusters.Resource {
	return clusters.Resource{
		Ctx:        c.ctx,
		Kubeconfig: c.kubeconfig,
		Project:    name,
		Request:    c.request,
	}
}
