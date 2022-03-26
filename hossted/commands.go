package hossted

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// CheckCommands check whether the app, and corresponding commands/subcommands are supported.
// Return error if the provided values are not in the pre-defined list
// TODO: cross check available apps in config
func CheckCommands() error {

	// Available commands map, kv as map[appName] -> available commands, []string
	// e.g. map["prometheus"] -> ["url", "xxx"]
	var commands []AvailableCommand
	available := `
apps:
  - app: server1
    commands: [pepito, asd]
  - app: server2
    commands: [juanito]
`

	err := yaml.Unmarshal([]byte(available), &commands)
	if err != nil {
		return err
	}
	fmt.Println(commands)
	return nil
}
