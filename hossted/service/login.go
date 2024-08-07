package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"os/exec"
	"runtime"

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
	authResp, err := acquireDeviceCode(develMode)
	if err != nil {
		return err
	}

	fmt.Printf("User Code: %s\n", authResp.UserCode)
	fmt.Printf("Verification URL Complete: %s\n", authResp.VerificationURIComplete)
	openBrowser(authResp.VerificationURIComplete)
	// Schedule pollAccessToken after authResp.Interval seconds

	interval := time.Duration(authResp.Interval) * time.Second
	for {
		time.Sleep(interval)
		err := pollAccessToken(develMode, authResp)
		if err != nil {
			log.Printf("Error polling access token: %v\n", err)
		} else {
			log.Println("Access token polled successfully.")
			break // Exit the loop if polling is successful
		}
	}

	return nil
}

func acquireDeviceCode(develMode bool) (authresp authResp, err error) {

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

	err = saveResponse(body, "auth.json")
	if err != nil {
		return authResp, err
	}

	return authResp, nil
}

func saveResponse(data []byte, fileName string) error {
	homeDir, err := os.UserHomeDir()

	folderPath := filepath.Join(homeDir, ".hossted")
	if err != nil {
		return err
	}

	err = os.WriteFile(folderPath+"/"+fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func pollAccessToken(develMode bool, auth authResp) error {
	var clientID, hosstedAuthUrl string

	// Override values in development mode
	if develMode {
		clientID = common.HOSSTED_DEV_CLIENT_ID
		hosstedAuthUrl = common.HOSSTED_DEV_AUTH_URL + "/device/token"

	} else {
		clientID = common.HOSSTED_CLIENT_ID
		hosstedAuthUrl = common.HOSSTED_AUTH_URL + "/device/token"
		//fmt.Printf("production mode:\nclientID: %s\nhosstedAuthUrl: %s\n", clientID, hosstedAuthUrl)
	}

	// Debugging prints
	if hosstedAuthUrl == "" {
		return fmt.Errorf("hosstedAuthUrl is not set")
	}
	if clientID == "" {
		return fmt.Errorf("clientID is not set")
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	data.Set("device_code", auth.DeviceCode)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, hosstedAuthUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Registration Failed, Error %s", string(body))
	}

	var authResp authResp
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return err
	}

	err = saveResponse(body, "authresp.json")
	if err != nil {
		return err
	}
	return nil
}


func openBrowser(url string) error {
    var cmd string
    var args []string

    switch runtime.GOOS {
    case "windows":
        cmd = "rundll32"
        args = append(args, "url.dll,FileProtocolHandler", url)
    case "darwin":
        cmd = "open"
        args = append(args, url)
    case "linux":
        cmd = "xdg-open"
        args = append(args, url)
    default:
        return fmt.Errorf("unsupported platform")
    }

    return exec.Command(cmd, args...).Start()
}