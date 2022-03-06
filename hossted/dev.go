package hossted

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

// For development only
func Dev() error {

	params := make(map[string]string)
	req := HosstedRequest{
		EndPoint:     "https://app.dev.hossted.com/api/register",
		Environment:  "dev",
		Params:       params,
		BearToken:    "Basic FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA",
		SessionToken: "",
	}
	res, err := req.Send()
	if err != nil {
		return err
	}
	fmt.Println(res)

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
