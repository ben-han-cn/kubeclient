package kubeclient

import (
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

//key will be used in url
//value is the resource real type
var resourceMap = map[string]runtime.Object{
	"pods":      &api.Pod{},
	"services":  &api.Service{},
	"endpoints": &api.Endpoints{},
}
