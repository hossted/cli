package hossted

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

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
	}
	selected = input

	// Return selected app config
	if val, ok := configMap[selected]; ok {
		app = val
	} else {
		return app, fmt.Errorf("Invalid selection for app. Available applications are [%v]", options)
	}

	return app, nil
}
