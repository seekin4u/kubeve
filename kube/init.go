package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Kinit sets up the Kubernetes client and returns the namespace, raw kubeconfig, clientset, and namespace list.
func Kinit(overrideNamespace string) (string, clientcmdapi.Config, *kubernetes.Clientset, []string, error) {
	// Load kubeconfig rules and overrides
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)

	// Determine namespace: override or default
	ns := overrideNamespace
	if ns == "" {
		var err error
		ns, _, err = clientConfig.Namespace()
		if err != nil {
			ns = metav1.NamespaceDefault
		}
	}

	// Load raw config
	rawCfg, err := clientConfig.RawConfig()
	if err != nil {
		return "", clientcmdapi.Config{}, nil, nil, err
	}

	// Build REST config & clientset
	restCfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return "", rawCfg, nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return "", rawCfg, nil, nil, err
	}

	// Retrieve namespace list
	var nsList []string
	nsItems, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		for _, item := range nsItems.Items {
			nsList = append(nsList, item.Name)
		}
	}

	return ns, rawCfg, clientset, nsList, nil
}
