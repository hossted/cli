package hossted

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// GetConfigPath gets the pre-defined config path. ~/.hossted/config.yaml
func GetConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	configPath := path.Join(home, ".hossted", "config.yaml")
	return configPath, nil
}

// GetConfigPath gets the config object
// TODO: Check which field is missing
func GetConfig() (Config, error) {
	var config Config
	cfgPath, err := GetConfigPath()
	if err != nil {
		return config, err
	}

	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}

	// Check if all the fields are set
	// TODO: Check which field is missing
	if (config.Email == "") || (config.Organization == "") || (config.UserToken == "") {
		return config, fmt.Errorf("One of the fields is null")
	}

	return config, nil
}

// WriteConfigWrapper is a wrapper function to call the underlying io.Writer function
func WriteConfigWrapper(config Config) error {

	// Get config path, and .hossted folder. Under user home
	cfgPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	folder := path.Dir(cfgPath)

	if _, err := os.Stat(cfgPath); err != nil {

		// Create directory if not exists
		if _, err := os.Stat(folder); err != nil {
			os.MkdirAll(folder, os.ModePerm)
		}

		fmt.Printf("\nNo existing config file. \nNew config file is created  - %s \n\n", cfgPath)

		return err
	}

	// Create file
	f, err := os.OpenFile(cfgPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = WriteConfig(f, config) // empty config
	if err != nil {
		return err
	}

	return nil
}

// WriteConfig writes the config to the config file (~/.hossted/config.yaml)
func WriteConfig(w io.Writer, config Config) error {

	// Read Template
	t, err := template.ParseFS(templates, "templates/config.tmpl")
	if err != nil {
		return err
	}

	// Write to template
	err = t.Execute(w, config)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(w)
	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// GetHosstedEnv gets the value of the env variable HOSSTED_ENV. Support dev/prod only.
// If it is not set, default as dev
func GetHosstedEnv() string {
	env := strings.TrimSpace(os.Getenv("HOSSTED_ENV"))
	switch env {
	case "dev":
		env = "dev"
	case "prod":
		env = "prod"
	case "":
		fmt.Printf("Environment variable (HOSSTED_ENV) is not set.\nUsing dev instead.\n")
		env = "dev"
	default:
		fmt.Printf("Only dev/prod is supported for (HOSSTED_ENV).\nUsing dev instead.\n")
		env = "dev"
	}
	return env
}
