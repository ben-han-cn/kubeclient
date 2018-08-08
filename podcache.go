package kubeclient

import (
	api "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var namespace = api.NamespaceAll

const (
	podIPIndex = "PodIP"
)

type PodCache struct {
	client        *kubernetes.Clientset
	indexer       cache.Indexer
	podController cache.Controller
	owner         ResourceController
	stopCh        chan struct{}
}

func NewPodCache(client *kubernetes.Clientset, controller ResourceController) *PodCache {
	opts := controller.GetOptions()
	indexer, podController := cache.NewIndexerInformer(
		cache.NewFilteredListWatchFromClient(client.CoreV1().RESTClient(),
			"pods",
			opts.Namespace,
			defaultOptionModifier(opts)),
		&api.Pod{},
		opts.ResyncPeriod,
		controller,
		cache.Indexers{podIPIndex: func(obj interface{}) ([]string, error) {
			p, _ := obj.(*api.Pod)
			return []string{p.Status.PodIP}, nil
		}})
	return &PodCache{
		client:        client,
		podController: podController,
		owner:         controller,
		indexer:       indexer,
		stopCh:        make(chan struct{}),
	}
}

func (c *PodCache) PodIndex(ip string) (pods []*api.Pod) {
	os, err := c.indexer.ByIndex(podIPIndex, ip)
	if err != nil {
		return nil
	}
	for _, o := range os {
		p, ok := o.(*api.Pod)
		if !ok {
			continue
		}
		pods = append(pods, p)
	}
	return pods
}

func (c *PodCache) Run() {
	c.podController.Run(c.stopCh)
}

func (c *PodCache) Stop() {
	close(c.stopCh)
}
