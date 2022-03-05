package hossted

import (
	"bufio"
	"embed"
	"errors"
	"regexp"

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

	// Prompt user for input
	email, _ := emailPrompt()
	company, _ := companyPrompt()

	// Assign back to config object
	config.Email = email
	config.Organization = company

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

// emailPromp prompt the user for email
func emailPrompt() (string, error) {

	// Regex for email checking
	re := regexp.MustCompile(`\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+`)

	validate := func(input string) error {
		if len(input) <= 5 {
			return errors.New("Invalid length. Must be larger than 5.")
		}

		if !re.MatchString(input) {
			return errors.New("Must be in valid email format.")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Email",
		Validate: validate,
	}

	res, err := prompt.Run()

	if err != nil {
		return "", err
	}
	return res, nil
}

// companyPrompt prompts the user for  organization
func companyPrompt() (string, error) {
	validate := func(input string) error {
		if len(input) <= 5 {
			return errors.New("Invalid length. Must be longer than 5 characters.")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Organization",
		Validate: validate,
	}

	res, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return res, nil
}
