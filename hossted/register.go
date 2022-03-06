package hossted

import (
	"embed"
	"errors"
	"regexp"

	"fmt"

	"github.com/manifoldco/promptui"
)

var (
	//go:embed templates
	templates embed.FS
)

// RegisterUsers updates email, organization, etc,.. in the yaml file
// TODO: Use the original values as default.
func RegisterUsers() error {

	config, _ := GetConfig() // Ignore error
	_ = config

	// Prompt user for input
	email, _ := emailPrompt()
	organization, _ := organizationPrompt()

	// Assign back to config object
	config.Email = email
	config.Organization = organization

	// Send request

	// Write back to file
	err := WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check.")
	}

	fmt.Println(fmt.Sprintf("Updated config. Registered User - [%s - %s]", email, organization))
	return nil
}

// Send request
func registerRequest(email, organization, uuid string) error {

	// Construct param map for input params
	params := make(map[string]string)
	params["email"] = email
	params["organization"] = organization
	params["uuid"] = uuid

	req := HosstedRequest{
		EndPoint:     "https://app.dev.hossted.com/api/register",
		Environment:  "dev",
		Params:       params,
		BearToken:    "Basic FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA",
		SessionToken: "",
	}
	resp, err := req.SendRequest()
	if err != nil {
		return err
	}
	fmt.Println(resp)
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

// organizationPrompt prompts the user for organization
func organizationPrompt() (string, error) {
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
