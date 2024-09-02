package compose

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

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

	appFilePath, err := getComposeFilePath("compose-request.json")
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
		osInfo.EmailID,
		enableMonitoring,
		osInfo.ProjectName)
	if err != nil {
		return err
	}
	err = sendComposeInfo(appFilePath, osInfo)
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

func getComposeFilePath(filename string) (string, error) {
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

func setClusterInfo(osInfo OsInfo, osFilePath string) (OsInfo, error) {
	info := osInfo
	osUUID, err := checkUUID(osFilePath)
	if err != nil {
		return info, err
	}

	if osUUID == "" {
		osUUID := "D-" + uuid.NewString()
		fmt.Println("Generating UUID for cluster: ", osUUID)
		info.OsUUID = osUUID
	}

	yamlData, err := yaml.Marshal(info)
	if err != nil {
		fmt.Printf("error in YAML marshaling: %s\n", err)
		return info, err
	}

	err = writeFile(osFilePath, yamlData)
	if err != nil {
		return info, err
	}

	return info, nil
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

	fmt.Println("Registering cluster with UUID: ", osData.OsUUID)
	return osData.OsUUID, nil
}

func writeComposeRequest2File(
	appFilePath string,
	containerList []types.Container,
	osUuid,
	email,
	enableMonitoring,
	projectName string) (bool, error) {
	// prepare appsInfo with updated container status
	appRequest, isComposeStateChange, err := prepareComposeRequest(
		appFilePath,
		containerList,
		osUuid,
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

type optionsState struct {
	Monitoring bool `json:"monitoring"`
	Logging    bool `json:"logging"`
	CVEScan    bool `json:"cve_scan"`
}

type URLInfo struct {
	URL      string `json:"url"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type AccessInfo struct {
	URLs []URLInfo `json:"urls"`
}

type request struct {
	UUID         string       `json:"uuid"`
	OsUUID       string       `json:"osuuid"`
	OrgID        string       `json:"org_id"`
	Email        string       `json:"email"`
	Type         string       `json:"type"`
	Product      string       `json:"product"`
	CPUNum       string       `json:"cpunum"`
	Memory       string       `json:"memory"`
	OptionsState optionsState `json:"options_state"`
	ComposeFile  string       `json:"compose_file"`
	AccessInfo   AccessInfo   `json:"access_info"`
}

type dockerInstance struct {
	DockerID   string      `json:"docker_id"`
	InstanceID string      `json:"instance_id"`
	ImageID    string      `json:"image_id"`
	Ports      interface{} `json:"ports"`
	Status     string      `json:"status"`
	Size       int64       `json:"size"`
	Names      string      `json:"names"`
	Mounts     interface{} `json:"mounts"`
	Networks   string      `json:"networks"`
	Tag        string      `json:"tag"`
	Image      string      `json:"image"`
}

func sendComposeInfo(appFilePath string, osInfo OsInfo) error {
	hosstedAPIUrl := osInfo.HosstedApiUrl
	orgID := osInfo.OrgID
	token := osInfo.Token

	composeInfo, err := readFile(appFilePath)
	if err != nil {
		return err
	}

	composeUrl := hosstedAPIUrl + "/compose/hosts"
	containersUrl := hosstedAPIUrl + "/compose/containers"

	access_info := getAccessInfo("/opt/" + osInfo.ProjectName + "/.env")

	var data map[string]AppRequest
	err = json.Unmarshal(composeInfo, &data)
	if err != nil {
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
	for appName, compose := range data {
		newReq := request{
			UUID:       compose.AppAPIInfo.AppUUID,
			OsUUID:     compose.AppAPIInfo.OsUUID,
			Email:      compose.AppAPIInfo.EmailID,
			AccessInfo: *access_info,
			OrgID:      orgID,
			Type:       "compose",
			Product:    appName,
			CPUNum:     cpu,
			Memory:     mem,
			OptionsState: optionsState{
				Monitoring: true,
				Logging:    true,
				CVEScan:    true,
			},
			ComposeFile: compose.AppInfo.ComposeFile,
		}

		body, err := json.Marshal(newReq)
		if err != nil {
			return err
		}

		err = common.HttpRequest("POST", composeUrl, token, body)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully registered app [%s] with appID [%s]\n", appName, compose.AppAPIInfo.AppUUID)
	}

	for appName, info := range data {
		for _, ci := range info.AppInfo.DockerInstance {
			newDI := dockerInstance{
				DockerID:   ci.DockerID,
				InstanceID: info.AppAPIInfo.AppUUID,
				ImageID:    ci.ImageID,
				Ports:      ci.Ports,
				Status:     ci.Status,
				Size:       ci.Size,
				Names:      ci.Names,
				Networks:   ci.Networks,
				Image:      ci.Image,
				Mounts:     ci.Mounts,
			}
			newDIBody, err := json.Marshal(newDI)
			if err != nil {
				return err
			}
			err = common.HttpRequest("POST", containersUrl, token, newDIBody)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully registered docker info, appName:[%s], dockerName:[%s], appID:[%s]\n", appName, ci.Names, info.AppAPIInfo.AppUUID)
		}
	}

	return nil
}

func listAllContainers(projectName string) ([]types.Container, error) {
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	// Define context
	ctx := context.Background()

	// Define filters to list only running containers
	filter := filters.NewArgs()
	//filter.Add("status", "running")

	filter.Add("label", "com.docker.compose.project="+projectName)

	filter.Add("label", "com.docker.compose.project")
	filter.Add("label", "com.docker.compose.config-hash")

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})

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
	osUuid,
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
						composeFiles := container.Labels["com.docker.compose.project.config_files"]
						composeFilePath := strings.Split(composeFiles, ",")[0]
						data, err := readFile(composeFilePath)
						if err != nil {
							fmt.Printf("Error in reading compose file %s: %s\n", composeFilePath, err)
						}
						composeFileContent = string(data)
					}
				}

				name := container.Names[0]
				if strings.HasPrefix(container.Names[0], "/") {
					name = container.Names[0][1:]
				}
				containerInfo := dockerInstance{
					Names:    name,
					Status:   container.State,
					Image:    container.Image,
					ImageID:  container.ImageID,
					Ports:    container.Ports,
					Size:     container.SizeRw,
					Networks: container.HostConfig.NetworkMode,
					DockerID: container.ID,
					Mounts:   container.Mounts,
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
			// create AppAPIInfo
			appUUID := "A-" + uuid.NewString()
			appAPIInfo := AppAPIInfo{
				AppUUID: appUUID,
				OsUUID:  osUuid,
				EmailID: email,
			}

			// create AppInfo(compose request)
			var newComposeRequest AppInfo
			var dockerInstances []dockerInstance

			for i, container := range list {
				if i == 0 {
					// Extract compose file content
					if _, ok := container.Labels["com.docker.compose.project.config_files"]; ok {
						composeFiles := container.Labels["com.docker.compose.project.config_files"]
						composeFilePath := strings.Split(composeFiles, ",")[0]
						data, err := readFile(composeFilePath)
						if err != nil {
							fmt.Printf("Error in reading compose file %s: %s\n", composeFilePath, err)
						}
						composeFileContent = string(data)
					}
				}

				name := container.Names[0]
				if strings.HasPrefix(container.Names[0], "/") {
					name = container.Names[0][1:]
				}
				dockerInfo := dockerInstance{
					Names:    name,
					Status:   container.State,
					Image:    container.Image,
					ImageID:  container.ImageID,
					Ports:    container.Ports,
					Size:     container.SizeRw,
					Networks: container.HostConfig.NetworkMode,
					DockerID: container.ID,
					Mounts:   container.Mounts,
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

	runMonitoringCompose(enableMonitoring, osUuid, appsData[projectName].AppAPIInfo.AppUUID)
	return appsData, isComposeStateChange, nil
}

// GetCPUUsage returns the average CPU usage percentage as a formatted string
func getCPUUsage() (string, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", err
	}

	if len(percentages) > 0 {
		return fmt.Sprintf("%.2f%%", percentages[0]), nil
	}

	return "", fmt.Errorf("no CPU usage data available")
}

// GetMemoryUsage returns the memory usage statistics as a formatted string
func getMemoryUsage() (string, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	memUsage := fmt.Sprintf("%.2f%%", vmStat.UsedPercent)
	return memUsage, nil
}

func runMonitoringCompose(monitoringEnable, osUUID, appUUID string) error {
	if monitoringEnable == "true" {
		configFilePath := os.Getenv("HOME") + "/.hossted/compose/monitoring/config.river"

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

func getAccessInfo(filepath string) *AccessInfo {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open file for Access Info: %s", err)
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
