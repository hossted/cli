package hossted

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ListCommands lists the available services and general commands on the virtual machine.
// The available application is defined under the file /opt/linoovate/run/uuid.txt.
// Also it would depends on whether the "action" is predefined as available of the certain app.
func ListCommands() error {

	// Get available apps
	config, err := GetConfig()
	if err != nil {
		return err
	}

	apps := config.Applications

	// Applications available on the vm
	vmAppsMap := make(map[string]bool) // e.g. map["prometheus"]true
	for _, app := range apps {
		vmAppsMap[app.AppName] = true
	}

	// Get all available apps/commands
	m, err := getCommandsMap(GAVAILABLE, AVAILABLE)
	if err != nil {
		return err
	}

	// Print general commands
	var generalCommands []Command
	for k, v := range m { // k: app.command, v: Command{}
		if strings.Contains(k, "general") {
			generalCommands = append(generalCommands, v)
		}
	}

	// Check matching commands
	// TODO: Check if this is still vaild after moving almost all commands to General Available (#39)
	var validCommands []Command
	for k, v := range m { // k: app.command, v: Command{}
		app := getAppNameFromKey(k)
		if _, ok := vmAppsMap[app]; ok { // append to validCommands only if its on vm
			validCommands = append(validCommands, v)
		}
	}

	// Sort
	sort.Slice(validCommands, func(i, j int) bool {
		return validCommands[i].App < validCommands[j].App
	})

	// Check
	allCommands := append(generalCommands, validCommands...)
	if len(allCommands) == 0 {
		return errors.New("No available commands. Please check with administrator.")
	}

	// Print the commands if available
	if len(generalCommands) > 0 {
		_ = printCommands(generalCommands)
	}
	if len(validCommands) > 0 {
		_ = printCommands(validCommands)
	}

	return nil
}

func getAppNameFromKey(key string) string {
	var app string
	s := strings.Split(key, ".")
	if len(s) > 0 {
		app = s[0] // Get the app/first part from the key
	}
	return app
}

// printCommands prints the list of available commands (general, app, etc,..) to the console
// to provide additional information to users on the available commands on the vm
func printCommands(commands []Command) error {

	if len(commands) == 0 {
		return errors.New("No available commands. Please check with administrator.")
	}

	// List the available commands (vm + predefined)
	var prev string // For formatting only. Group same apps together.
	for _, c := range commands {
		app := c.App
		if prev != app {
			prev = app
			fmt.Println("")
			fmt.Println(app)
			fmt.Println("------------")
		}
		// Handle general commands
		if c.App == "general" {
			c.App = ""
		}

		msg := fmt.Sprintf("hossted %s %s %s %s\n", c.CommandGroup, c.Command, c.App, c.Value)

		// Replace multiple spaces to one
		space := regexp.MustCompile(`\s+`)
		s := space.ReplaceAllString(msg, " ")
		fmt.Println(s)
	}
	fmt.Println("")
	return nil
}
