package hossted

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"

	"context"
	"reflect"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

	// Check if all the fields are set.
	// Ticket: Removed checking for Issue #20.
	// TODO: Check which field is missing. May be add UserToken back for checking
	// if config.Email == "" {
	// 	return config, errors.New("One of the fields [Email] is null.\nPlease call `hossted register` first.\n")
	// }

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

		//fmt.Printf("\nNo existing config file. \nNew config file is created  - %s \n\n", cfgPath)

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
		// fmt.Println(err)
	}

	return nil
}

// GetHosstedEnv determines which env. Support dev/prod only.
// If it is not set, default as dev
func GetHosstedEnv(env string) string {

	switch env {
	case "dev":
		env = "dev"
	case "prod":
		env = "prod"
	case "":
		// fmt.Printf("Environment variable (HOSSTED_ENV) is not set.\nUsing dev instead.\n")
		env = "dev"
	default:
		fmt.Printf("Only dev/prod is supported for env.\nUsing dev instead.\n")
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

// GetAppInfo gets the application related information from predefined path
// /opt/hossted/run/software.txt or /opt/linnovate/run/software.txt
// Returns the App name, and the corresponding path. e.g. Linnovate-AWS-wikijs -> wikijs
// TODO Assume single application for now
// TODO Remove dummy app (for dropdown select demo)
func GetAppInfo() ([]ConfigApplication, error) {
	var (
		appName string // Application name, e.g. wikijs
		appPath string // Application folder, e.g. /opt/wikijs
		apps    []ConfigApplication
	)

	// Predefined path. Assume single line in /opt/hossted/run/software.txt or /opt/linnovate/run/software.txt
	path, err := GetSoftwarePath()
	if err != nil {
		return apps, err
	}

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

	appPath, _ = GetDockerComposeDir()
	if appPath == "" {
		appPath = filepath.Join("/opt", appName)
	}

	// Check if path exists
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return apps, fmt.Errorf("App path does not exists - %s. Please check.\n%w", appPath, err)
	}

	app := ConfigApplication{
		AppName: appName,
		AppPath: appPath,
	}
	apps = append(apps, app)

	return apps, nil
}

// updateEndpointEnv replace the place holder with the environment specified
// TODO: Review later. Now only use prod link.
func updateEndpointEnv(endpoint, env string) string {
	if env == "prod" {
		// prod: https://app.hossted.com/api/register
		endpoint = strings.ReplaceAll(endpoint, "__ENV__", "")
	} else if env == "dev" {
		// dev: https://app.dev.hossted.com/api/register
		endpoint = strings.ReplaceAll(endpoint, "__ENV__", "dev.")
	}
	return endpoint
}

// verifyInputFormat verify different types of user input like, email, url, etc..
func verifyInputFormat(in, format string) bool {

	// Reference: https://stackoverflow.com/questions/10306690/what-is-a-regular-expression-which-will-match-a-valid-domain-name-without-a-subd
	if format == "domain" {

		// Replace https and http
		if strings.HasPrefix(in, "https://") {
			in = strings.Replace(in, "https://", "", 1)
		} else if strings.HasPrefix(in, "http://") {
			in = strings.Replace(in, "http://", "", 1)
		} else {
			// pass
		}
		re := regexp.MustCompile(`^[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,6}\/?$`)
		if re.MatchString(in) {
			return true
		}
	} else {
		panic("Input format is not supported. Please check")
	}
	return false
}

// TODO: Find out what is the deal for the space in value
// TODO: Fix additonal lines for result
func replaceYamlSettings(b []byte, setting YamlSetting) (string, error) {
	var (
		pattern  = setting.Pattern  // regex pattern. e.g. `(PROJECT_BASE_URL=).*`
		value    = setting.NewValue // New values. e.g. "$1abc"
		newLines []string
		result   string // result content of the file
		matched  bool   // Match exactly once only
	)
	_ = value
	content := strings.ReplaceAll(string(b), "\r\n", "\n") // For windows
	lines := strings.Split(content, "\n")

	// Compile regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("Invaid regex. Pattern (%s). %w", pattern, err)
	}

	// For each line, check with the pattern
	for _, line := range lines {

		// Check if match, then replace
		if !matched && re.MatchString(line) { // only once
			new := re.ReplaceAllString(line, value)
			new = strings.ReplaceAll(new, " ", "")
			newLines = append(newLines, new)
			matched = true
		} else {
			newLines = append(newLines, line)
		}

	}
	// If no matched, return error
	if matched == false {
		return "", fmt.Errorf("No matching pattern for [%s]. Please check", pattern)
	}

	// Join back the lines
	result = strings.Join(newLines, "\n")

	return result, nil
}

