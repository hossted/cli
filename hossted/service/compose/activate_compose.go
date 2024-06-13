package compose

import (
	"github.com/hossted/cli/hossted/service/common"
)

type OsInfo struct {
	OsUUID               string `yaml:"osUUID"`
	EmailID              string `yaml:"emailId,omitempty"`
	ClusterRegisteration bool   `yaml:"clusterRegisteration,omitempty"`
}

type AppRequest struct {
	AppAPIInfo AppAPIInfo `json:"app_api_info"`
	AppInfo    AppInfo    `json:"app_info"`
}

// AppAPIInfo contains basic information about the application API.
type AppAPIInfo struct {
	AppUUID string `json:"app_uuid,omitempty"`
	OsUUID  string `json:"os_uuid,omitempty"`
	EmailID string `json:"email_id,omitempty"`
}

type AppInfo struct {
	ComposeFile    string           `json:"compose_file,omitempty"`
	DockerInstance []dockerInstance `json:"container_info,omitempty"`
}

type ContainerInfo struct {
	Name       string      `json:"names,omitempty"`
	Status     string      `json:"status,omitempty"`
	AppVersion string      `json:"app_version,omitempty"`
	Image      string      `json:"image,omitempty"`
	ImageID    string      `json:"image_id,omitempty"`
	Ports      interface{} `json:"ports,omitempty"`
	Size       int64       `json:"size,omitempty"`
	Networks   string      `json:"networks,omitempty"`
	Mounts     string      `json:"mounts,omitempty"`
	DockerID   string      `json:"docker_id,omitempty"`
}

func ActivateCompose() error {

	emailID, err := common.GetEmail()
	if err != nil {
		return err
	}

	resp, err := common.GetLoginResponse()
	if err != nil {
		return err
	}

	// validate auth token
	err = common.ValidateToken(resp)
	if err != nil {
		return err
	}
	// handle usecases for orgID selection
	orgID, err := common.UseCases(resp)
	if err != nil {
		return err
	}

	err = reconcileCompose(orgID, emailID, resp.Token)
	if err != nil {
		return err
	}

	return nil

}

// provide prompt to enable monitoring and vulnerability scan

//

// func setUUID(email string) error {

// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		fmt.Printf("Error getting home directory: %s\n", err)
// 		return err
// 	}
// 	// Construct the full path to compose-status.json
// 	filePath := filepath.Join(homeDir, ".hossted", "compose-status.json")
// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		err = updateUUID(filePath, email)
// 		if err != nil {
// 			return err
// 		}
// 	} else if err != nil {
// 		return err
// 	} else {
// 		prompt := promptui.Select{
// 			Label: "Config file already exists do want to override it?",
// 			Items: []string{"Yes", "No"},
// 		}
// 		_, result, err := prompt.Run()
// 		if err != nil {
// 			fmt.Printf("Prompt failed %v\n", err)
// 			return err
// 		}
// 		if result == "No" {
// 			fmt.Println("Exiting.")
// 			os.Exit(0)
// 		}
// 		err = updateUUID(filePath, email)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func updateUUID(filepath, email string) error {
// 	clusterUUID := "D-" + uuid.NewString()
// 	fmt.Println("Generating Cluster UUID", clusterUUID)
// 	info := AppsInfo{
// 		ClusterInfo: ClusterInfo{
// 			ClusterUUID: clusterUUID,
// 			EmailID:     email,
// 		},
// 		ComposeStatus: []ComposeStatus{},
// 	}
// 	jsonData, err := json.MarshalIndent(info, "", "    ")
// 	if err != nil {
// 		fmt.Printf("Error marshaling JSON: %s\n", err)
// 		return err
// 	}
// 	// Write the JSON data to the file
// 	err = ioutil.WriteFile(filepath, jsonData, 0644)
// 	if err != nil {
// 		fmt.Printf("Error writing JSON to file: %s\n", err)
// 		return err
// 	}
// 	return nil
// }

// func checkUUID(filePath string) (string, error)  {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			fmt.Printf("File %s does not exist.\n", filePath)
// 		} else {
// 			fmt.Printf("Error opening file: %s\n", err)
// 		}
// 		return "", err
// 	}
// 	defer file.Close()

// 	data, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		fmt.Printf("Error reading file: %s\n", err)
// 		return "", err
// 	}
// 	var appData AppsInfo
// 	err = json.Unmarshal(data, &appData)
// 	if err != nil {
// 		fmt.Printf("Error unmarshaling JSON: %s\n", err)
// 		return "", err
// 	}
// 	fmt.Println("UUID IS", appData.ClusterInfo.ClusterUUID)
// 	return appData.ClusterInfo.ClusterUUID, nil
// }

//GET CLUSTER ID

// func getClusterUUID() (string, error) {
// 	var clusterUUID string
// 	var err error
// 	yellow := color.New(color.FgYellow).SprintFunc()

// 	// Retry for 120 seconds
// 	for i := 0; i < 120; i++ {
// 		clusterUUID, err = getClusterUUIDFromDocker()
// 		if err == nil {
// 			return clusterUUID, nil
// 		}
// 		fmt.Println(yellow("Info:"), "Waiting for Hossted Operator to get into running state")
// 		time.Sleep(4 * time.Second) // Wait for 1 second before retrying
// 	}

// 	return "", fmt.Errorf("Failed to get ClusterUUID after 120 seconds: %v", err)
// }

// func getClusterUUIDFromDocker() (string, error) {
// 	cs := getKubeClient()
// 	hp, err := cs.Resource(hpGVK).Get(context.TODO(), "hossted-operator-cr", metav1.GetOptions{})
// 	if err != nil {
// 		return "", err
// 	}
// 	clusterUUID, _, err := unstructured.NestedString(hp.Object, "status", "clusterUUID")
// 	if err != nil || clusterUUID == "" {
// 		return "", fmt.Errorf("ClusterUUID is nil, func errored")
// 	}
// 	return clusterUUID, nil
// }
