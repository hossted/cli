package hossted

import (
	"fmt"
	"os/exec"
)

// ListAppPS goes to the app directory, then calls docker-compose ps
// if input is "", call prompt to get user input, otherwise look up the application in the config
func ListAppPS(input string) error {
	var app ConfigApplication

	config, err := GetConfig()
	if err != nil {
		// fmt.Printf("Please call the command `hossted register` first.\n%w", err)
		// os.Exit(0)
		return fmt.Errorf("Please call the command `hossted register` first.\n%w", err)
	}

	app, err = appPrompt(config.Applications, input)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "docker-compose", "ps")
	cmd.Dir = app.AppPath

	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return nil
}
