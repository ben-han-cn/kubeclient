package kubeclient

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func DefaultKubeconfig() string {
	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	} else {
		return ""
	}
}

func NewClientSet(master, kubeconfig string) (*kubernetes.Clientset, error) {
	c, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		return nil, err
	}
	return clientForConfig(c)
}

func NewClientSetInsideCluster() (*kubernetes.Clientset, error) {
	c, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return clientForConfig(c)
}

func clientForConfig(c *rest.Config) (*kubernetes.Clientset, error) {
	c.ContentType = "application/vnd.kubernetes.protobuf"
	return kubernetes.NewForConfig(c)
}
