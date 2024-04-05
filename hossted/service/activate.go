package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"k8s.io/client-go/tools/clientcmd"
)

// Response represents the structure of the JSON response.
type response struct {
	Success bool                `json:"success"`
	OrgIDs  []map[string]string `json:"org_ids"`
	Token   string              `json:"token"`
	Message string              `json:"message"`
}

// ActivateK8s imports Kubernetes clusters.
func ActivateK8s() error {
	// Prompt user for email
	prompt := promptui.Prompt{
		Label: "Enter your email:",
	}

	emailID, err := prompt.Run()
	if err != nil {
		return err
	}

	//read file
	homeDir, err := os.UserHomeDir()

	folderPath := filepath.Join(homeDir, ".hossted")
	if err != nil {
		return err
	}

	fileData, err := os.ReadFile(folderPath + "/" + emailID + ".json")
	if err != nil {
		return fmt.Errorf("User not registered, Please run hossted login to register")
	}

	var file response
	err = json.Unmarshal(fileData, &file)
	if err != nil {
		return err
	}

	// Retrieve available Kubernetes contexts

	if file.Success {
		if len(file.OrgIDs) == 0 {
			fmt.Println("We have just sent the confirmation link to", emailID, ". Once you confirm it, you'll be able to continue the activation.")
		} else if len(file.OrgIDs) == 1 {
			for orgID, email := range file.OrgIDs[0] {
				prompt := promptui.Select{
					Label: fmt.Sprintf("Are you sure you want to register this cluster with org_name %s", email),
					Items: []string{"Yes", "No"},
				}
				_, value, err := prompt.Run()
				if err != nil {
					return err
				}
				if value == "Yes" {
					fmt.Println("Your orgID is ", orgID)
				} else {
					return nil
				}
			}
		} else if len(file.OrgIDs) > 1 {
			// Handle cases where len(resp.OrgIDs) > 1
			fmt.Println("There are multiple organization IDs. Handling multiple org IDs logic here.")
		}
	} else {
		return fmt.Errorf("Cluster registration failed to hossted platform")
	}

	contexts := getK8sContext()

	// // Prompt user to select Kubernetes context
	promptK8s := promptui.Select{
		Label: "Select your Kubernetes context:",
		Items: contexts,
	}
	_, clusterName, err := promptK8s.Run()
	if err != nil {
		return err
	}

	fmt.Println("Your cluster name is ", clusterName)
	// fmt.Println("Email:", emailID, "| ClusterName:", clusterName, "| ClusterUUID:", clusterUUID, "| Response:", resp)

	operator := promptui.Select{
		Label: fmt.Sprintf("Do you wish to install the operator in %s", clusterName),
		Items: []string{"Yes", "No"},
	}
	_, value, err := operator.Run()
	if err != nil {
		return err
	}

	if value == "Yes" {
		monitoring := promptui.Select{
			Label: fmt.Sprintf("Do you wish to enable monitoring in operator"),
			Items: []string{"Yes", "No"},
		}
		_, monitoringEnable, err := monitoring.Run()
		if err != nil {
			return err
		}

		if monitoringEnable == "Yes" {
			fmt.Println("Enabled Monitoring ", monitoringEnable)
		}

		cve := promptui.Select{
			Label: fmt.Sprintf("Do you wish to enable cve scan in operator"),
			Items: []string{"Yes", "No"},
		}
		_, cveEnable, err := cve.Run()
		if err != nil {
			return err
		}

		if cveEnable == "Yes" {
			fmt.Println("Enabled CVE Scanning ", cveEnable)
		}

	}
	return nil
}

// getK8sContext retrieves Kubernetes contexts from kubeconfig.
func getK8sContext() []string {
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
	for i := range currentContext {
		contexts = append(contexts, i)
	}
	return contexts
}

// registerClusterUUID registers cluster UUID with provided email ID and cluster name.
func registerClusterUUID(clusterUUID, clusterName, emailID string) (response, error) {
	clusterUUIDRegPath := os.Getenv("HOSSTED_API_URL") + "/clusters/" + clusterUUID + "/register"

	type clusterUUIDBody struct {
		Email       string `json:"email"`
		ReqType     string `json:"type"`
		ClusterName string `json:"cluster_name"`
	}

	clusterUUIDBodyReq := clusterUUIDBody{
		Email:       emailID,
		ReqType:     "k8s",
		ClusterName: clusterName,
	}

	body, err := json.Marshal(clusterUUIDBodyReq)
	if err != nil {
		return response{}, err
	}

	res, err := httpRequest(body, http.MethodPost, clusterUUIDRegPath)
	if err != nil {
		return response{}, err
	}

	responseBody, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(responseBody))
	var resp response
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return response{}, err
	}

	return resp, nil
}

// registerClusterUUID registers cluster UUID with provided email ID and cluster name.
func patchCluster(clusterUUID, clusterName, emailID, orgID string) (response, error) {
	clusterUUIDRegPath := os.Getenv("HOSSTED_API_URL") + "/clusters/" + clusterUUID + "/register"

	type clusterUUIDBody struct {
		Email       string `json:"email"`
		ReqType     string `json:"type"`
		OrgID       string `json:"org_id"`
		ClusterName string `json:"cluster_name"`
	}

	clusterUUIDBodyReq := clusterUUIDBody{
		Email:       emailID,
		ReqType:     "k8s",
		ClusterName: clusterName,
	}

	body, err := json.Marshal(clusterUUIDBodyReq)
	if err != nil {
		return response{}, err
	}

	res, err := httpRequest(body, http.MethodPatch, clusterUUIDRegPath)
	if err != nil {
		return response{}, err
	}

	responseBody, err := ioutil.ReadAll(res.Body)

	var resp response
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return response{}, err
	}

	return resp, nil
}

// httpRequest performs an HTTP POST request.
func httpRequest(body []byte, reqType string, url string) (*http.Response, error) {
	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+os.Getenv("HOSSTED_AUTH_TOKEN"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}
