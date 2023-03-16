package hossted

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

var (
	images string
)

func Scan(env string) error {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("No docker containers found")
	}

	for _, currentContainer := range containers {
		image, _, err := cli.ImageInspectWithRaw(ctx, currentContainer.Image)
		if err != nil {
			panic(err)
		}

		imageName := image.RepoTags[0]
		fileName := strings.Replace(imageName, "/", "-", -1)
		// Skip any images that don't have a name or tag
		if imageName == "<none>:<none>" || imageName == "anchore/grype:latest" {
			continue
		}

		// Print the name of the image being scanned
		fmt.Println("Scanning image:", imageName)

		reader, err := cli.ImagePull(ctx, "anchore/grype:latest", types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		defer reader.Close()

		_, err = cli.ContainerInspect(ctx, "grype")

		if err != nil {
			if !client.IsErrNotFound(err) {
				panic(err)
			}
			// Container does not exist
		} else {
			// Container exists
			if err := cli.ContainerRemove(ctx, "grype", types.ContainerRemoveOptions{Force: true}); err != nil {
				panic(err)
			}
		}

		//create grype docker container
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "anchore/grype",
			Tty:   false,
			Cmd:   []string{imageName, "-o", "json"},
		}, &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   "volume",
					Source: "docker_sock", // use an alphanumeric name for the volume
					Target: "/var/run/docker.sock",
				},
			},
		}, nil, nil, "grype")
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}

		if _, err := os.Stat("vulnerabilities"); os.IsNotExist(err) {
			// create the directory if it does not exist
			if err := os.Mkdir("vulnerabilities", 0755); err != nil {
				panic(err)
			}
		}
		file, err := os.Create(fmt.Sprintf("vulnerabilities/scan_%s.json", fileName))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Copy the container logs to the file
		stdcopy.StdCopy(file, os.Stderr, out)

		if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
			panic(err)
		}
		//send to ScanRequest
		ScanRequest(imageName, fileName, currentContainer.ID, env)
	}

	return nil
}

// ScanRequest sends security request
func ScanRequest(imageName, fileName, containerId, env string) (pingResponse, error) {

	var response pingResponse

	// Read file content
	scan, err := ioutil.ReadFile(fmt.Sprintf("vulnerabilities/scan_%s.json", fileName))
	if err != nil {
		return response, fmt.Errorf("Failed to read dockers file. %w", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fileName, "file")
	part.Write(scan)
	writer.Close()

	// Construct param map for input params
	params := make(map[string]string)
	params["file_name"] = fileName
	params["image_name"] = imageName
	params["docker_id"] = containerId
	req := HosstedRequest{
		// Endpoint env needs to replace in runtime for url parse to work. Otherwise runtime error.
		//EndPoint:     "https://api.__ENV__hossted.com/v1/instances/dockers",
		EndPoint:     "http://localhost:3004/v1/security", //"https://api.hossted.com/v1/instances/dockers", //"https://api.stage.hossted.com/v1/instances/dockers", //, // ,//
		Environment:  env,
		Params:       params,
		BearToken:    "Basic y5TXKDY4kTKbFcFtz9aD1pa2irmzhoziKPnEBcA8",
		SessionToken: "",
		TypeRequest:  "PATCH",
		ContentType:  writer.FormDataContentType(),
		Body:         body,
	}

	fmt.Printf("%s-Sending vulnerabilities. Please wait a moment... \n", imageName)
	resp, err := req.SendRequest()
	if err != nil {
		fmt.Println("err", err)
		return response, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return response, fmt.Errorf("Failed to parse JSON. %w", err)
	}

	return response, nil
}
