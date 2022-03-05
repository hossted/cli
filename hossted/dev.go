package hossted

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/hossted/utils"
)

// For development only
func Dev() error {
	config, err := GetConfig()
	if err != nil {
		return err
	}
	fmt.Println(utils.PrettyPrint(config))

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
