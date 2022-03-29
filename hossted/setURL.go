package hossted

import (
	"fmt"
	"os/exec"
)

// SetURL set the url for different apps
// TODO: check whether the function is generic for different apps. Divide to different cases if not.
// TODO: check error for sed command
// TODO: restart app
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

	check := verifyInputFormat(url, "url")
	if !check {
		return fmt.Errorf("Invalid url input. Expecting domain name (e.g. example.com).\nInput - %s\n", url)
	}

	// Get .env file
	appConfig, err := config.GetAppConfig(app)
	if err != nil {
		return err
	}
	appPath, err := getAppFilePath(appConfig.AppPath, ".env")
	if err != nil {
		return err
	}

	// Use sed to change the url
	// TODO: check if the line really exists in the file first
	text := fmt.Sprintf("s/(PROJECT_BASE_URL=)(.*)/\\1%s/", url)
	cmd := exec.Command("sudo", "sed", "-i", "-E", text, appPath)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

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
