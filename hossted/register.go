package hossted

import (
	"embed"
	"encoding/json"
	"errors"
	"os"
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
func RegisterUsers(env string) error {

	config, _ := GetConfig() // Ignore error

	// Prompt user for input
	email, _ := EmailPrompt()

	// Get uuid, env. Env default to be dev, if env varible
	env = GetHosstedEnv(env)
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}

	// Send request
	response, err := registerRequest(email, uuid, env)
	if err != nil {
		return err
	}

	// TODO: Check response status
	jwt := strings.TrimSpace(response.JWT)
	endpoint := strings.TrimSpace(response.URL)

	// Assign back to config object
	config.Email = email
	config.SessionToken = jwt
	config.EndPoint = endpoint

	// Write back to file
	err = WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}
	dashboardUrl := endpoint + "/verify?token=" + config.SessionToken
	fmt.Println(fmt.Sprintf("Updated config. Registered User - [%s]\n", email))
	fmt.Println(fmt.Sprintf("Please visit the dashboard link - %s\n", dashboardUrl))
	return nil
}

// registerRequest sends register request based on the input, env/email/organization, etc..
// TODO: Set BearToken to env variable
func registerRequest(email, uuid, env string) (RegisterResponse, error) {

	var response RegisterResponse

	// Construct param map for input params
	params := make(map[string]string)
	params["email"] = email
	params["uuid"] = uuid

	req := HosstedRequest{
		// Endpoint env needs to replace in runtime for url parse to work. Otherwise runtime error.
		// Old EndPoint:     "https://app.__ENV__hossted.com/api/register",
		EndPoint:     "https://api.__ENV__hossted.com/v1/instances/registerUser",
		Environment:  env,
		Params:       params,
		BearToken:    "Basic FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA",
		SessionToken: "",
		TypeRequest:  "POST",
	}

	fmt.Println("Registering user. Please wait a second...")
	resp, err := req.SendRequest()
	if err != nil {
		return response, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return response, fmt.Errorf("Failed to parse JSON. %w", err)
	}

	// Check if the sessionToken is null
	if strings.TrimSpace(response.JWT) == "" {
		if strings.Contains(response.Message, "already been registered") {

			fmt.Printf("\nInstance is already registered. Please contact support@hossted.com if you'd like to add a user.\n\n")
			os.Exit(0)

		} else {
			fmt.Printf("Empty Session Token. Some kind of error occoured - %s.\n", resp)
			os.Exit(0)
		}
	}

	return response, nil
}

// emailPrompt prompts the user for email
func EmailPrompt() (string, error) {

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
