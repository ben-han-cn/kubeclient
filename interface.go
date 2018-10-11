package kubeclient

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

type ControllerOptions struct {
	LabelSelector labels.Selector
	FieldSelector fields.Selector
	ResyncPeriod  time.Duration
	Namespace     string
}

type ResourceController interface {
	cache.ResourceEventHandler
	GetResourceIndexers(string) cache.Indexers
}

func defaultOptionModifier(opts *ControllerOptions) func(opts *metav1.ListOptions) {
	return func(opts_ *metav1.ListOptions) {
		if opts.LabelSelector != nil {
			opts_.LabelSelector = opts.LabelSelector.String()
		}
		if opts.FieldSelector != nil {
			opts_.FieldSelector = opts.FieldSelector.String()
		}
	}
}
