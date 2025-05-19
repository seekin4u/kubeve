package kube

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func WatchEvents(namespace string, eventHandler func(event *corev1.Event)) {
	_, _, clientset, _, err := Kinit(namespace)
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
