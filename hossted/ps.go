package hossted

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

// ListAppPS goes to the app directory, then calls docker-compose ps
// if input is "", call prompt to get user input, otherwise look up the application in the config
func ListAppPS(input string) error {
	var app ConfigApplication

	config, err := GetConfig()
	if err != nil {
		fmt.Printf("Please call the command `hossted register` first.\n%w", err)
		os.Exit(0)
	}

	app, err = appPrompt(config.Applications, input)
	if err != nil {
		return err
	}

	cmd := exec.Command("docker-compose", "ps")
	cmd.Dir = app.AppPath
	fmt.Printf("Called command: %v\n", strings.Join(cmd.Args, " "))

	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(out)

	return nil
}

// appPrompt prompt the user for which app to select
func appPrompt(apps []ConfigApplication, input string) (ConfigApplication, error) {
	var (
		options   []string                     // options for applications
		configMap map[string]ConfigApplication // e.g. map[wikijs] -> ConfigApplication{}
		app       ConfigApplication            // Config for selected App
		selected  string
	)
	configMap = make(map[string]ConfigApplication)
	input = strings.TrimSpace(input)

	// Build select options and mapping
	for _, app := range apps {
		name := strings.TrimSpace(app.AppName)
		options = append(options, name)
		configMap[name] = app
	}

	// If input is empty, prompt for user to input
	// Otherwise use the provided value as application
	if input == "" {
		// Prompt for selection
		prompt := promptui.Select{
			Label: "Applications",
			Items: options,
		}

		_, selected, err := prompt.Run()
		input = selected

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return app, err
		}
	} else {
		selected = input
	}

	// Return selected app config
	if val, ok := configMap[selected]; ok {
		app = val
	} else {
		return app, fmt.Errorf("Invalid selection for app. Available applications are [%v]", options)
	}

	return app, nil
}
