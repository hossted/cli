/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"

	// "github.com/docker/docker/api/types"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

func Reconcile() error {
	emailsID, err := getEmail()
	if err != nil {
		return err
	}

	filePath, err := getComposeStatusFilePath()
	if err != nil {
		return err
	}

	err = setClusterUUID(emailsID, filePath)
	if err != nil {
		return err
	}

	err = setAppStatus(filePath)
	if err != nil {
		return err
	}

	return nil

}

func getComposeStatusFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %s\n", err)
		return "", err
	}

	// Construct the full path to compose-status.json
	filePath := filepath.Join(homeDir, ".hossted", "compose-status.json")
	return filePath, nil
}

func getComposeStatusFileJson(filePath string) (AppInfo, error) {
	var appData AppInfo
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist.\n", filePath)
		} else {
			fmt.Printf("Error opening file: %s\n", err)
		}
		return appData, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return appData, err
	}

	err = json.Unmarshal(data, &appData)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %s\n", err)
		return appData, err
	}
	return appData, nil
}

func writeComposeStatusFileJson(filePath string, appData AppInfo) error {
	jsonData, err := json.MarshalIndent(appData, "", "    ")
	if err != nil {
		fmt.Printf("Write compose-status.json - error marshaling JSON: %s\n", err)
		return err
	}

	// Write the JSON data to the file
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Write compose-status.json - error writing JSON: %s\n", err)
		return err
	}
	return nil
}

func setClusterUUID(email string, filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = updateUUID(filePath, email)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		uuid, err := checkUUID(filePath)
		if err != nil {
			return err
		}
		if uuid == "" {
			err = updateUUID(filePath, email)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateUUID(filePath string, email string) error {
	clusterUUID := "D-" + uuid.NewString()
	fmt.Println("Generating UUID for cluster: ", clusterUUID)
	info := AppInfo{
		ClusterInfo: ClusterInfo{
			ClusterUUID: clusterUUID,
			EmailID:     email,
		},
		ContainerStatus: []ContainerStatus{},
	}

	err := writeComposeStatusFileJson(filePath, info)
	if err != nil {
		return err
	}

	return nil
}

func checkUUID(filePath string) (string, error) {
	appData, err := getComposeStatusFileJson(filePath)
	if err != nil {
		return "", err
	}
	fmt.Println("Registering cluster with UUID: ", appData.ClusterInfo.ClusterUUID)
	return appData.ClusterInfo.ClusterUUID, nil
}

func setAppStatus(filePath string) error {
	// list the running containers
	list, err := listContainers()
	if err != nil {
		return err
	}

	// prepare appInfo with updated container status
	appInfo, err := prepareContainerStatus(filePath, list)
	if err != nil {
		return err
	}

	// write compose status file
	err = writeComposeStatusFileJson(filePath, appInfo)
	if err != nil {
		return err
	}

	return nil
}

func listContainers() ([]types.Container, error) {
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

func prepareContainerStatus(filePath string, containerList []types.Container) (AppInfo, error) {
	var containerStatusList []ContainerStatus
	appData, err := getComposeStatusFileJson(filePath)
	if err != nil {
		return appData, err
	}

	// Create a map for existing containers
	existingContainer := make(map[string]ContainerStatus)
	for _, status := range appData.ContainerStatus {
		existingContainer[status.Name] = status
	}

	for _, container := range containerList {
		containerName := container.Names[0]
		if strings.HasPrefix(container.Names[0], "/") {
			containerName = container.Names[0][1:]
		}

		if status, found := existingContainer[containerName]; found {
			fmt.Println("container already exists: ", status)
			// Update existing status
			status.Status = container.State
			containerStatusList = append(containerStatusList, status)
		} else {
			// Add new container and its status
			appUUID := "A-" + uuid.NewString()
			status = ContainerStatus{
				AppUUID: appUUID,
				Name:    containerName,
				Image:   container.Image,
				Status:  container.State,
			}
			fmt.Println("adding new container: ", status)
			containerStatusList = append(containerStatusList, status)
		}
	}

	appData.ContainerStatus = containerStatusList
	return appData, nil
}