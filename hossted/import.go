package hossted

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
)

// hossted ping - send docker ,sbom and security infor to hossted API
func Import(env string) error {
	// Get the path to the kubeconfig file
	kubeconfig := getKubeconfigPath()

	// Create the Kubernetes client.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// List the namespaces.
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	listOptions := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=Helm",
	}

	// Print the namespaces.
	for _, namespace := range namespaces.Items {
		fmt.Printf("Helm releases in namespace %s:\n", namespace.Name)
		releases, err := clientset.AppsV1().Deployments(namespace.Name).List(context.TODO(), listOptions)
		if err != nil {
			fmt.Printf("Error listing Helm releases: %v\n", err)
			os.Exit(1)
		}
		for _, release := range releases.Items {
			fmt.Println(release.Name)
		}
	}
	return err
}



func getKubeconfigPath() string {
	home := homedir.HomeDir()
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	return kubeconfig
}