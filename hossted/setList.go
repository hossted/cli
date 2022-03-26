package hossted

import (
	"errors"
	"fmt"
	"os"
)

// TODO: list available commands for the available services on the machine
func ListCommands() error {

	// Get available apps
	config, err := GetConfig()
	apps := config.Applications

	if err != nil {
		fmt.Printf("Please call the command `hossted register` first.\n%w", err)
		os.Exit(0)
	}

	// Get all available apps/commands
	m, err := getCommandsMap(AVAILABLE)
	if err != nil {
		return err
	}

	// Only print the application available on the vm
	var check bool // check if there is any apps available
	for _, app := range apps {
		if commands, ok := m[app.AppName]; ok {
			check = true
			fmt.Println(app)
			fmt.Println("--------------")
			fmt.Println(commands)
			fmt.Println("")
		}
	}
	if !check {
		return errors.New("No available applications. Please check the file /opt/linnovate/run/uuid.txt.\n")
	}

	fmt.Println(m)
	return nil
}
