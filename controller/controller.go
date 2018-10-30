package controller

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func AddToManager(mgr manager.Manager) error {
	r := &ReconcilePod{client: mgr.GetClient(), scheme: mgr.GetScheme()}
	c, err := controller.New("mycontroller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	return c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcilePod{}

type ReconcilePod struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcilePod) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconciling pod: ---> %s/%s\n", request.Namespace, request.Name)

	instance := &corev1.Pod{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("---> pod is deleted\n")
		} else {
			log.Printf("---> get unexpected error %v\n", err)
		}
	} else {
		log.Printf("pod event %v\n", instance.ObjectMeta)
	}

	return reconcile.Result{}, nil
}
