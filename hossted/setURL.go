package hossted

import (
	"fmt"
	"io/ioutil"
)

// SetURL set the url for different apps
// TODO: check whether the function is generic for different apps. Divide to different cases if not.
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

	// Read file, and replace url
	b, err := ioutil.ReadFile(appPath)
	if err != nil {
		return err
	}

	setting := YamlSetting{
		Pattern:  `(PROJECT_BASE_URL=).*$`,
		NewValue: fmt.Sprintf("$1 %s", url), // New additional space for subgroup match
	}

	newContent, err := replaceYamlSettings(b, setting)
	if err != nil {
		return err
	}

	// Write back result
	err = overwriteFile(appPath, newContent)
	if err != nil {
		return err
	}
	fmt.Printf("Updated config file - %s\n", appPath)

	return nil

}
