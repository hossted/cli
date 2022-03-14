package hossted

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

// ListAppPS goes to the app directory, then calls docker-compose ps
func ListAppPS() error {

	config, err := GetConfig()
	if err != nil {
		fmt.Printf("Please call the command `hossted register` first.\n%w", err)
		os.Exit(0)
	}

	// Get App from prompt
	app, err := appPrompt(config.Applications)
	if err != nil {
		return err
	}
	fmt.Println(app.AppName)

	cmd := exec.Command("docker-compose", "ps")
	cmd.Dir = app.AppPath
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(out)

	return nil
}

// appPrompt prompt the user for which app to select
func appPrompt(apps []ConfigApplication) (ConfigApplication, error) {
	var (
		options   []string                     // options for applications
		configMap map[string]ConfigApplication // e.g. map[wikijs] -> ConfigApplication{}
		app       ConfigApplication            // Config for selected App
	)
	configMap = make(map[string]ConfigApplication)

	// Build select options and mapping
	for _, app := range apps {
		name := strings.TrimSpace(app.AppName)
		options = append(options, name)
		configMap[name] = app
	}

	// Prompt for selection
	prompt := promptui.Select{
		Label: "Applications",
		Items: options,
	}

	_, input, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return app, err
	}
	fmt.Printf("Application: %q\n", input)

	// Return selected app config
	if val, ok := configMap[input]; ok {
		app = val
	} else {
		return app, fmt.Errorf("Invalid selection for app. Available applications are [%v]", options)
	}

	return app, nil
}
