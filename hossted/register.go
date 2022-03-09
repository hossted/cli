package hossted

import (
	"embed"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"fmt"

	"github.com/manifoldco/promptui"
)

var (
	//go:embed templates
	templates embed.FS
)

// RegisterUsers updates email, organization in the yaml file, if successful.
// Also it will get the uuid of the machine, and environment from $HOSSTED_ENV.
// It will then send an request to the API server, to get the JWT sessin token.
// TODO: Use the original values as default.
func RegisterUsers() error {

	config, _ := GetConfig() // Ignore error

	// Prompt user for input
	email, _ := emailPrompt()
	organization, _ := organizationPrompt()

	// Get uuid, env. Env default to be dev, if env varible
	env := GetHosstedEnv()
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}

	// Send request
	response, err := registerRequest(email, organization, uuid, env)
	if err != nil {
		return err
	}

	// TODO: Check response status
	jwt := strings.TrimSpace(response.jwt)
	endpoint := strings.TrimSpace(response.url)

	// Assign back to config object
	config.Email = email
	config.Organization = organization
	config.SessionToken = jwt
	config.EndPoint = endpoint

	// Write back to file
	err = WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}

	fmt.Println(fmt.Sprintf("Updated config. Registered User - [%s - %s]\n", email, organization))
	fmt.Println(fmt.Sprintf("Please visit the dashboard link - %s\n", endpoint))
	return nil
}

// registerRequest sends register request based on the input, env/email/organization, etc..
// TODO: Set BearToken to env variable
func registerRequest(email, organization, uuid, env string) (RegisterResponse, error) {

	var response RegisterResponse

	// Construct param map for input params
	params := make(map[string]string)
	params["email"] = email
	params["organization"] = organization
	params["uuid"] = uuid

	req := HosstedRequest{
		// Endpoint env needs to replace in runtime for url parse to work. Otherwise runtime error.
		// currently disabling env support - LK
		EndPoint:     "https://app.hossted.com/api/register",
		Environment:  env,
		Params:       params,
		BearToken:    "Basic FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA",
		SessionToken: "",
	}

	fmt.Println("Registering user. Please wait a second.")
	resp, err := req.SendRequest()
	if err != nil {
		return response, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return response, fmt.Errorf("Failed to parse JSON. %w", err)
	}

	// Check if the sessionToken is null
	if strings.TrimSpace(response.jwt) == "" {
		return response, fmt.Errorf("Empty Session Token. Please check the api requests. %s.\n", resp)
	}

	return response, nil
}

// emailPromp prompt the user for email
func emailPrompt() (string, error) {

	// Regex for email checking
	re := regexp.MustCompile(`\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+`)

	validate := func(input string) error {
		if len(input) <= 5 {
			return errors.New("Invalid length. Must be larger than 5")
		}

		if !re.MatchString(input) {
			return errors.New("Must be in valid email format")
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
		if len(input) < 3 {
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
