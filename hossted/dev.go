package hossted

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/manifoldco/promptui"
)

// For development only
func Dev() error {
	m, err := getCommandsMap(GAVAILABLE, AVAILABLE)
	_ = m
	if err != nil {
		return err
	}

	app := "gitbucket"

	// Check command
	err = CheckCommands(app, "auth")
	if err != nil {
		return err
	}

	// fmt.Println(PrettyPrint(m))

	return nil

}

// For development only
func prompt() (string, error) {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
}

// checkCurl is a quick func to check if the api is working as expected.
func checkCurl() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", "https://app.hossted.com/api/register?email=billy%40hossted.com&organization=adf&uuid=55cdfdae-ce22-4c36-8513-b09df945734a", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil

}

func testCommand() error {
	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("Please call the command `hossted register` first.\n%w", err)
	}

	input := "prometheus"
	apps := config.Applications
	var appPath string
	for _, app := range apps {
		if app.AppName == input {
			appPath = app.AppPath
		}
	}
	filepath := fmt.Sprintf("%s/.env", appPath)
	fmt.Println(filepath)

	return nil
}

func testSed() error {
	test := "s/(PROJECT_BASE_URL=)(.*)/\\1some/"
	cmd := exec.Command("sed", "-i", "-E", test, "/tmp/ddd.txt")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}
