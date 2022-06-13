package hossted

import (
	"fmt"
	"os/exec"
	"regexp"
)

// SetDomain set the domain for different apps
// TODO: check whether the function is generic for different apps. Divide to different cases if not.
// TODO: check error for sed command
func SetDomain(app, domain string) error {
	command := "domain"

	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("Please call the command `hossted register` first.\n%w", err)

	}
	err = CheckCommands(app, command)
	if err != nil {
		return fmt.Errorf("\n\n%w", err)
	}

	check := verifyInputFormat(domain, "domain")
	if !check {
		return fmt.Errorf("Invalid domain input. Expecting domain name (e.g. example.com).\nInput - %s\n", domain)
	}

	// Get .env file and appDir
	appConfig, err := config.GetAppConfig(app)
	if err != nil {
		return err
	}
	appDir := appConfig.AppPath
	envPath, err := getAppFilePath(appConfig.AppPath, ".env")
	if err != nil {
		return err
	}

	// Use sed to change the domain
	// TODO: check if the line really exists in the file first
	fmt.Println("Changing settings...")
	text := fmt.Sprintf("s/(PROJECT_BASE_URL=)(.*)/\\1%s/", domain)
	cmd := exec.Command("sudo", "sed", "-i", "-E", text, envPath)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	// Try command
	fmt.Println("Stopping traefik...")
	err = stopTraefik(appDir)
	if err != nil {
		return err
	}

	err = dockerUp(appDir)
	if err != nil {
		return err
	}

	fmt.Printf("Service Restarted - %s\n", app)

	return nil

}

// ChangeMOTD changes the content of the MOTD file, to match the set domain changes
// TODO: Allow domain to be something other than .com by changing the regex patten
func ChangeMOTD(domain string) error {

	filepath := "/etc/motd"
	b, err := readSth(filepath)
	if err != nil {
		return fmt.Errorf("Can't read the /etc/motd file. Please check - %s and contact administrator.\n%w\n", filepath, err)
	}
	content := string(b)

	// Currently only .com is supported. Looking for line like
	// Your ^[[01;32mgitbucket^[[0m is available under ^[[01;34m http://3.215.23.221.c.hossted.com ^[[0m
	re, err := regexp.Compile(`.*available under (.*https?:\/\/.*\.com).*`)
	if err != nil {
		return err
	}

	matches := re.FindAllStringSubmatch(content, -1)
	fmt.Println(matches)

	return nil
}

// readSth read the file content with sudo right
func readSth(filepath string) ([]byte, error) {

	cmd := exec.Command("cat", filepath)
	out, err := cmd.Output()
	if err != nil {
		return []byte{}, fmt.Errorf("Protected file does not exists. Please check - %s.\n%w\n", filepath, err)
	}

	return out, nil
}
