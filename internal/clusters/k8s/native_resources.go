package k8s

import (
	"context"

	authenticationv1 "k8s.io/api/authentication/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resource struct {
	kubeconfig string
}

func (r Resource) Pods(ns string) (*v1.PodList, error) {
	clientSet, err := ClientSet(r.kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
}

func (r Resource) Secrets(ns string) (*v1.SecretList, error) {
	clientSet, err := ClientSet(r.kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientSet.CoreV1().Secrets(ns).List(context.Background(), metav1.ListOptions{})
}

func (r Resource) ConfigMaps(ns string) (*v1.ConfigMapList, error) {
	clientSet, err := ClientSet(r.kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientSet.CoreV1().ConfigMaps(ns).List(context.Background(), metav1.ListOptions{})
}

func (r Resource) CreateToken(ns, serviceAccount string) (*authenticationv1.TokenRequest, error) {
	clientSet, err := ClientSet(r.kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientSet.CoreV1().ServiceAccounts(ns).CreateToken(context.Background(), serviceAccount, &authenticationv1.TokenRequest{}, metav1.CreateOptions{})
}
