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
	appPath := appConfig.AppPath
	if err != nil {
		return err
	}

	fmt.Println("Some sed statement")
	fmt.Println(app.AppName)

	fmt.Printf("App Path: %s\n", appPath)
	fmt.Println("Stopping traefik...")
	err = stopTraefik(appPath)
	if err != nil {
		return err
	}

	err = dockerUp(appPath)
	if err != nil {
		return err
	}

	fmt.Printf("Service Restarted - %s\n", app)

	return nil
}
