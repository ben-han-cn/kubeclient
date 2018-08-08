package kubeclient

import (
	api "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	IPTOPod = "iptopod"
)

func PodIPIndexers() cache.Indexers {
	return cache.Indexers{
		IPTOPod: func(obj interface{}) ([]string, error) {
			p, _ := obj.(*api.Pod)
			return []string{p.Status.PodIP}, nil
		},
	}
}

func PodsWithIP(indexer cache.Indexer, ip string) ([]*api.Pod, error) {
	objs, err := indexer.ByIndex(IPTOPod, ip)
	if err != nil {
		return nil, err
	}

	var pods []*api.Pod
	for _, o := range objs {
		p, ok := o.(*api.Pod)
		if !ok {
			panic("wrong object in indexer")
		}
		pods = append(pods, p)
	}
	return pods, nil
}
