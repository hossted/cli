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

	// Check command
	err = CheckCommands(app.AppName, "httpopen")
	if err != nil {
		return err
	}

	// Get appPath
	appConfig, err := config.GetAppConfig(app.AppName)
	if err != nil {
		return err
	}
	name := appConfig.AppName // app name
	path := appConfig.AppPath // app path. e.g. /opt/gitbucket
	if err != nil {
		return err
	}

	fmt.Println("Some sed statement")

	fmt.Printf("App Path: %s\n", path)
	err = stopTraefik(path)
	if err != nil {
		return err
	}

	err = dockerUp(path)
	if err != nil {
		return err
	}

	fmt.Printf("Service Restarted - %s\n", name)

	return nil
}
