package hossted

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

// Generall available commands, for all apps
var GAVAILABLE = `
apps:
  - app: general
    commands: [domain, auth]
    values: [example.com, false]
`

// Available commands in yaml format. If a new set of apps/commands needs to be supported,
// need to append the values here
// TODO: Add general command
// TODO: Handle logic for command group
var AVAILABLE = `
apps:
  - app: prometheus
    group: set
    commands: [domain, auth]
    values: [example.com, false]

  - app: airflow
    group: set
    commands: [domain, auth]
    values: [example.com, false]

  - app: wordpress
    group: set
    commands: [domain, auth]
    values: [example.com, false]

  - app: wph
    group: set
    commands: [domain, auth]
    values: [example.com, false]

  - app: gitbucket
    group:
    commands: [httpopen]
    values: [""]

  - app: demo
    group:
    commands: [abc, def]
    values: [abc, def]
`

// CheckCommands check whether the app, and corresponding commands/subcommands are supported.
// Return error if the provided values are not in the pre-defined list
func CheckCommands(app, command string) error {

	// Get the map of available apps and commands
	m, err := getCommandsMap(AVAILABLE)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s.%s", app, command) // e.g. promethus.domain

	if _, ok := m[key]; ok {
		// happy path. app.command is available
		return nil
	} else {
		// app not supported
		return fmt.Errorf("Provided application is not supported - %s\nPlease use `hossted set list` to check the available commands.\n", key)
	}

	return nil
}

// getCommandsMap gets a mapping for available apps and commands mapping
// input as the yaml formatted available commands
func getCommandsMap(input string) (AvailableCommandMap, error) {

	// Available commands map, kv as map[appName.command] -> available commands, []Command
	// e.g. map["prometheus.domain"] -> []Command[{prometheus domain example.com}]
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
		cg := app.CommandGroup

		if len(app.Commands) != len(app.Values) {
			return m, errors.New("Length of commands does not equal to the length of sample values.\n Please check the available command yaml.")
		}
		for i, _ := range app.Commands {
			name := fmt.Sprintf("%s.%s", appName, app.Commands[i]) // e.g. prometheus.domain
			c := Command{
				App:          appName,
				CommandGroup: cg,
				Command:      app.Commands[i],
				Value:        app.Values[i],
			}
			m[name] = c
		}
	}

	return m, nil
}
