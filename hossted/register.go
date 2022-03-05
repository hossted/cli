package hossted

import (
	"bufio"
	"embed"
	"errors"

	"fmt"
	"html/template"
	"io"

	"github.com/manifoldco/promptui"
	"github.com/spf13/hossted/utils"
)

var (
	//go:embed templates
	templates embed.FS
)

// RegisterUsers updates email, organization, etc,.. in the yaml file
func RegisterUsers() error {

	config, _ := GetConfig() // Ignore error
	_ = config

	email, _ := emailPrompt()

	// Assign back to config object
	config.Email = email
	config.Organization = "Axa"

	// Write back to file
	fmt.Println(utils.PrettyPrint(config))
	fmt.Println("Register User")
	return nil
}

// WriteConfig writes the config to the config file (~/.hossted/config.yaml)
func WriteConfig(w io.Writer, config Config) error {

	// Read Template
	t, err := template.ParseFS(templates, "templates/config.tmpl")
	if err != nil {
		return err
	}

	// Write to template
	err = t.Execute(w, config)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(w)
	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// registerPromp prompt the user for email and organization
func emailPrompt() (string, error) {
	validate := func(input string) error {
		if len(input) <= 5 {
			return errors.New("Invalid length. Must be larger than 5.")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Email",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
}
