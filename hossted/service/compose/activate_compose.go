package compose

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hossted/cli/hossted/service/common"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

type OsInfo struct {
	OsUUID               string `yaml:"osUUID,omitempty"`
	EmailID              string `yaml:"emailId,omitempty"`
	ClusterRegisteration bool   `yaml:"clusterRegisteration,omitempty"`
	OrgID                string `yaml:"orgID,omitempty"`
	Token                string `yaml:"token,omitempty"`
	ProjectName          string `yaml:"projectName,omitempty"`
	HosstedApiUrl        string `yaml:"hosstedAPIUrl,omitempty"`
	MimirUsername        string `yaml:"MIMIR_USERNAME,omitempty"`
	MimirPassword        string `yaml:"MIMIR_PASSWORD,omitempty"`
	MimirUrl             string `yaml:"MIMIR_URL,omitempty"`
	LokiUsername         string `yaml:"LOKI_USERNAME,omitempty"`
	LokiPassword         string `yaml:"LOKI_PASSWORD,omitempty"`
	LokiUrl              string `yaml:"LOKI_URL,omitempty"`
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

func ActivateCompose(composeFilePath string, develMode bool) error {

	osFilePath, err := getComposeFilePath("compose.yaml")
	if err != nil {
		return err
	}

	tr, err := common.GetTokenResp()
	if err != nil {
		return err
	}

	orgs, err := common.GetOrgs(tr.AccessToken)
	if err != nil {
		return err
	}

	orgID, err := common.OrgUseCases(orgs)
	if err != nil {
		return err
	}

	projectName, err := getProjectName()
	if err != nil {
		return err
	}

	osInfo := OsInfo{
		OrgID:         orgID,
		Token:         tr.AccessToken,
		ProjectName:   projectName,
		HosstedApiUrl: common.HOSSTED_API_URL,
		MimirUsername: common.MIMIR_USERNAME,
		MimirPassword: common.MIMIR_PASSWORD,
		MimirUrl:      common.MIMIR_URL,
		LokiUsername:  common.LOKI_USERNAME,
		LokiPassword:  common.LOKI_PASSWORD,
		LokiUrl:       common.LOKI_URL,
	}

	// Override values in development mode
	if develMode {

		if devUrl := common.HOSSTED_DEV_API_URL; devUrl != "" {
			osInfo.HosstedApiUrl = devUrl
		}
		if devUrl := common.MIMIR_DEV_URL; devUrl != "" {
			osInfo.MimirUrl = devUrl
		}
		if devUrl := common.LOKI_DEV_URL; devUrl != "" {
			osInfo.LokiUrl = devUrl
		}
	}

	fmt.Println(osFilePath)
	osData, err := setClusterInfo(osInfo, osFilePath)
	if err != nil {
		return err
	}

	enableMonitoring, err := askPromptsToInstall()
	if err != nil {
		return err
	}

	err = ReconcileCompose(osData, enableMonitoring)
	if err != nil {
		return err
	}

	return nil

}

func askPromptsToInstall() (string, error) {
	green := color.New(color.FgGreen).SprintFunc()

	monitoringEnabled := "false"

	//------------------------------ Monitoring ----------------------------------
	monitoring := promptui.Select{
		Label: "Do you wish to enable monitoring in hossted platform",
		Items: []string{"Yes", "No"},
	}
	_, monitoringEnable, err := monitoring.Run()
	if err != nil {
		return monitoringEnable, err
	}

	if monitoringEnable == "Yes" {
		fmt.Println("Enabled Monitoring :", green(monitoringEnable))
		monitoringEnabled = "true"
		AddComposeFile()
	}

	return monitoringEnabled, nil
}

// getProjectName takes a file path and returns the final directory name
func GetProjectName(filePath string) string {
	// Clean the path to handle any extraneous characters
	cleanPath := filepath.Clean(filePath)

	// Get the final directory name
	return filepath.Base(cleanPath)
}

func AddComposeFile() {
	files := map[string]string{
		"https://raw.githubusercontent.com/hossted/cli/main/compose/monitoring/config.river":        "config.river",
		"https://raw.githubusercontent.com/hossted/cli/main/compose/monitoring/docker-compose.yaml": "docker-compose.yaml",
	}

	// Define the base directory where the files will be saved
	baseDir := filepath.Join(os.Getenv("HOME"), ".hossted/compose/monitoring")

	// Ensure the base directory exists
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Download and save each file
	for url, fileName := range files {
		savePath := filepath.Join(baseDir, fileName)
		err := DownloadFile(url, savePath)
		if err != nil {
			log.Fatalf("Failed to download %s: %v", url, err)
		}
		log.Printf("File successfully downloaded to %s", savePath)
	}

}

// DownloadFile downloads a file from a URL and saves it to a local path.
func DownloadFile(url, filePath string) error {
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func GetClusterInfo() (OsInfo, error) {
	var osInfo OsInfo
	//read file
	homeDir, err := os.UserHomeDir()

	folderPath := filepath.Join(homeDir, ".hossted")
	if err != nil {
		return osInfo, err
	}

	fileData, err := os.ReadFile(folderPath + "/" + "compose.yaml")
	if err != nil {
		return osInfo, fmt.Errorf("unable to read %s file", folderPath+"/compose.yaml")
	}

	err = yaml.Unmarshal(fileData, &osInfo)
	if err != nil {
		return osInfo, err
	}

	return osInfo, nil
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
