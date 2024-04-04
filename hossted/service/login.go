package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hossted/cli/hossted"
)

func Login() error {

	email, err := hossted.EmailPrompt()
	if err != nil {
		return err
	}
	err = postRequest(email)
	if err != nil {
		return err
	}
	return nil
}

func postRequest(email string) error {
	authToken := os.Getenv("HOSSTED_AUTH_TOKEN")

	payloadStr := fmt.Sprintf(`{"email": "%s"}`, email)

	url := os.Getenv("HOSSTED_API_URL")

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payloadStr)))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Add Authorization header with Basic authentication
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", []byte(authToken)))
	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	saveResponse(email, body)

	return nil
}

func saveResponse(email string, data []byte) error {
	homeDir, err := os.UserHomeDir()

	folderPath := filepath.Join(homeDir, ".hossted")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(folderPath+"/"+email+".json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
