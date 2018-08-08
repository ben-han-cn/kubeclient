package main

import (
	"fmt"
	"log"
	"time"

	"cement/signal"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/kubeclient"
)

type DumbController struct {
	opts *kubeclient.ControllerOptions
}

var _ kubeclient.ResourceController = &DumbController{}

func newDumbController() *DumbController {
	opts := &kubeclient.ControllerOptions{
		ResyncPeriod: 10 * time.Second,
		Namespace:    "default",
	}
	return &DumbController{
		opts: opts,
	}
}

func (c *DumbController) OnAdd(obj interface{}) {
	pod := obj.(*api.Pod)
	fmt.Printf("new pod: %v\n", pod)
}

func (c *DumbController) OnUpdate(old, new interface{}) {
	if old.(metav1.Object).GetResourceVersion() == new.(metav1.Object).GetResourceVersion() {
		return
	}

	oldPod := old.(*api.Pod)
	newPod := new.(*api.Pod)
	fmt.Printf("update pod from %v to : %v\n", oldPod, newPod)
}

func (c *DumbController) OnDelete(obj interface{}) {
	pod := obj.(*api.Pod)
	fmt.Printf("delete pod: %v\n", pod)
}

func (c *DumbController) GetOptions() *kubeclient.ControllerOptions {
	return c.opts
}

func main() {
	client, err := kubeclient.NewClientSet("", kubeclient.DefaultKubeconfig())
	if err != nil {
		log.Fatalf("connect to api server failed:%s", err.Error)
	}

	cache := kubeclient.NewPodCache(client, newDumbController())
	go cache.Run()

	signal.WaitForInterrupt(func() {
		cache.Stop()
		log.Print("existing !!")
	})
}
