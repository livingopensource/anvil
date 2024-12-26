package clusters

import (
	"context"
	"encoding/json"
	"net/http"
)

type Resource struct {
	Ctx        context.Context
	Kubeconfig string
	Project    string
	Request    *http.Request
}

type Compute struct {
	Name      string      `json:"name,omitempty"`
	CPU       float64     `json:"vcpu,omitempty"`
	RAM       string      `json:"ram,omitempty"`
	Storage   string      `json:"storage,omitempty"`
	Instances float64     `json:"instances,omitempty"`
	State     string      `json:"state,omitempty"`
	SSHKey    string      `json:"ssh_key,omitempty"`
	URL       string      `json:"url,omitempty"`
	Container []Container `json:"containers,omitempty"`
}

type User struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type Container struct {
	Image string `json:"image,omitempty"`
	Port  []Port `json:"ports,omitempty"`
	Env   []Env  `json:"env,omitempty"`
}

type Port struct {
	ContainerPort int32 `json:"containerPort,omitempty"`
}

type Env struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type ResourceDetails struct {
	ID        string  `json:"id,omitempty"`
	Compute   Compute `json:"compute,omitempty"`
	User      User    `json:"user,omitempty"`
}

// Payload is a decoded json request payload
func Payload(r *http.Request) (ResourceDetails, error) {
	var payload ResourceDetails
	err := json.NewDecoder(r.Body).Decode(&payload)
	return payload, err
}
