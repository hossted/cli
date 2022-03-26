package hossted

import (
	"fmt"
)

func SetURL(app, url string) error {
	command := "url"
	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("Please call the command `hossted register` first.\n%w", err)

	}
	err = CheckCommands(app, command)
	if err != nil {
		return fmt.Errorf("\n\n%w", err)
	}
	_ = config

	return nil

}
