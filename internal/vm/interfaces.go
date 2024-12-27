package vm

import (
	"cloud/internal/clusters"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kvV1 "kubevirt.io/client-go/generated/kubevirt/clientset/versioned/typed/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

type VirtualMachine struct {
	ctx        context.Context
	kubeconfig string
	project    string
	request    *http.Request
}

func NewCluster(req clusters.Resource) *VirtualMachine {
	return &VirtualMachine{
		ctx:        req.Ctx,
		kubeconfig: req.Kubeconfig,
		project:    req.Project,
		request:    req.Request,
	}
}

func (vm *VirtualMachine) Create() error {
	// TODO: Prevent users from creating an insance with name watch, 
	// as creating such an instance will prevent thr router from listing 
	// the instance, instead serving the endpoint for watching websockets
	payload, err := clusters.Payload(vm.request)
	if err != nil {
		return err
	}
	name := payload.User.Name
	passwd := payload.User.Password
	cloudInitConfig := fmt.Sprintf(`#cloud-config
users:
  - name: %s
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: users
    home: /home/%s
    shell: /bin/bash
    lock_passwd: false
chpasswd:
  list: |
    %s:%s
  expire: False`, name, name, name, passwd)
	cloudInitBase64 := base64.StdEncoding.EncodeToString([]byte(cloudInitConfig))
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kubevirt.io/v1",
			"kind":       "VirtualMachine",
			"metadata": map[string]interface{}{
				"name": payload.Compute.Name,
			},
			"spec": map[string]interface{}{
				"runStrategy": "RerunOnFailure",
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"kubevirt.io/vm": payload.Compute.Name,
						},
					},
					"spec": map[string]interface{}{
						"domain": map[string]interface{}{
							"cpu": map[string]interface{}{
								"cores": payload.Compute.CPU,
							},
							"devices": map[string]interface{}{
								"disks": []map[string]interface{}{
									{
										"name": "os-disk-" + payload.Compute.Name,
										"disk": map[string]interface{}{
											"bus": "virtio",
										},
									},
									{
										"name": "cloudinitdisk",
										"cdrom": map[string]interface{}{
											"bus": "sata",
										},
									},
								},
								"interfaces": []map[string]interface{}{
									{
										"name":       "default",
										"masquerade": map[string]interface{}{},
									},
								},
							},
							"resources": map[string]interface{}{
								"limits": map[string]interface{}{
									"memory": payload.Compute.RAM,
								},
							},
						},
						"networks": []map[string]interface{}{
							{
								"name": "default",
								"pod":  map[string]interface{}{},
							},
						},
						"volumes": []map[string]interface{}{
							{
								"name": "os-disk-" + payload.Compute.Name,
								"dataVolume": map[string]interface{}{
									"name": "os-volume-disk-" + payload.Compute.Name,
								},
							},
							{
								"name": "cloudinitdisk",
								"cloudInitNoCloud": map[string]interface{}{
									"userDataBase64": cloudInitBase64,
								},
							},
						},
					},
				},
				"dataVolumeTemplates": []map[string]interface{}{
					{
						"apiVersion": "cdi.kubevirt.io/v1beta1",
						"kind":       "DataVolume",
						"metadata": map[string]interface{}{
							"name": "os-volume-disk-" + payload.Compute.Name,
						},
						"spec": map[string]interface{}{
							"storage": map[string]interface{}{
								"accessModes": []string{
									"ReadWriteOnce",
								},
								"resources": map[string]interface{}{
									"requests": map[string]interface{}{
										"storage": payload.Compute.Storage,
									},
								},
							},
							"source": map[string]interface{}{
								"http": map[string]interface{}{
									"url": payload.Compute.URL,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = clusters.CreateResourceSchema(obj, vm.kubeconfig, vm.project)
	if err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachine) Delete() error {
	vars := mux.Vars(vm.request)
	name := vars["name"]
	return clusters.DeleteResourceSchema(schema.GroupVersionKind{
		Group:   "kubevirt.io",
		Version: "v1",
		Kind:    "VirtualMachine",
	}, name, vm.kubeconfig, vm.project)
}

func (vm *VirtualMachine) Find() (map[string]interface{}, error) {
	vars := mux.Vars(vm.request)
	name := vars["name"]
	gvk := schema.GroupVersionKind{
		Group:   "kubevirt.io",
		Version: "v1",
		Kind:    "VirtualMachine",
	}
	if vm.request.URL.Query().Get("state") == "up" {
		gvk = schema.GroupVersionKind{
			Group:   "kubevirt.io",
			Version: "v1",
			Kind:    "VirtualMachineInstance",
		}
	}
	response, err := clusters.GetResourceSchema(gvk, name, vm.kubeconfig, vm.project)
	if err != nil {
		return nil, err
	}
	return response.Object, nil
}

func (vm *VirtualMachine) FindAll() ([]map[string]interface{}, error) {
	gvk := schema.GroupVersionKind{
		Group:   "kubevirt.io",
		Version: "v1",
		Kind:    "VirtualMachine",
	}
	if vm.request.URL.Query().Get("state") == "up" {
		gvk = schema.GroupVersionKind{
			Group:   "kubevirt.io",
			Version: "v1",
			Kind:    "VirtualMachineInstance",
		}
	}
	response, err := clusters.ListResourceSchema(gvk, vm.kubeconfig, vm.project)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, len(response.Items))
	for i, item := range response.Items {
		result[i] = item.Object
	}
	return result, nil
}

func (vm *VirtualMachine) Patch() (map[string]interface{}, error) {
	return nil, errors.New("not implemented yet")
}

func (vm *VirtualMachine) Watch() (watch.Interface, error) {
	return nil, errors.New("not implemented yet")
}

func (vm *VirtualMachine) VNC() (kvV1.StreamInterface, error) {
	vars := mux.Vars(vm.request)
	name := vars["name"]
	kubevirt, err := clusters.KubevirtResourceSchema(vm.kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubevirt.VirtualMachineInstance(vm.project).VNC(name)
}
