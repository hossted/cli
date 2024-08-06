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

type authResp struct {
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

	fmt.Printf("User Code: %s\n", authResp.UserCode)
	fmt.Printf("Verification URL Complete: %s\n", authResp.VerificationURIComplete)

	return nil
}

func postRequest(develMode bool) (authresp authResp, err error) {

	var clientID, hosstedAuthUrl string

	// Override values in development mode
	if develMode {
		clientID = common.HOSSTED_DEV_CLIENT_ID
		hosstedAuthUrl = common.HOSSTED_DEV_AUTH_URL + "/device/authorize"
		fmt.Printf("devel mode: true \nclientID: %s\nhosstedAuthUrl: %s\n", clientID, hosstedAuthUrl)
		fmt.Println("------------------------------------------------------------------------------")

	} else {
		clientID = common.HOSSTED_CLIENT_ID
		hosstedAuthUrl = common.HOSSTED_AUTH_URL + "/device/authorize"
		//fmt.Printf("production mode:\nclientID: %s\nhosstedAuthUrl: %s\n", clientID, hosstedAuthUrl)
	}

	// Debugging prints
	if hosstedAuthUrl == "" {
		return authResp{}, fmt.Errorf("hosstedAuthUrl is not set")
	}
	if clientID == "" {
		return authResp{}, fmt.Errorf("clientID is not set")
	}

	data := url.Values{}
	data.Set("client_id", clientID)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, hosstedAuthUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return authResp{}, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return authResp{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return authResp{}, err
	}

	if resp.StatusCode != 200 {
		return authResp{}, fmt.Errorf("Registration Failed, Error %s", string(body))
	}

	var authResp authResp
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return authResp, err
	}

	err = saveResponse(body)
	if err != nil {
		return authResp, err
	}

	return authResp, nil
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
