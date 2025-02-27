package compose

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/hossted/cli/hossted/service/common"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"gopkg.in/yaml.v2"
)

func ReconcileCompose(osInfo OsInfo, enableMonitoring string) error {

	appFilePath, err := getHosstedComposeFilePath("compose-request.json")
	if err != nil {
		return err
	}

	list, err := listAllContainers(osInfo.ProjectName)
	if err != nil {
		return err
	}

	isComposeStateChange, err := writeComposeRequest2File(
		appFilePath,
		list,
		osInfo.OsUUID,
		osInfo.AppUUID,
		osInfo.EmailID,
		enableMonitoring,
		osInfo.ProjectName)
	if err != nil {
		return err
	}

	if isComposeStateChange {
		// send compose info
		err = sendComposeInfo(appFilePath, osInfo)
		if err != nil {
			return err
		}
	}
	return nil

}

func getHosstedComposeFilePath(filename string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %s\n", err)
		return "", err
	}

	// Construct the full path to file
	filePath := filepath.Join(homeDir, ".hossted", filename)
	return filePath, nil
}

func readFile(filePath string) ([]byte, error) {
	var data []byte
	file, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	defer file.Close()

	data, err = io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return data, err
	}
	return data, nil
}

func writeFile(filePath string, data []byte) error {
	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("error creating %s file: %s\n", filePath, err)
		return err
	}
	defer file.Close()

	// Write the data to the file
	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("error writing %s file: %s\n", filePath, err)
		return err
	}

	return nil

}

func getOSUUID() (string, error) {

	ok, _ := isMarketplaceVM()
	if ok {
		osUUID, _ := getMarketplaceOSUUID()
		fmt.Println("Found existing osUUID for marketplace VM: ", osUUID)
		return osUUID, nil
	} else {
		osUUID := "D-" + uuid.NewString()
		fmt.Println("Generating osUUID for VM: ", osUUID)
		return osUUID, nil
	}
}

func checkUUID(osFilePath string) (string, error) {
	var osData OsInfo

	data, err := readFile(osFilePath)
	if err != nil {
		//return "", err
	}

	err = yaml.Unmarshal(data, &osData)
	if err != nil {
		fmt.Printf("Error unmarshaling Yaml: %s\n", err)
		return "", err
	}

	return osData.OsUUID, nil
}

func writeComposeRequest2File(
	appFilePath string,
	containerList []types.Container,
	osUUID,
	appUUID,
	email,
	enableMonitoring,
	projectName string) (bool, error) {
	// prepare appsInfo with updated container status
	appRequest, isComposeStateChange, err := prepareComposeRequest(
		appFilePath,
		containerList,
		osUUID,
		appUUID,
		email,
		enableMonitoring,
		projectName,
	)

	if err != nil {
		return isComposeStateChange, err
	}

	jsonData, err := json.MarshalIndent(appRequest, "", "    ")
	if err != nil {
		fmt.Printf("error marshaling JSON: %s\n", err)
		return isComposeStateChange, err
	}

	// write compose status file
	err = writeFile(appFilePath, jsonData)
	if err != nil {
		return isComposeStateChange, err
	}

	return isComposeStateChange, nil
}

type OptionsState struct {
	Monitoring bool `json:"monitoring,omitempty"`
	Logging    bool `json:"logging,omitempty"`
	CVEScan    bool `json:"cve_scan,omitempty"`
}

