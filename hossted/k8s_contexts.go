package hossted

import (
	"fmt"
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

func GetK8sContext() []string {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}

	// Load kubeconfig file
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Get current context
	currentContext := config.Contexts
	var contexts []string
	for i, _ := range currentContext {
		contexts = append(contexts, i)
	}
	return contexts
}
