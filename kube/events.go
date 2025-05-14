package kube

import (
	"context"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func GetKubeClient() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		fmt.Println("Error loading kubeconfig:", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating clientset:", err)
		os.Exit(1)
	}
	return clientset
}

func WatchEvents(namespace string, eventHandler func(event *v1.Event)) {
	clientset := GetKubeClient()
	watcher, err := clientset.CoreV1().Events(namespace).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Failed to watch events:", err)
		return
	}
	ch := watcher.ResultChan()

	for evt := range ch {
		if event, ok := evt.Object.(*v1.Event); ok {
			eventHandler(event)
		}
	}
}