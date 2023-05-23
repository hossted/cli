package hossted

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func SetMonitoring(env string, flag bool) error {

	if !HasContainerRunning() {
		fmt.Println("The application still in configuration")
		os.Exit(0)
	}

	config, _ := GetConfig()

	config.Monitoring = flag

	err := WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}
	fmt.Println("monitoring set to", flag)

	//send activity log about the command
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}
	fullCommand := "hossted set monitoring " + fmt.Sprint(flag)
	options := `{"monitoring":` + fmt.Sprint(flag) + `}`
	typeActivity := "set_monitoring"
	sendActivityLog(env, uuid, fullCommand, options, typeActivity)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	if flag {
		//check if container exists in container list
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}
		for _, container := range containers {
			if container.Names[0] == "/monitoring" {
				fmt.Println("Monitoring already running, you can see your dashboard at https://app.hossted.com")
				os.Exit(0)
			}
		}
	}

	//check if container exists
	_, err = cli.ContainerInspect(ctx, "monitoring")

	if err != nil {
		if !client.IsErrNotFound(err) {
			panic(err)
		}
		// Container does not exist
	} else {
		// Container exists
		stopAndRemoveMonitoring(cli)
	}

	if flag {

		auth := types.AuthConfig{
			Username: "Username",
			Password: "Password",
		}

		authData, err := json.Marshal(auth)
		if err != nil {
			return err
		}

		auths := base64.URLEncoding.EncodeToString(authData)
		out, err := cli.ImagePull(
			ctx,
			"linnovate.azurecr.io/hossted/grafana-agent:latest",
			types.ImagePullOptions{
				RegistryAuth: auths,
			})
		if err != nil {
			return fmt.Errorf("failed to pull image: %w", err)
		}
		defer out.Close()

		_, err = io.Copy(io.Discard, out)
		if err != nil {
			return fmt.Errorf("failed to read image logs: %w", err)
		}

		// Create the hossted-agent container
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image:      "linnovate.azurecr.io/hossted/grafana-agent:latest",
			Entrypoint: []string{"/bin/grafana-agent", "-config.file=/etc/agent/agent.yaml", "-metrics.wal-directory=/etc/agent/data"},
		}, &container.HostConfig{
			Privileged: true,
			Binds: []string{
				"/:/rootfs:ro",
				"/var/run:/var/run:rw",
				"/sys:/sys:ro",
				"/var/lib/docker/:/var/lib/docker:ro",
			},
		}, nil, nil, "monitoring")
		if err != nil {
			panic(err)
		}

		// Start the container
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		// Wait for the container to finish running
		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		// Retrieve the container logs
		out, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}

		io.Copy(io.Discard, out)

	}

	return nil
}

func stopAndRemoveMonitoring(cli *client.Client) error {
	ctx := context.Background()

	err := cli.ContainerStop(ctx, "monitoring", nil)
	if err != nil {
		panic(err)
	}
	err = cli.ContainerRemove(ctx, "monitoring", types.ContainerRemoveOptions{Force: true})
	if err != nil {

		panic(err)
	}

	return nil
}