func overwriteFile(filepath string, content string) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(content)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

// readProtected read the file content with sudo right
// TODO: Add sudo flag
func readProtected(filepath string) ([]byte, error) {

	cmd := exec.Command("sudo", "cat", filepath)
	out, err := cmd.Output()
	if err != nil {
		return []byte{}, fmt.Errorf("Protected file does not exists. Please check - %s.\n%w\n", filepath, err)
	}

	return out, nil
}

// writeProtected write the file content with sudo right
// TODO: Add sudo flag
// TODO: Remove last line break
func writeProtected(path string, b []byte) error {

	// Check if the file exists first
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Protected file does not exist. Please check - %s.\n%w\n", path, err)
	}

	// Write to file
	content := string(b)
	cmd := exec.Command("sudo", "bash", "-c", fmt.Sprintf("echo '%s' > '%s'", content, path))
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func getAppFilePath(base, relative string) (string, error) {
	path := filepath.Join(base, relative)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("file not exists. please check. %w", err)
	}
	return path, nil
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func stopTraefik(appDir string) error {
	fmt.Println("Stopping traefik...")

	// Construct the docker compose down command
	cmd := exec.Command("sudo", "docker", "compose", "down")
	cmd.Dir = appDir // Set the working directory for the command

	// Run the command and capture stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error stopping traefik: %v\nOutput: %s", err, string(output))
	}

	// Print the trimmed output for debugging
	fmt.Println(string(output))
	fmt.Println("Traefik stopped")
	return nil
}

func dockerUp(appDir string) error {
	fmt.Println("Restarting service...")

	// Construct the docker compose command
	cmd := exec.Command("sudo", "docker", "compose", "up", "-d")
	cmd.Dir = appDir // Set the working directory for the command

	// Run the command and capture stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running docker compose: %v\nOutput: %s", err, string(output))
	}

	// Print the trimmed output for debugging
	fmt.Println(string(output))
	return nil
}

// Shell calls the bash shell command in a particular directory
func Shell(appDir string, commands []string) (error, string, string) {

	const ShellToUse = "bash"
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
		sout   []string // List of stdout in string format
		serr   []string // List of stderr in string format
	)

	for _, command := range commands {
		cmd := exec.Command(ShellToUse, "-c", command)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = appDir
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("Can not call Shell Command [%s]. %w\n", command, err), strings.Join(sout, "\n"), strings.Join(serr, "\n")
		}

		// Append stdout and stderr, if any
		if strings.TrimSpace(stdout.String()) == "" {
			sout = append(sout, stdout.String())
		}
		if strings.TrimSpace(stdout.String()) == "" {
			serr = append(serr, stderr.String())
		}
	}

	return nil, strings.Join(sout, "\n"), strings.Join(serr, "\n")
}

// Wrapper for getting current directory
func GetCurrentDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Can not get current directory.")
		return ""
	}
	return pwd
}

// ConvertBool converts the string input to boolean, return error if the input is not true/false
func ConvertBool(in string) (bool, error) {

	var flag bool
	in = strings.ToLower(in)

	if in == "true" {
		flag = true
	} else if in == "false" {
		flag = false
	} else {
		return flag, fmt.Errorf("Only true/false is supported. Input - %s\n", in)
	}

	return flag, nil
}

// trimOuput remove the last (double line breaks) from the string
// usually use before printing out stderr
// TODO: Test why not working
func trimOutput(in string) string {
	s := strings.Replace(in, "\n\n", "\n", -1)
	return s
}

// GetSoftwarePath gets the software related path, it could either be
// /opt/hossted/run/software.txt (Preferred) or /opt/linnovate/run/software.txt
// If neither of that exists, return error
func GetSoftwarePath() (string, error) {
	path := "/opt/hossted/run/software.txt"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// pass
	} else {
		return path, nil
	}
	// Try another path
	path = "/opt/linnovate/run/software.txt"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		//return "", fmt.Errorf("Config file does not exists in both /opt/hossted/run/software.txt or /opt/linnovate/run/software.txt. Please check.\n%w\n", err)
		return "", err
	} else {
		return path, nil
	}

}

