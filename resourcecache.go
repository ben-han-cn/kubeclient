package kubeclient

import (
	"errors"

	api "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var namespace = api.NamespaceAll
var (
	ErrUnknownResourceType = errors.New("resource type is unknown")
)

type ResourceCache struct {
	client     *kubernetes.Clientset
	indexer    cache.Indexer
	controller cache.Controller
	stopCh     chan struct{}
}

func NewResourceCache(client *kubernetes.Clientset, resource string, owner ResourceController, opts *ControllerOptions) (*ResourceCache, error) {
	resourceType, ok := resourceMap[resource]
	if ok == false {
		return nil, ErrUnknownResourceType
	}
	indexer, controller := cache.NewIndexerInformer(
		cache.NewFilteredListWatchFromClient(client.CoreV1().RESTClient(),
			resource,
			opts.Namespace,
			defaultOptionModifier(opts)),
		resourceType,
		opts.ResyncPeriod,
		owner,
		owner.GetResourceIndexers(resource),
	)

	return &ResourceCache{
		client:     client,
		controller: controller,
		indexer:    indexer,
		stopCh:     make(chan struct{}),
	}, nil
}

func (c *ResourceCache) Indexer() cache.Indexer {
	return c.indexer
}

func (c *ResourceCache) Run() {
	c.controller.Run(c.stopCh)
}

func (c *ResourceCache) Stop() {
	close(c.stopCh)
}
