package hossted

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Available commands in yaml format. If a new set of apps/commands needs to be supported,
// need to append the values here
var AVAILABLE = `
apps:
  - app: prometheus
    commands: [url]
    values: [example.com]

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
	if ac, ok := m[app]; ok { // available commands
		// Check the input command in the available ones.
		for _, c := range ac {
			if c == command {
				return nil
			}
		}

		// commands not supported for certain app
		return fmt.Errorf("Provided command is not supported in  - %s\nPlease use `hossted set list` to check the available commands.\n", command)

	} else {
		// TODO: Add supported apps
		// app not supported
		return fmt.Errorf("Provided application is not support - %s\nPlease use `hossted set list` to check the available commands.\n", app)
	}

	return nil
}

// getCommandsMap gets a mapping for available apps and commands mapping
// input as the yaml formatted available commands
func getCommandsMap(input string) (AvailableCommandMap, error) {

	// Available commands map, kv as map[appName] -> available commands, []string
	// e.g. map["prometheus"] -> ["url", "xxx"]
	var (
		m         AvailableCommandMap // result available map
		available AvailableCommand    // For parsing yaml
	)
	m = make(AvailableCommandMap)
	err := yaml.Unmarshal([]byte(input), &available)
	if err != nil {
		return m, fmt.Errorf("Can not parse avilable commands yaml. %w", err)
	}

	// k as app, v as commands
	for _, val := range available.Commands {
		m[val.App] = val
	}

	return m, nil
}
