package hossted

import (
	"fmt"
	"os"
)

func SetURL(app, url string) error {
	command := "url"
	config, err := GetConfig()
	if err != nil {
		fmt.Printf("Please call the command `hossted register` first.\n%w", err)
		os.Exit(0)
	}
	err = CheckCommands(app, command)
	if err != nil {
		return fmt.Errorf("\n\n%w", err)
	}
	_ = config

	return nil

}
