package hossted

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

	// Exit function if no config path.
	if _, err := os.Stat(cfgPath); err != nil {
		return config, fmt.Errorf("Can not open config file - %s. Please check.\n%w", cfgPath, err)
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
	// TODO: Check which field is missing. May be add UserToken back for checking
	if config.Email == "" {
		return config, errors.New("One of the fields [Email] is null.\nPlease call `hossted register` first.\n")
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
		// fmt.Printf("Environment variable (HOSSTED_ENV) is not set.\nUsing dev instead.\n")
		env = "dev"
	default:
		fmt.Printf("Only dev/prod is supported for (HOSSTED_ENV).\nUsing dev instead.\n")
		env = "dev"
	}
	return env
}

// GetHosstedUUID gets the uuid from the file /opt/linnovate/run/uuid.txt
// Return error if it's empty or file not exists
// TODO: Ask whether it's a request to get the uuid or somewhat being saved to the VM during creation
func GetHosstedUUID(path string) (string, error) {
	var uuid string
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return uuid, err
	}
	uuid = strings.TrimSpace(string(b))
	if uuid == "" {
		return uuid, fmt.Errorf("uuid is empty. Please check the content of the file - %s.", path)
	}
	return uuid, nil
}

// GetAppInfo gets the application related information from predefined path /opt/linnovate/run/software.txt
// Returns the App name, and the corresponding path. e.g. Linnovate-AWS-wikijs -> wikijs
// TODO Assume single application for now
// TODO Remove dummy app (for dropdown select demo)
func GetAppInfo() ([]ConfigApplication, error) {
	var (
		appName string // Application name, e.g. wikijs
		appPath string // Application folder, e.g. /opt/wikijs
		apps    []ConfigApplication
	)
	path := "/opt/linnovate/run/software.txt" // Predefined path. Assume single line
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return apps, fmt.Errorf("Can not open %s. Please check.\n%w", path, err)
	}
	text := string(b)

	// Assume single line only
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if len(lines) >= 1 {
		text = lines[0] // First line only
	}

	// Grep cloud env and app name with regex
	re := regexp.MustCompile(`\w*\-(\w*)\-(\w*)`) // e.g. Linnovate-AWS-wikijs
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 3 {
		cloudEnv := matches[1] // e.g. AWS
		_ = cloudEnv
		appName = matches[2]
	}
	appName = strings.ToLower(strings.TrimSpace(appName))
	if appName == "" {
		return apps, fmt.Errorf("Empty appName. Please check the file - %s\n%w", path, err)
	}

	// Check if path exists
	appPath = filepath.Join("/opt", appName)
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return apps, fmt.Errorf("App path does not exists - %s. Please check.\n%w", appPath, err)
	}

	app := ConfigApplication{
		AppName: appName,
		AppPath: appPath,
	}
	apps = append(apps, app)

	// TODO: Demo purpose. Remove later
	demo := ConfigApplication{
		AppName: "demoapp",
		AppPath: "/tmp/demoapp",
	}
	apps = append(apps, demo)

	return apps, nil
}

// updateEndpointEnv replace the place holder with the environment specified
// TODO: Review later. Now only use prod link.
func updateEndpointEnv(endpoint, env string) string {
	endpoint = strings.ReplaceAll(endpoint, "__ENV__", env)
	return endpoint
}

// verifyInputFormat verify different types of user input like, email, url, etc..
func verifyInputFormat(in, format string) bool {

	// Reference: https://stackoverflow.com/questions/10306690/what-is-a-regular-expression-which-will-match-a-valid-domain-name-without-a-subd
	if format == "url" {

		// Replace https and http
		if strings.HasPrefix(in, "https://") {
			in = strings.Replace(in, "https://", "", 1)
		} else if strings.HasPrefix(in, "http://") {
			in = strings.Replace(in, "http://", "", 1)
		} else {
			// pass
		}
		fmt.Println(in)
		re := regexp.MustCompile(`^[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,6}\/?$`)
		if re.MatchString(in) {
			return true
		}
	} else {
		panic("Input format is not supported. Please check")
	}
	return false
}

func replaceYamlSetting(b []byte) error {
	fmt.Println(string(b))
	return nil
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