type URLInfo struct {
	URL      string `json:"url"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type AccessInfo struct {
	URLs []URLInfo `json:"urls"`
}

type Request struct {
	UUID         string       `json:"uuid"`
	OsUUID       string       `json:"osuuid"`
	OrgID        string       `json:"org_id"`
	Email        string       `json:"email,omitempty"`
	Type         string       `json:"type,omitempty"`
	Product      string       `json:"product,omitempty"`
	CPUNum       string       `json:"cpunum,omitempty"`
	Memory       string       `json:"memory,omitempty"`
	OptionsState OptionsState `json:"options_state,omitempty"`
	ComposeFile  string       `json:"compose_file,omitempty"`
	AccessInfo   AccessInfo   `json:"access_info,omitempty"`
}

type dockerInstance struct {
	DockerID   string      `json:"docker_id"`
	InstanceID string      `json:"instance_id"`
	ImageID    string      `json:"image_id"`
	Ports      interface{} `json:"ports"`
	Status     string      `json:"status"`
	Size       string      `json:"size"`
	Names      []string    `json:"names"`
	Mounts     interface{} `json:"mounts"`
	Networks   string      `json:"networks"`
	Tag        string      `json:"tag"`
	Image      string      `json:"image"`
	CreatedAt  string      `json:"created_at"`
}

func sendComposeInfo(appFilePath string, osInfo OsInfo) error {
	hosstedAPIUrl := osInfo.HosstedApiUrl
	orgID := osInfo.OrgID
	token := osInfo.Token

	composeInfo, err := readFile(appFilePath)
	if err != nil {
		return err
	}

	projectName, err := getProjectName()
	if err != nil {
		return err
	}

	accessInfo := GetAccessInfo("/opt/" + projectName + "/.env")

	var data map[string]AppRequest
	if err = json.Unmarshal(composeInfo, &data); err != nil {
		return err
	}

	cpu, err := getCPUUsage()
	if err != nil {
		return err
	}
	mem, err := getMemoryUsage()
	if err != nil {
		return err
	}

	isMarketplace, err := isMarketplaceVM()
	if err != nil {
		return fmt.Errorf("isMarketplaceVM func errored: %v", err)
	}

	if isMarketplace {
		if err = submitPatchRequest(osInfo, data, *accessInfo, cpu, mem); err != nil {
			return fmt.Errorf("error in patch request: %v", err)
		}
		fmt.Println("Successfully submitted PATCH request if marketplace VM")
	} else {
		if err = registerApplications(data, osInfo, *accessInfo, cpu, mem, orgID, token, hosstedAPIUrl+"/compose/hosts"); err != nil {
			return err
		}
	}

	return registerDockerInstances(data, osInfo, token, hosstedAPIUrl+"/compose/containers", isMarketplace)
}

// registerApplications registers all applications with the specified API URL.
func registerApplications(data map[string]AppRequest, osInfo OsInfo, accessInfo AccessInfo, cpu, mem, orgID, token, composeUrl string) error {
	for appName, compose := range data {
		newReq := Request{
			UUID:       compose.AppAPIInfo.AppUUID,
			OsUUID:     compose.AppAPIInfo.OsUUID,
			Email:      compose.AppAPIInfo.EmailID,
			AccessInfo: accessInfo,
			OrgID:      orgID,
			Type:       "compose",
			Product:    appName,
			CPUNum:     cpu,
			Memory:     mem,
			OptionsState: OptionsState{
				Monitoring: true,
				Logging:    true,
				CVEScan:    true,
			},
			ComposeFile: compose.AppInfo.ComposeFile,
		}

		if err := SendRequest("POST", composeUrl, token, newReq); err != nil {
			return err
		}
		fmt.Printf("Successfully registered app [%s] with appID [%s]\n", appName, compose.AppAPIInfo.AppUUID)
	}
	return nil
}

// registerDockerInstances registers all docker instances based on whether it is a marketplace VM.
func registerDockerInstances(data map[string]AppRequest, osInfo OsInfo, token, containersUrl string, isMarketplace bool) error {
	for appName, info := range data {
		for _, ci := range info.AppInfo.DockerInstance {
			instanceID := osInfo.AppUUID
			if isMarketplace {
				instanceID = osInfo.OsUUID
			}
			newDI := dockerInstance{
				DockerID:   ci.DockerID,
				InstanceID: instanceID,
				ImageID:    ci.ImageID,
				Ports:      ci.Ports,
				Status:     ci.Status,
				Size:       ci.Size,
				Names:      ci.Names,
				Networks:   ci.Networks,
				Image:      ci.Image,
				Mounts:     ci.Mounts,
				Tag:        ci.Tag,
				CreatedAt:  ci.CreatedAt,
			}

			if err := SendRequest("POST", containersUrl, token, newDI); err != nil {
				return err
			}

			fmt.Printf("Successfully registered docker info, appName:[%s], dockerName:[%s], appID:[%s]\n", appName, ci.Names, info.AppAPIInfo.AppUUID)
		}
	}
	return nil
}

// submitPatchRequest sends a PATCH request with VM info for marketplace setups.
func submitPatchRequest(osInfo OsInfo, compose map[string]AppRequest, accessInfo AccessInfo, cpu, mem string) error {
	composeUrl := osInfo.HosstedApiUrl + "/compose/hosts/" + osInfo.OsUUID

	var composeFile, applicationName string
	for appName, appReq := range compose {
		composeFile = appReq.AppInfo.ComposeFile
		applicationName = appName
	}

	newReq := Request{
		UUID:       osInfo.AppUUID,
		OsUUID:     osInfo.OsUUID,
		AccessInfo: accessInfo,
		OrgID:      osInfo.OrgID,
		Type:       "vm",
		Product:    applicationName,
		CPUNum:     cpu,
		Memory:     mem,
		OptionsState: OptionsState{
			Monitoring: true,
			Logging:    true,
			CVEScan:    false,
		},
		ComposeFile: composeFile,
	}

	return SendRequest(http.MethodPatch, composeUrl, osInfo.Token, newReq)
}

// sendRequest handles HTTP requests for a given method, URL, token, and request body.
func SendRequest(method, url, token string, reqBody interface{}) error {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	err = common.HttpRequest(method, url, token, body)
	if err != nil {
		return fmt.Errorf("error in %s request to %s: %v", method, url, err)
	}

	return nil
}

func listAllContainers(projectName string) ([]types.Container, error) {
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	// Define filters to list only running containers
	filter := filters.NewArgs()
	//filter.Add("status", "running")

	filter.Add("label", "com.docker.compose.project="+projectName)

	filter.Add("label", "com.docker.compose.project")
	filter.Add("label", "com.docker.compose.config-hash")

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filter})

	if err != nil {
		return nil, err
	}
	return containers, nil
}

func listProjectContainers(projectName string) ([]types.Container, error) {
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	// Define context
	ctx := context.Background()

	// Define filters to list only running containers
	filter := filters.NewArgs()
	filter.Add("label", "com.docker.compose.project="+projectName)
	filter.Add("label", "com.docker.compose.config-hash")
	filter.Add("label", "com.docker.compose.project.config_files")

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})

	if err != nil {
		return nil, err
	}
	return containers, nil
}

func getUniqueComposeProjects(containerList []types.Container) (map[string]bool, error) {
	uniqueProjects := make(map[string]bool)

	for _, container := range containerList {
		if project, ok := container.Labels["com.docker.compose.project"]; ok {
			uniqueProjects[project] = true
		}
	}

	return uniqueProjects, nil
}

func prepareComposeRequest(
	appFilePath string,
	containerList []types.Container,
	osUUID,
	appUUID,
	email,
	enableMonitoring,
	projectName string) (map[string]AppRequest, bool, error) {
	var appsData map[string]AppRequest
	isComposeStateChange := false

	if _, err := os.Stat(appFilePath); os.IsNotExist(err) {
		appsData = make(map[string]AppRequest)
	} else {
		data, err := readFile(appFilePath)
		if err != nil {
			fmt.Printf("Error in reading %s file: %s\n", appFilePath, err)
			return appsData, isComposeStateChange, err
		}
		err = json.Unmarshal(data, &appsData)
		if err != nil {
			fmt.Printf("Error in JSON unmarshaling %s file: %s\n", appFilePath, err)
			return appsData, isComposeStateChange, err
		}
	}

	uniqueProjects, err := getUniqueComposeProjects(containerList)
	if err != nil {
		return appsData, isComposeStateChange, err
	}

	// Create a slice of existing apps
	existingApps := make(map[string]bool)
	for app := range appsData {
		existingApps[app] = true
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return appsData, isComposeStateChange, err
	}

	for project := range uniqueProjects {
		var composeFileContent string
		list, err := listProjectContainers(project)
		if err != nil {
			return appsData, isComposeStateChange, err
		}

		// Check if project is present in the existingApps
		if found := existingApps[project]; found {
			// create only AppInfo(compose request)
			var newComposeRequest AppInfo
			var dockerInstances []dockerInstance

			for i, container := range list {
				if i == 0 {
					// Extract compose file content
					if _, ok := container.Labels["com.docker.compose.project.config_files"]; ok {
						composeDir := container.Labels["com.docker.compose.project.working_dir"]
						data, err := readFile(composeDir + "/docker-compose.yml")
						if err != nil {
							fmt.Printf("Error in reading compose here file %s: %s\n", composeDir, err)
						}
						composeFileContent = string(data)
					}
				}

				meta, err := cli.ContainerInspect(context.Background(), container.ID)
				if err != nil {
					return appsData, isComposeStateChange, err
				}

				imagemeta, _, err := cli.ImageInspectWithRaw(context.Background(), container.ImageID)
				if err != nil {
					return appsData, isComposeStateChange, err
				}

				containerInfo := dockerInstance{
					Names:     container.Names,
					Status:    container.State,
					Image:     container.Image,
					ImageID:   container.ImageID,
					Ports:     container.Ports,
					Size:      convertInt64ToString(imagemeta.Size),
					Networks:  container.HostConfig.NetworkMode,
					DockerID:  container.ID,
					Mounts:    container.Mounts,
					Tag:       container.Image,
					CreatedAt: meta.Created,
				}

				dockerInstances = append(dockerInstances, containerInfo)
			}

			newComposeRequest.DockerInstance = dockerInstances
			newComposeRequest.ComposeFile = composeFileContent

			// get previous compose request
			prevComposeRequest := appsData[project].AppInfo

			//compare newComposeRequest and prevComposeRequest
			res := reflect.DeepEqual(prevComposeRequest.ComposeFile, newComposeRequest.ComposeFile)
			if !res {
				appsData[project] = AppRequest{
					AppAPIInfo: appsData[project].AppAPIInfo,
					AppInfo:    newComposeRequest,
				}
				isComposeStateChange = true
			}
		} else {
			appAPIInfo := AppAPIInfo{
				AppUUID: appUUID,
				OsUUID:  osUUID,
				EmailID: email,
			}

			// create AppInfo(compose request)
			var newComposeRequest AppInfo
			var dockerInstances []dockerInstance

			for i, container := range list {
				if i == 0 {
					// Extract compose file content
					if _, ok := container.Labels["com.docker.compose.project.config_files"]; ok {
						composeDir := container.Labels["com.docker.compose.project.working_dir"]
						data, err := readFile(composeDir + "/docker-compose.yml")
						if err != nil {
							fmt.Printf("Error in reading compose file %s: %s\n", composeDir, err)
						}
						composeFileContent = string(data)
					}
				}

				meta, err := cli.ContainerInspect(context.Background(), container.ID)
				if err != nil {
					return appsData, isComposeStateChange, err
				}

				imagemeta, _, err := cli.ImageInspectWithRaw(context.Background(), container.ImageID)
				if err != nil {
					return appsData, isComposeStateChange, err
				}

				dockerInfo := dockerInstance{
					Names:     container.Names,
					Status:    container.State,
					Image:     container.Image,
					ImageID:   container.ImageID,
					Ports:     container.Ports,
					Size:      convertInt64ToString(imagemeta.Size),
					Networks:  container.HostConfig.NetworkMode,
					DockerID:  container.ID,
					Mounts:    container.Mounts,
					Tag:       container.Image,
					CreatedAt: meta.Created,
				}
				dockerInstances = append(dockerInstances, dockerInfo)
			}

			newComposeRequest.DockerInstance = dockerInstances
			newComposeRequest.ComposeFile = composeFileContent
			appsData[project] = AppRequest{
				AppAPIInfo: appAPIInfo,
				AppInfo:    newComposeRequest,
			}
			isComposeStateChange = true
		}
	}

	runMonitoringCompose(enableMonitoring, osUUID, appUUID)
	return appsData, isComposeStateChange, nil
}

// convertInt64ToString converts an int64 value to a string
func convertInt64ToString(value int64) string {
	return strconv.FormatInt(value, 10) // Base 10 conversion
}

// getCPUUsage returns the average CPU usage percentage and the number of CPU cores.
func getCPUUsage() (string, error) {

	// Get the number of CPU cores
	numCores, err := cpu.Counts(true)
	if err != nil {
		return "", err
	}

	// Convert the number of cores to a string
	return fmt.Sprintf("%d", numCores), nil
}

// getMemoryUsage returns a string summarizing the total memory of the server.
func getMemoryUsage() (string, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	// Summarize total memory in a single string (in GB)
	return fmt.Sprintf("%.2f GB", float64(vmStat.Total)/(1024*1024*1024)), nil
}

func runMonitoringCompose(monitoringEnable, osUUID, appUUID string) error {
	if monitoringEnable == "true" {
		configFilePath := os.Getenv("HOME") + "/.hossted/compose/monitoring/config.alloy"

		// Read the Grafana Agent config file
		configData, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatalf("Failed to read the Grafana Agent config file: %v", err)
		}

		// Replace the UUID placeholder with the actual UUID
		configStr := string(configData)

		configStr = strings.Replace(configStr, "${UUID}", fmt.Sprintf("\"%s\"", osUUID), -1)
		configStr = strings.Replace(configStr, "${APP_UUID}", fmt.Sprintf("\"%s\"", appUUID), -1)

		// Replace MIMIR_USERNAME and MIMIR_PASSWORD placeholders
		mimirUsername := common.MIMIR_USERNAME
		mimirPassword := common.MIMIR_PASSWORD
		mimirURL := common.MIMIR_URL
		lokiUsername := common.LOKI_USERNAME
		lokiPassword := common.LOKI_PASSWORD
		lokiURL := common.LOKI_URL

		if mimirUsername == "" || mimirPassword == "" || mimirURL == "" {
			log.Fatalf("MIMIR_USERNAME, MIMIR_URL and  MIMIR_PASSWORD environment variables must be set")
		}

		configStr = strings.Replace(configStr, "${MIMIR_USERNAME}", fmt.Sprintf("\"%s\"", mimirUsername), -1)
		configStr = strings.Replace(configStr, "${MIMIR_PASSWORD}", fmt.Sprintf("\"%s\"", mimirPassword), -1)
		configStr = strings.Replace(configStr, "${MIMIR_URL}", fmt.Sprintf("\"%s\"", mimirURL), -1)
		configStr = strings.Replace(configStr, "${LOKI_USERNAME}", fmt.Sprintf("\"%s\"", lokiUsername), -1)
		configStr = strings.Replace(configStr, "${LOKI_PASSWORD}", fmt.Sprintf("\"%s\"", lokiPassword), -1)
		configStr = strings.Replace(configStr, "${LOKI_URL}", fmt.Sprintf("\"%s\"", lokiURL), -1)

		// Write the updated config back to the file
		err = os.WriteFile(configFilePath, []byte(configStr), 0644)
		if err != nil {
			log.Fatalf("Failed to write the updated Grafana Agent config file: %v", err)
		}

		// Define the path to the Docker Compose file
		composeFile := os.Getenv("HOME") + "/.hossted/compose/monitoring/docker-compose.yaml"

		// Create the command to run Docker Compose
		cmd := exec.Command("docker", "compose", "-f", composeFile, "up", "-d")

		// Set the command's output to the standard output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Failed to execute Docker Compose: %v", err)
		}

		fmt.Println("Docker Compose executed successfully")
	}
	return nil
}

func GetAccessInfo(filepath string) *AccessInfo {
	file, err := os.Open(filepath)
	if err != nil {
		return &AccessInfo{}
	}
	defer file.Close()

	config := AccessInfo{
		URLs: []URLInfo{
			{},
		},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid line:", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "PROJECT_BASE_URL":
			config.URLs[0].URL = value
		case "H_EMAIL":
			config.URLs[0].User = value
		case "APP_PASSWORD":
			config.URLs[0].Password = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return &config
}

func getSoftwarePath() (string, error) {
	path := "/opt/hossted/run/software.txt"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", nil
	} else {
		return path, nil
	}

}

func getProjectName() (string, error) {
	path, err := getSoftwarePath()
	if err != nil {
		fmt.Println("Error getting software path", err)
	}

	// its a market place VM, access info object will exist
	if path == "/opt/hossted/run/software.txt" {
		// read the file in this path
		// file will have this convention - Linnovate-AWS-keycloak
		// capture the last word ie keycloak in this case.
		// and use this last work ie instead of osInfo.ProjectName
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return "", err
		}

		// The file will have the convention Linnovate-AWS-keycloak
		// Capture the last word (i.e., keycloak in this case)
		softwareName := strings.TrimSpace(string(data))
		words := strings.Split(softwareName, "-")
		if len(words) > 0 {
			projectName := words[len(words)-1]
			// Use this last word (i.e., keycloak) instead of osInfo.ProjectName
			return projectName, nil
		}
	} else if path == "" {
		fmt.Println("Contact Hossted support to add Access Info object")
		return "", nil
	}
	return "", nil
}

func isMarketplaceVM() (bool, error) {
	path, err := getSoftwarePath()
	if err != nil {
		return false, err
	}
	return path == "/opt/hossted/run/software.txt", nil
}

func getMarketplaceOSUUID() (string, error) {
	path := "/opt/hossted/run/uuid.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading UUID file: %w", err)
	}

	uuid := strings.TrimSpace(string(data))
	return uuid, nil
}

func getMarketplaceAppUUID(projectName string) (string, error) {

	ok, err := isMarketplaceVM()
	if err != nil {
		return "", err
	}
	if ok {
		path := "/opt/" + projectName + "/hossted/uuid.txt"
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("error reading UUID file: %w", err)
		}

		uuid := strings.TrimSpace(string(data))
		fmt.Println("Found existing appUUID for marketplace VM: ", uuid)
		return uuid, nil
	} else {
		fmt.Println("Generating appUUID for marketplace VM")
		return "A-" + uuid.NewString(), nil
	}
}
