package main

import (
	"fmt"
	"log"
	"os"
	osig "os/signal"
	"syscall"
	"time"

	"github.com/ben-han-cn/kubeclient"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

type DumbController struct {
}

var _ kubeclient.ResourceController = &DumbController{}

func (c *DumbController) OnAdd(obj interface{}) {
	switch o := obj.(type) {
	case *api.Service:
		fmt.Printf("new service: %v\n", o.Name)
	case *api.Pod:
		fmt.Printf("new pod: %v\n", o.Name)
	default:
		fmt.Printf("unknown obj\n")
	}
}

func (c *DumbController) OnUpdate(old, new interface{}) {
	if old.(metav1.Object).GetResourceVersion() == new.(metav1.Object).GetResourceVersion() {
		return
	}

	switch old.(type) {
	case *api.Service:
		oldService := old.(*api.Service)
		newService := new.(*api.Service)
		fmt.Printf("update serfvice from %v to : %v\n", oldService.Name, newService.Name)
	case *api.Pod:
		oldPod := old.(*api.Pod)
		newPod := new.(*api.Pod)
		fmt.Printf("update pod from %v to : %v\n", oldPod.Name, newPod.Name)
	default:
		fmt.Printf("unknown obj\n")
	}

}

func (c *DumbController) OnDelete(obj interface{}) {
	switch o := obj.(type) {
	case *api.Service:
		fmt.Printf("delete service: %v\n", o.Name)
	case *api.Pod:
		fmt.Printf("delete pod: %v\n", o.Name)
	default:
		fmt.Printf("unknown obj\n")
	}
}

func (c *DumbController) GetResourceIndexers(resource string) cache.Indexers {
	return nil
}

func WaitForInterrupt(cb func()) {
	signalCh := make(chan os.Signal)
	osig.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh

	if cb != nil {
		cb()
	}
}

func main() {
	client, err := kubeclient.NewClientSet("", kubeclient.DefaultKubeconfig())
	if err != nil {
		log.Fatalf("connect to api server failed:%s", err.Error())
	}

	opts := &kubeclient.ControllerOptions{
		ResyncPeriod: 10 * time.Second,
		Namespace:    "default",
	}

	controller := &DumbController{}
	podCache, _ := kubeclient.NewResourceCache(client, "pods", controller, opts)
	go podCache.Run()
	serviceCache, _ := kubeclient.NewResourceCache(client, "services", controller, opts)
	go serviceCache.Run()

	signal.WaitForInterrupt(func() {
		podCache.Stop()
		serviceCache.Stop()
		log.Print("existing !!")
	})
}
