package kubeclient

import (
	apps "k8s.io/api/apps/v1"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

//key will be used in url
//value is the resource real type
var resourceMap = map[string]runtime.Object{
	"pods":            &api.Pod{},
	"services":        &api.Service{},
	"endpoints":       &api.Endpoints{},
	"nodes":           &api.Node{},
	"namespaces":      &api.Namespace{},
	"serviceaccounts": &api.ServiceAccount{},
	"events":          &api.Event{},
	"configmaps":      &api.ConfigMap{},

	"daemonsets":  &apps.DaemonSet{},
	"deployments": &apps.Deployment{},
}
