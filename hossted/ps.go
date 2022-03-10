package hossted

import (
	"fmt"
	"os"
	"os/exec"
)

// ListDockerPS goes to the app dire, then calls docker-compose ps
func ListDockerPS() error {

	config, err := GetConfig()
	if err != nil {
		fmt.Printf("Please call the command `hossted register` first.\n%w", err)
		os.Exit(0)
	}

	cmd := exec.Command("docker-compose", "ps")
	cmd.Dir = config.AppPath
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}
