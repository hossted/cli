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
	_ = app

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

// appPrompt prompt the user for which app to select
func appPrompt(apps []ConfigApplication) (ConfigApplication, error) {
	var (
		options   []string                     // options for applications
		configMap map[string]ConfigApplication // e.g. map[wikijs] -> ConfigApplication{}
	)
	configMap = make(map[string]ConfigApplication)

	for _, app := range apps {
		name := strings.TrimSpace(app.AppName)
		options = append(options, name)
		configMap[name] = app
	}

	prompt := promptui.Select{
		Label: "Applications",
		Items: options,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	fmt.Printf("Application: %q\n", result)
	return result, nil
}
