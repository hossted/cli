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
	err := CheckCommands("prometheus", "url")
	if err != nil {
		return err
	}

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
	cmd := exec.Command("docker-compose", "ps")
	cmd.Dir = "/tmp"
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func testNestedYAML() error {
	fmt.Println("Test")

	return nil
}
