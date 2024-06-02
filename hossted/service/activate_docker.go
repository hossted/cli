/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package service

type VmInfo struct {
	ClusterUUID string `json:"clusterUUID"`
	EmailID     string `json:"email_id,omitempty"`
}

// AppAPIInfo contains basic information about the application API.
type ContainerStatus struct {
	Name       string `json:"name,omitempty"`
	AppUUID    string `json:"appUUID,omitempty"`
	Status     string `json:"status,omitempty"`
	AppVersion string `json:"appVersion,omitempty"`
	Image      string `json:"image,omitempty"`
}

type AppInfo struct {
	VmInfo          VmInfo            `json:"vm_info"`
	ContainerStatus []ContainerStatus `json:"container_status"`
}

func ActivateCompose() error {
	// emailsID, err := getEmail()
	// if err != nil {
	// 	return err
	// }
	// getResponse from reading file in .hossted/email@id.json
	resp, err := getResponse()
	if err != nil {
		return err
	}
	// validate auth token

	err = validateToken(resp)
	if err != nil {
		return err
	}
	// handle usecases for orgID selection
	_, err = useCases(resp)
	if err != nil {
		return err
	}
	// VM REGISTRATION
	// reconiler( )
	// 	// perform cluster registeration
	// - create UUID if not exists

	// app registeration
	//  -

	// )
	// ADD MONITORING

	err = Reconcile()
	if err != nil {
		return err
	}
	// SET CLUSTER UUID
	// UUID, err :=
	// if err!=nil {

	// }
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
// 	info := AppInfo{
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
// 	var appData AppInfo
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
