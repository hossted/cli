package hossted

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/hossted/utils"
)

// For development only
func Dev() error {

	// Construct param struct for input params
	params := make(map[string]string)
	params["uuid"] = "55cdfdae-ce22-4c36-8513-b09df945734a"
	params["email"] = "billy@hossted.com"
	params["organization"] = "asdf"

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
	var response RegisterResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return err
	}
	fmt.Println(utils.PrettyPrint(resp))

	return nil
}

// For development only
func Prompt() (string, error) {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
}
