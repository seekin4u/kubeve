package kube

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

func WatchEvents(namespace string, eventHandler func(event *corev1.Event)) {
	clientset := GetKubeClient()
	ctx := context.TODO()
	evList, err := clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error listing events:", err)
		return
	}
	resourceVersion := evList.ResourceVersion

	watcher, err := clientset.CoreV1().Events(namespace).Watch(ctx, metav1.ListOptions{
		ResourceVersion: resourceVersion,
	})
	if err != nil {
		fmt.Println("Failed to watch events:", err)
		return
	}
	ch := watcher.ResultChan()

	for evt := range ch {
		if event, ok := evt.Object.(*corev1.Event); ok {
			eventHandler(event)
		}
	}
}

// package kube

// import (
// 	"context"
// 	"fmt"
// 	"k8s.io/client-go/tools/clientcmd"
// 	"k8s.io/client-go/kubernetes"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"os"
// )

// func GetKubeClient() *kubernetes.Clientset {
// 	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
// 	if err != nil {
// 		fmt.Println("Error loading kubeconfig:", err)
// 		os.Exit(1)
// 	}
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		fmt.Println("Error creating clientset:", err)
// 		os.Exit(1)
// 	}
// 	return clientset
// }

// func WatchEvents(namespace string, includePast bool, eventHandler func(event *corev1.Event)) {
// 	clientset := GetKubeClient()
// 	ctx := context.TODO()

// 	// Determine starting resource version and optionally emit existing events
// 	var resourceVersion string
// 	if includePast {
// 		// Fetch and emit all existing events
// 		evList, err := clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
// 		if err != nil {
// 			fmt.Println("Error listing events:", err)
// 			return
// 		}
// 		for i := range evList.Items {
// 			eventHandler(&evList.Items[i])
// 		}
// 		resourceVersion = evList.ResourceVersion
// 	} else {
// 		// Start watching only new events
// 		resourceVersion = "0"
// 	}

// 	watcher, err := clientset.CoreV1().Events(namespace).Watch(ctx, metav1.ListOptions{
// 		ResourceVersion: resourceVersion,
// 	})
// 	if err != nil {
// 		fmt.Println("Failed to watch events:", err)
// 		return
// 	}
// 	ch := watcher.ResultChan()

// 	for evt := range ch {
// 		if event, ok := evt.Object.(*corev1.Event); ok {
// 			eventHandler(event)
// 		}
// 	}
// }
