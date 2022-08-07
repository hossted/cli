/*
Copyright © 2022 Lior Kesos lior@hossted.com

*/
package cmd

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/hossted/cli/hossted"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	VERSION     = "dev" // Update during build time
	ENVIRONMENT = "dev" // Update during build time
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:     "hossted",
		Version: VERSION,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var err error
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	_, err = checkConfigFilePath()

	if err != nil {
		fmt.Println(err)
	}

	// Set greetings
	greetings := fmt.Sprintf(`
Hossted CLI %s - for help please contact us at support@hossted.com

Usage:
  hossted [command]

Available Commands:

  |-------------------------------+-------------------------+------------------------------------------------------|
  | Commands                      | Alias                   | Description                                          |
  |-------------------------------+-------------------------+------------------------------------------------------|
  | register                      | hossted r               | Register your application with the hossted ecosystem |
  | set auth false/true           | hossted s a false/true  | Enable / disable HTTP Basic Auth                     |
  | set remote-support false/true | hossted s r false/true  | Enable / disable ssh access for Hossted support team |
  | set domain <domain>           | hossted s d example.com | Set a custom domain                                  |
  | help                          | hossted help            | Help about any command                               |
  | logs                          | hossted log             | View application logs                                |
  | ps                            | hossted ps              | docker-compose ps of the application                 |
  | version                       | hossted version         | Get the version of the hossted CLI program           |
  |-------------------------------+-------------------------+------------------------------------------------------|

Flags:
  -h, --help      help for hossted
  -t, --toggle    Help message for toggle
  -v, --version   version for hossted

Use "hossted [command] --help" for more information about a command.

`, VERSION)
	rootCmd.SetHelpTemplate(greetings)
}

func initConfig() {

	// Not allowed to change cfg path anyway
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	folder := fmt.Sprintf("%s/%s", home, ".hossted")
	viper.AddConfigPath(folder)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	return
}

// checkConfigFilePath checks if the ~/.hossted/config.yaml is created under home folder
// Create it if it doesnt exist. Will create folder recursively. Also it will init the config file yaml.
// Also it will check for the new fields in the Config Struct and write to config.yaml again
// TODO: Use the util function one instead
func checkConfigFilePath() (string, error) {

	// Get config path, and .hossted folder. Under user home
	cfgPath, err := hossted.GetConfigPath()
	if err != nil {
		return "", err
	}

	folder := path.Dir(cfgPath)

	if _, err := os.Stat(cfgPath); err != nil {

		// Create directory if not exists
		if _, err := os.Stat(folder); err != nil {
			os.MkdirAll(folder, os.ModePerm)
		}

		// Create file
		f, err := os.OpenFile(cfgPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return "", err
		}
		defer f.Close()

		// Write init config from template
		var config hossted.Config
		err = hossted.WriteConfig(f, config) // empty config
		if err != nil {
			return "", err
		}

		fmt.Printf("\nNo existing config file. \nNew config file is created  - %s \n\n", cfgPath)

	} else {
		// Normal case
		// Do nothing. config.yaml exists
	}

	// Check for new fields, In case new config is available.
	// Write back to file with the original content and new fields (if any).
	config, _ := hossted.GetConfig()

	// Update App related info anyway
	apps, err := hossted.GetAppInfo()
	if err != nil {
		return cfgPath, err
	}

	// Assume Single Application for now
	config.UUIDPath, err = hossted.GetUUIDPath()
	if err != nil {
		return "", err
	}
	// Populate the config with hosts UUID
	var HostUUIDObj, readFileErr = ioutil.ReadFile(config.UUIDPath)
	config.HostUUID = string(HostUUIDObj)

	if err != nil {
		return "", readFileErr
	}
	config.Applications = apps

	// Just write back to the config file, new fields should be written as well
	err = hossted.WriteConfigWrapper(config)
	if err != nil {
		return cfgPath, fmt.Errorf("Can not write the new config to config.yaml. %w", err)
	}

	return cfgPath, nil
}
