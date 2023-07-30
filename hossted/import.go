package hossted

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
)

// hossted ping - send docker ,sbom and security infor to hossted API
func Import(env string, authorization string,  kluster KCluster ) error {
	// Get the path to the kubeconfig file
	kubeconfig := getKubeconfigPath()

	klusterJson, _ := json.Marshal(kluster)
	fmt.Println("k8s Cluster json:", string( klusterJson ))


	// Create the Kubernetes client.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	// Print current context


	// List the namespaces.
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//ActiveContext := `hardcoded`

	//TODO get current Context
	// Print the namespaces.
	listOptions := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=Helm",
	}

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
/*
func generateKlusterJson(ActiveContext string ,  namespaces metav1.Namespaces ){
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
}
*/

func ImportRequest(env, authorization, kluster, klusterJson string) (importResponse, error) {

	var response importResponse

	// Construct param map for input params
	var data map[string]string
	err := json.Unmarshal([]byte(klusterJson), &data)
	println(data)
	if err != nil {
		panic(err)
	}

	params := make(map[string]string)
	for k, v := range data {
		params[k] = v
	}
	req := HosstedRequest{
		// Endpoint env needs to replace in runtime for url parse to work. Otherwise runtime error.
		//EndPoint:     "https://api.__ENV__hossted.com/v1/instances/dockers",
		EndPoint:     "https://api.hossted.com/v1/instances/registry", //"https://api.stage.hossted.com/v1/instances/registry", // "https://api.dev.hossted.com/v1/instances/registry", //,
		Environment:  env,
		Params:       params,
		BearToken:    "Basic " + authorization,
		SessionToken: "",
		TypeRequest:  "POST",
	}

	resp, err := req.SendRequest()
	if err != nil {
		//fmt.Println("\033[31m", "Error:", "\033[0m")
		fmt.Printf("\033[0;31mError: \033[0m")
		fmt.Printf("%s\n\033[0m", err)
		return response, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return response, fmt.Errorf("Failed to parse JSON. %w", err)
	}

	return response, nil
}

func getKubeconfigPath() string {
	home := homedir.HomeDir()
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	return kubeconfig
}