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
	AppUUID              string `yaml:"appUUID,omitempty"`
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

	osFilePath, err := getHosstedComposeFilePath("compose.yaml")
	if err != nil {
		return err
	}

	tr, err := common.GetTokenResp()
	if err != nil {
		return err
	}

	orgs, _, err := common.GetOrgs(tr.AccessToken)
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

	fmt.Println("ProjectName: ", projectName)

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

	osUUID, err := getOSUUID()
	if err != nil {
		return err
	}

	appUUID, err := getMarketplaceAppUUID(osInfo.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to get appUUID, err: %v", appUUID)
	}

	osInfo.OsUUID = osUUID
	osInfo.AppUUID = appUUID

	yamlData, err := yaml.Marshal(osInfo)
	if err != nil {
		return fmt.Errorf("error in YAML marshaling: %s", err)
	}

	err = writeFile(osFilePath, yamlData)
	if err != nil {
		return err
	}

	enableMonitoring, err := askPromptsToInstall(develMode)
	if err != nil {
		return err
	}

	err = ReconcileCompose(osInfo, enableMonitoring)
	if err != nil {
		return err
	}

	return nil

}

func askPromptsToInstall(develMode bool) (string, error) {
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
		AddComposeFile(develMode)
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

func GetMode(develMode bool) string {
	if develMode {
		return "dev"
	} else {
		return "main"
	}
}

func AddComposeFile(develMode bool) {
	// Set your environment to "prod" or "dev"
	branch := GetMode(develMode)
	files := map[string]string{
		fmt.Sprintf("https://raw.githubusercontent.com/hossted/cli/%s/compose/monitoring/config.river", branch):        "config.river",
		fmt.Sprintf("https://raw.githubusercontent.com/hossted/cli/%s/compose/monitoring/docker-compose.yaml", branch): "docker-compose.yaml",
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
