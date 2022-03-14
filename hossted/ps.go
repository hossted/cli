package hossted

import (
	"fmt"
	"os"
	"os/exec"
)

// ListDockerPS goes to the app directory, then calls docker-compose ps
func ListDockerPS() error {

	config, err := GetConfig()
	if err != nil {
		fmt.Printf("Please call the command `hossted register` first.\n%w", err)
		os.Exit(0)
	}

	cmd := exec.Command("docker-compose", "ps")
	for _, app := range config.Applications {
		cmd.Dir = app.AppPath
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		fmt.Println(out)
	}

	return nil
}