// GetUUIDPath is similar to GetSoftwarePath
func GetUUIDPath() (string, error) {
	path := "/opt/hossted/run/uuid.txt"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// pass
	} else {
		return path, nil
	}
	// Try another path
	path = "/opt/linnovate/run/uuid.txt"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("Config file does not exists in both /opt/hossted/run/uuid.txt or /opt/linnovate/run/uuid.txt. Please check.\n%w\n", err)
	} else {
		return path, nil
	}

	return path, nil
}

func GetDockersInfo() (string, error) {

	//fmt.Printf("Start look for dockers\n")

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	if len(containers) == 0 {
		return "", fmt.Errorf("No docker containers found")
	}

	dockers := ""
	var docker Docker
	for _, container := range containers {
		imageName := container.Image

		image, _, err := cli.ImageInspectWithRaw(ctx, imageName)
		if err != nil {
			panic(err)
		}

		docker = Docker{
			ID:        container.ID,
			Image:     container.Image,
			ImageID:   container.ImageID,
			CreatedAt: container.Created,
			Ports:     container.Ports,
			Status:    container.Status,
			Size:      image.Size,
			Names:     container.Names,
			Mounts:    container.Mounts,
			Networks:  reflect.ValueOf(container.NetworkSettings.Networks).MapKeys()[0].String(),
			Tag:       image.RepoTags[0],
		}

		dockerjson, err := json.Marshal(docker)
		if err != nil {
			return "", fmt.Errorf("Error occured during marshaling. Error: %s", err.Error())
		}
		dockers = dockers + string(dockerjson) + ","
	}
	dockers = dockers[:len(dockers)-1]
	dockers = "[" + dockers + "]"
	//fmt.Printf("dockers: %s\n", dockers)

	return string(dockers), nil
}

func sendActivityLog(env, uuid, fullCommand, options, typeActivity string) (activityLogResponse, error) {

	var response activityLogResponse

	user, err := user.Current()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	userName := user.Username
	// Construct param map for input params
	params := make(map[string]string)
	params["uuid"] = uuid
	params["command"] = fullCommand
	params["options"] = options
	params["user_name"] = userName
	params["type"] = typeActivity
	req := HosstedRequest{
		// Endpoint env needs to replace in runtime for url parse to work. Otherwise runtime error.
		//EndPoint:     "https://api.__ENV__hossted.com/v1/instances/dockers",
		EndPoint:     "https://api.hossted.com/v1/instances/activityLog", //"https://api.stage.hossted.com/v1/instances/activityLog", // // ,//"https://api.dev.hossted.com/v1/instances/activityLog", //"http://localhost:3004/v1/activityLog", //,
		Environment:  env,
		Params:       params,
		BearToken:    "Basic y5TXKDY4kTKbFcFtz9aD1pa2irmzhoziKPnEBcA8",
		SessionToken: "",
		TypeRequest:  "PATCH",
	}

	resp, err := req.SendRequest()
	if err != nil {
		return response, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return response, fmt.Errorf("failed to parse JSON. %w", err)
	}

	return response, nil

}

// HasContainerRunning checks if there is a container with the name "trafik" in the list of containers
func HasContainerRunning() bool {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if strings.Contains(container.Names[0], "traefik") {
			return true
		}
	}

	return false
}

// GetDockerComposeDir gets the docker compose directory
func GetDockerComposeDir() (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", nil
	}

	// Get a list of all running containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return "", nil
	}

	// Get the first container name from the list
	if len(containers) == 0 {
		return "", nil
	}
	containerName := containers[0].Names[0]
	// Run docker container inspect command
	inspect, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return "", nil
	}

	// Get label value
	labelValue := inspect.Config.Labels["com.docker.compose.project.working_dir"]
	return labelValue, nil
}

// AddDomainToMotd appends a given domain name to /etc/motd if it does not already exist
func AddDomainToMotd(domain string) error {
	// Open /etc/motd in read mode to check if the domain already exists
	file, err := os.Open("/etc/motd")
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to open /etc/motd for reading: %w", err)
		}
	} else {
		defer file.Close()

		// Check if the domain is already present
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) == domain {
				fmt.Printf("Domain '%s' already exists in /etc/motd\n", domain)
				return nil
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading /etc/motd: %w", err)
		}
	}

	// Open /etc/motd in append mode to add the new domain
	file, err = os.OpenFile("/etc/motd", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open /etc/motd for writing: %w", err)
	}
	defer file.Close()

	// Write the domain name to the file
	if _, err := file.WriteString(domain + "\n"); err != nil {
		return fmt.Errorf("failed to write to /etc/motd: %w", err)
	}

	fmt.Printf("Domain '%s' successfully added to /etc/motd\n", domain)
	return nil
}
