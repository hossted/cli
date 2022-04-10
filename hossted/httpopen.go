package hossted

import (
	"fmt"
)

// HttpOpen allows the http access to the frontend
func HttpOpen(input string) error {
	var app ConfigApplication

	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("Please call the command `hossted register` first.\n%w", err)
	}

	app, err = appPrompt(config.Applications, input)
	if err != nil {
		return err
	}

	fmt.Println(string(app.AppName))

	return nil
}
