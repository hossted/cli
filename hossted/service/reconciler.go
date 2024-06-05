/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/docker/docker/api/types"
	"gopkg.in/yaml.v2"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

func Reconcile() error {
	emailsID, err := getEmail()
	if err != nil {
		return err
	}

	osFilePath, err := getComposeFilePath("compose.yaml")
	if err != nil {
		return err
	}

	appFilePath, err := getComposeFilePath("compose-request.json")
	if err != nil {
		return err
	}

	osUuid, err := setClusterUUID(emailsID, osFilePath)
	if err != nil {
		return err
	}

	list, err := listAllContainers()
	if err != nil {
		return err
	}

	err = setComposeRequest(appFilePath, list, osUuid, emailsID)
	if err != nil {
		return err
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

func setClusterUUID(email string, osFilePath string) (string, error) {
	var uuid string
	if _, err := os.Stat(osFilePath); os.IsNotExist(err) {
		uuid, err = updateUUID(osFilePath, email)
		if err != nil {
			return uuid, err
		}
	} else if err != nil {
		return uuid, err
	} else {
		uuid, err := checkUUID(osFilePath)
		if err != nil {
			return uuid, err
		}
		if uuid == "" {
			uuid, err = updateUUID(osFilePath, email)
			if err != nil {
				return uuid, err
			}
		}
	}
	return uuid, nil
}

func updateUUID(osFilePath string, email string) (string, error) {
	osUUID := "D-" + uuid.NewString()
	fmt.Println("Generating UUID for cluster: ", osUUID)
	info := OsInfo{
		OsUUID:  osUUID,
		EmailID: email,
	}

	yamlData, err := yaml.Marshal(info)
	if err != nil {
		fmt.Printf("error in YAML marshaling: %s\n", err)
		return osUUID, err
	}

	err = writeFile(osFilePath, yamlData)
	if err != nil {
		return osUUID, err
	}

	return osUUID, nil
}

func checkUUID(osFilePath string) (string, error) {
	var osData OsInfo
	data, err := readFile(osFilePath)
	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal(data, &osData)
	if err != nil {
		fmt.Printf("Error unmarshaling Yaml: %s\n", err)
		return "", err
	}

	fmt.Println("Registering cluster with UUID: ", osData.OsUUID)
	return osData.OsUUID, nil
}

func setComposeRequest(appFilePath string, containerList []types.Container, osUuid string, email string) error {
	// prepare appsInfo with updated container status
	appRequest, err := prepareComposeRequest(appFilePath, containerList, osUuid, email)
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(appRequest, "", "    ")
	if err != nil {
		fmt.Printf("error marshaling JSON: %s\n", err)
		return err
	}

	// write compose status file
	err = writeFile(appFilePath, jsonData)
	if err != nil {
		return err
	}

	return nil
}

func listAllContainers() ([]types.Container, error) {
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

func prepareComposeRequest(appFilePath string, containerList []types.Container, osUuid string, email string) (map[string]AppRequest, error) {
	var appsData map[string]AppRequest

	if _, err := os.Stat(appFilePath); os.IsNotExist(err) {
		appsData = make(map[string]AppRequest)
	} else {
		data, err := readFile(appFilePath)
		if err != nil {
			fmt.Printf("Error in reading %s file: %s\n", appFilePath, err)
			return appsData, err
		}
		err = json.Unmarshal(data, &appsData)
		if err != nil {
			fmt.Printf("Error in JSON unmarshaling %s file: %s\n", appFilePath, err)
			return appsData, err
		}
	}

	uniqueProjects, err := getUniqueComposeProjects(containerList)
	if err != nil {
		return appsData, err
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
			return appsData, err
		}

		// Check if project is present in the existingApps
		if found := existingApps[project]; found {
			// create only AppInfo(compose request)
			var newComposeRequest AppInfo
			var containersInfo []ContainerInfo

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
				containerInfo := ContainerInfo{
					Name:   name,
					Status: container.State,
					Image:  container.Image,
				}
				containersInfo = append(containersInfo, containerInfo)
			}

			newComposeRequest.ContainerInfo = containersInfo
			newComposeRequest.ComposeFile = composeFileContent

			// get previous compose request
			prevComposeRequest := appsData[project].AppInfo

			//compare newComposeRequest and prevComposeRequest
			res := reflect.DeepEqual(prevComposeRequest, newComposeRequest)
			if !res {
				appsData[project] = AppRequest{
					AppAPIInfo: appsData[project].AppAPIInfo,
					AppInfo:    newComposeRequest,
				}
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
			var containersInfo []ContainerInfo

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
				containerInfo := ContainerInfo{
					Name:   name,
					Status: container.State,
					Image:  container.Image,
				}
				containersInfo = append(containersInfo, containerInfo)
			}

			newComposeRequest.ContainerInfo = containersInfo
			newComposeRequest.ComposeFile = composeFileContent
			appsData[project] = AppRequest{
				AppAPIInfo: appAPIInfo,
				AppInfo:    newComposeRequest,
			}
		}
	}

	return appsData, nil
}
