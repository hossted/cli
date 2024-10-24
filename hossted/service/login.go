package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hossted/cli/hossted/service/common"
)

type authLoginResp struct {
	DeviceCode              string `json:"device_code,omitempty"`
	UserCode                string `json:"user_code,omitempty"`
	VerificationURI         string `json:"verification_uri,omitempty"`
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int    `json:"expires_in,omitempty"`
	Interval                int    `json:"interval,omitempty"`
}

type authResp struct {
	AccessToken           string `json:"access_token,omitempty"`
	RefreshToken          string `json:"refresh_token,omitempty"`
	TokenType             string `json:"token_type,omitempty"`
	State                 string `json:"state,omitempty"`
	ExpiresIn             int    `json:"expires_in,omitempty"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in,omitempty"`
	AccessTokenTimestamp  int64  `json:"access_token_timestamp,omitempty"`
	RefreshTokenTimestamp int64  `json:"refresh_token_timestamp,omitempty"`
}

func Login(develMode bool) error {
	loginResp, err := acquireDeviceCode(develMode)
	if err != nil {
		return err
	}

	log.Printf("\033[32mVerification URL: %s\033[0m", loginResp.VerificationURIComplete)
	log.Printf("\033[32mUser Code: %s\033[0m", loginResp.UserCode)
	openBrowser(loginResp.VerificationURIComplete)
	// Schedule pollAccessToken after loginResp.Interval seconds

	interval := time.Duration(loginResp.Interval) * time.Second
	for {
		time.Sleep(interval)
		err := pollAccessToken(develMode, loginResp)
		if err != nil {
			log.Println("\033[33m Please visit the above verification URL to complete sign in and paste the user code\033[0m")
		} else {
			log.Println("\033[32m Access token polled successfully.\033[0m")
			break // Exit the loop if polling is successful
		}
	}

	return nil
}

func acquireDeviceCode(develMode bool) (authloginresp authLoginResp, err error) {

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
		return authLoginResp{}, fmt.Errorf("hosstedAuthUrl is not set")
	}
	if clientID == "" {
		return authLoginResp{}, fmt.Errorf("clientID is not set")
	}

	data := url.Values{}
	data.Set("client_id", clientID)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, hosstedAuthUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return authLoginResp{}, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return authLoginResp{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return authLoginResp{}, err
	}

	if resp.StatusCode != 200 {
		return authLoginResp{}, fmt.Errorf("registration failed, error %s", string(body))
	}

	var loginresp authLoginResp
	err = json.Unmarshal(body, &loginresp)
	if err != nil {
		return loginresp, err
	}

	err = saveResponse(body, "auth.json")
	if err != nil {
		return loginresp, err
	}

	return loginresp, nil
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

func pollAccessToken(develMode bool, loginResp authLoginResp) error {
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
	data.Set("device_code", loginResp.DeviceCode)

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
		return fmt.Errorf("registration failed, error %s", string(body))
	}

	var pollResp authResp
	err = json.Unmarshal(body, &pollResp)
	if err != nil {
		return err
	}

	currentTimestamp := time.Now().Unix()
	pollResp.AccessTokenTimestamp = currentTimestamp
	pollResp.RefreshTokenTimestamp = currentTimestamp

	modifiedData, err := json.Marshal(pollResp)
	if err != nil {
		return fmt.Errorf("error marshalling struct to JSON: %v", err)
	}

	err = saveResponse(modifiedData, "authresp.json")
	if err != nil {
		return err
	}
	return nil
}

func refreshAccessToken(develMode bool, authPollResp authResp) error {
	var clientID, hosstedAuthUrl string

	// Override values in development mode
	if develMode {
		clientID = common.HOSSTED_DEV_CLIENT_ID
		hosstedAuthUrl = common.HOSSTED_DEV_AUTH_URL + "/device/token"

	} else {
		clientID = common.HOSSTED_CLIENT_ID
		hosstedAuthUrl = common.HOSSTED_AUTH_URL + "/device/token"
	}

	// Debugging prints
	if hosstedAuthUrl == "" {
		return fmt.Errorf("hosstedAuthUrl is not set")
	}
	if clientID == "" {
		return fmt.Errorf("clientID is not set")
	}

	data := url.Values{}
	data.Set("refresh_token", authPollResp.RefreshToken)
	data.Set("grant_type", "refresh_token")
	data.Set("state", authPollResp.State)

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
		return fmt.Errorf("refresh access token failed, error %s", string(body))
	}

	var refreshTokenResp authResp
	err = json.Unmarshal(body, &refreshTokenResp)
	if err != nil {
		return err
	}

	currentTimestamp := time.Now().Unix()
	authPollResp.AccessToken = refreshTokenResp.AccessToken
	authPollResp.AccessTokenTimestamp = currentTimestamp

	modifiedData, err := json.Marshal(authPollResp)
	if err != nil {
		return fmt.Errorf("error marshalling struct to JSON: %v", err)
	}

	err = saveResponse(modifiedData, "authresp.json")
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
