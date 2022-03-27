package hossted

import (
	"fmt"
)

func SetURL(app, url string) error {
	command := "url"
	_, err := GetConfig()
	if err != nil {
		return fmt.Errorf("Please call the command `hossted register` first.\n%w", err)

	}
	err = CheckCommands(app, command)
	if err != nil {
		return fmt.Errorf("\n\n%w", err)
	}

	check := verifyInputFormat(url, "url")
	if !check {
		return fmt.Errorf("Invalid url input. Expecting domain name (e.g.example.com).\nInput - %s\n", url)
	}

	return nil

}
