package hossted

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GetAppLogs goes to the app directory, then calls docker-compose logs
// Similar to ListAppPS func
func GetAppLogs(input string, followFlag bool) error {

	var app ConfigApplication
	config, err := GetConfig()
	if err != nil {
		fmt.Printf("Please call the command `hossted register` first.\n%w", err)
		os.Exit(0)
	}

	// Get App from prompt
	app, err = appPrompt(config.Applications, input)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if followFlag {
		cmd = exec.Command("sudo", "docker-compose", "logs", "-f")
	} else {
		cmd = exec.Command("sudo", "docker-compose", "logs")
	}

	cmd.Dir = app.AppPath
	fmt.Printf("Called command: %v\n", strings.Join(cmd.Args, " "))

	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return nil
}
