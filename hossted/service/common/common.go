package common

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

var (
	LOKI_PASSWORD      = "-"
	LOKI_URL           = "-"
	LOKI_USERNAME      = "-"
	MIMIR_PASSWORD     = "-"
	MIMIR_URL          = "-"
	MIMIR_USERNAME     = "-"
	HOSSTED_API_URL    = "-"
	HOSSTED_AUTH_TOKEN = "-"
	HOSSTED_AUTH_URL   = "-"
	HOSSTED_CLIENT_ID  = "-"
	///////////////////////////
	HOSSTED_DEV_API_URL   = "-"
	MIMIR_DEV_URL         = "-"
	LOKI_DEV_URL          = "-"
	HOSSTED_DEV_CLIENT_ID = "-"
	HOSSTED_DEV_AUTH_URL  = "-"
)

func HttpRequest(method, url, token string, body []byte) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Add Authorization header with Basic authentication
	req.Header.Set("Authorization", "Bearer "+token)
	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// ApiResponse represents the structure of the JSON response from the API
	type ApiResponse struct {
		Success bool        `json:"success"`
		Message interface{} `json:"message"`
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error sending event, errcode: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var apiResponse ApiResponse

	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return err
	}

	if !apiResponse.Success {
		return fmt.Errorf("API response indicates failure: %v\n", apiResponse)
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

type org struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// TokenResponse represents the complete structure of the JSON
type tokenResponse struct {
	AccessToken           string    `json:"access_token"`
	ExpiresIn             int       `json:"expires_in"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresIn int       `json:"refresh_token_expires_in"`
	TokenType             string    `json:"token_type"`
	Orgs                  []org     `json:"orgs"`
	IssuedAt              time.Time `json:"iat"`
	ExpiresAt             time.Time `json:"exp"`
	Issuer                string    `json:"iss"`
}

func GetTokenResp() (tokenResponse, error) {
	var tr tokenResponse
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return tokenResponse{}, err
	}
	folderPath := filepath.Join(homeDir, ".hossted")
	fileData, err := os.ReadFile(folderPath + "/" + "authresp.json")
	if err != nil {
		return tokenResponse{}, err
	}

	// Parse the JSON data into Config struct
	err = json.Unmarshal(fileData, &tr)
	if err != nil {
		return tokenResponse{}, err
	}
	return tr, nil
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Orgs   []org  `json:"orgs"`
	Iat    int64  `json:"iat"`
	Exp    int64  `json:"exp"`
	Iss    string `json:"iss"`
}

func GetOrgs(tokenString string) ([]org, error) {

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {

		return nil, fmt.Errorf("Invalid token format")
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Error decoding payload: %s", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(decodedPayload, &claims); err != nil {
		return nil, fmt.Errorf("Error unmarshalling payload: %s", err)

	}

	return claims.Orgs, nil
}

func OrgUseCases(orgs []org) (orgID string, err error) {

	if len(orgs) == 0 {
		for _, orgID := range orgs {
			fmt.Println("We have just sent the confirmation link registered emailID", ". Once you confirm it, you'll be able to continue the activation.")
			return orgID.ID, nil
		}
	} else if len(orgs) > 1 {
		fmt.Println("You have multiple organisations to choose from:")

		var items []string
		for i, org := range orgs {
			decodedName, err := base64.StdEncoding.DecodeString(org.Name)
			if err != nil {
				return "", err
			}
			items = append(items, fmt.Sprintf("%d: %s", i+1, decodedName))
		}

		prompt := promptui.Select{
			Label: "Select Your Organisation",
			Items: items,
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return "", err
		}

		userOrgName, err := removePrefix(result)
		if err != nil {
			return "", err
		}

		fmt.Printf("OrgName: %s\n", userOrgName)

		var selectedOrgID string

		for _, org := range orgs {
			sorgName, err := base64.StdEncoding.DecodeString(org.Name)
			if err != nil {
				return "", err
			}

			if string(sorgName) == userOrgName {
				selectedOrgID = org.ID
			}
		}

		if selectedOrgID == "" {
			return "", fmt.Errorf("selected organization not found")
		}

		return selectedOrgID, nil

	}

	return "", nil
}

func removePrefix(text string) (string, error) {
	// Define a regular expression to match a number followed by a colon and a space
	regex := regexp.MustCompile(`^\d+:\s+`)

	match := regex.FindStringSubmatch(text)
	if match != nil {
		// Extract the captured prefix (number and colon)
		prefix := match[0]
		return strings.TrimPrefix(text, prefix), nil
	}

	return text, nil
}
