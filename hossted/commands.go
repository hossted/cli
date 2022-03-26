package hossted

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

// Available commands in yaml format. If a new set of apps/commands needs to be supported,
// need to append the values here
var AVAILABLE = `
apps:
  - app: prometheus
    commands: [url, auth]
    values: [example.com, false]

  - app: demo
    commands: [abc, def]
    values: [abc, def]
`

// CheckCommands check whether the app, and corresponding commands/subcommands are supported.
// Return error if the provided values are not in the pre-defined list
// TODO: cross check available apps in config
// TODO: Error handling. Add a list of availabe app/commands, etc,..
func CheckCommands(app, command string) error {

	// Get the map of available apps and commands
	m, err := getCommandsMap(AVAILABLE)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s.%s", app, command) // e.g. promethus.url

	if _, ok := m[key]; ok {
		// happy path. app.command is available
		return nil
	} else {
		// TODO: Add supported apps
		// app not supported
		return fmt.Errorf("Provided application is not support - %s\nPlease use `hossted set list` to check the available commands.\n", key)
	}

	return nil
}

// getCommandsMap gets a mapping for available apps and commands mapping
// input as the yaml formatted available commands
func getCommandsMap(input string) (AvailableCommandMap, error) {

	// Available commands map, kv as map[appName.command] -> available commands, []Command
	// e.g. map["prometheus.url"] -> []Command[{prometheus url example.com}]
	var (
		m         AvailableCommandMap // result available map
		available AvailableCommand    // For parsing yaml
	)
	m = make(AvailableCommandMap)
	err := yaml.Unmarshal([]byte(input), &available)
	if err != nil {
		return m, fmt.Errorf("Can not parse avilable commands yaml. %w", err)
	}
	if len(available.Apps) == 0 {
		return m, errors.New("No available apps and commands. Please check.")
	}

	// k as app, v as commands
	for _, app := range available.Apps {
		appName := app.App

		if len(app.Commands) != len(app.Values) {
			return m, errors.New("Length of commands does not equal to the length of sample values.\n Please check the available command yaml.")
		}
		for i, _ := range app.Commands {
			name := fmt.Sprintf("%s.%s", appName, app.Commands[i]) // e.g. prometheus.url
			c := Command{
				App:     appName,
				Command: app.Commands[i],
				Value:   app.Values[i],
			}
			m[name] = c
		}
	}

	return m, nil
}
