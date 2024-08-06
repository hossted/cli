package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hossted/cli/hossted/service/common"
)

type AuthResp struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

func Login(develMode bool) error {
	authResp, err := postRequest(develMode)
	if err != nil {
		return err
	}

	fmt.Printf("User Code: %s\n", authResp)
	return nil
}

func postRequest(develMode bool) (usercode string, err error) {

	var clientID, hosstedAuthUrl string

	// Override values in development mode
	if develMode {
		fmt.Printf("devel mode clientID\n hosstedAuthUrl: %s\n%s\n", clientID, hosstedAuthUrl)
		clientID = common.HOSSTED_DEV_CLIENT_ID
		hosstedAuthUrl = common.HOSSTED_DEV_AUTH_URL
	}

	clientID = common.HOSSTED_CLIENT_ID
	hosstedAuthUrl = common.HOSSTED_AUTH_URL

	data := url.Values{}
	data.Set("client_id", clientID)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, hosstedAuthUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Registration Failed, Error %s", string(body))
	}

	var authResp AuthResp
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return "", err
	}

	err = saveResponse(body)
	if err != nil {
		return "", err
	}

	return authResp.UserCode, nil
}

func saveResponse(data []byte) error {
	homeDir, err := os.UserHomeDir()

	folderPath := filepath.Join(homeDir, ".hossted")
	if err != nil {
		return err
	}

	err = os.WriteFile(folderPath+"/"+"auth.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
