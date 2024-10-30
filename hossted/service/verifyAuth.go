package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Verify auth tokens
func VerifyAuth(develMode bool) error {
	// Read authresp.json file
	authRes, err := readAuthRespFile()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("\033[33mExisting auth token not found. Proceeding with login...\033[0m")
			if loginErr := Login(develMode); loginErr != nil {
				return fmt.Errorf("\033[31mLogin failed: %v\033[0m", loginErr)
			}
			return nil
		}
		return fmt.Errorf("\033[31mError reading auth response file: %v\033[0m", err)
	}

	// Check if access token is expired
	if checkAccessTokenExpiration(authRes) {
		// If access token expired, check refresh token expiration
		if checkRefreshTokenExpiration(authRes) {
			fmt.Println("\033[31mBoth tokens expired. Logging in again...\033[0m")
			if loginErr := Login(develMode); loginErr != nil {
				return fmt.Errorf("\033[31mLogin failed: %v\033[0m", loginErr)
			}
		} else {
			// Get new access token using existing refresh token
			if refreshErr := refreshAccessToken(develMode, authRes); refreshErr != nil {
				return fmt.Errorf("\033[31mError refreshing access token: %v\033[0m", refreshErr)
			}
			fmt.Println("\033[32mAccess token refreshed successfully.\033[0m")
		}
	}

	return nil
}

func readAuthRespFile() (authResp, error) {
	var authRespFileData authResp
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return authRespFileData, fmt.Errorf("\033[31m%v\033[0m", err)
	}
	folderPath := filepath.Join(homeDir, ".hossted")
	fileData, err := os.ReadFile(folderPath + "/" + "authresp.json")
	if err != nil {
		return authRespFileData, err // no error color added to allow caller to check os.IsNotExist
	}

	err = json.Unmarshal(fileData, &authRespFileData)
	if err != nil {
		return authRespFileData, fmt.Errorf("\033[31m%v\033[0m", err)
	}

	return authRespFileData, nil
}

func checkAccessTokenExpiration(authres authResp) bool {
	expired := true
	// Current timestamp in seconds
	currentTimestamp := time.Now().Unix()

	accessTokenTS := authres.AccessTokenTimestamp
	accessTokenExpiry := authres.ExpiresIn

	expirationTime := accessTokenTS + int64(accessTokenExpiry)

	// Check if the token has expired
	if currentTimestamp < expirationTime {
		expired = false
	} else {
		fmt.Println("\033[31maccess token has expired, going to refresh access token\033[0m")
	}
	return expired
}

func checkRefreshTokenExpiration(authres authResp) bool {
	expired := true
	// Current timestamp in seconds
	currentTimestamp := time.Now().Unix()

	refreshTokenTS := authres.RefreshTokenTimestamp
	refreshTokenExpiry := authres.RefreshTokenExpiresIn

	expirationTime := refreshTokenTS + int64(refreshTokenExpiry)

	// Check if the token has expired
	if currentTimestamp >= expirationTime {
		fmt.Println("\033[31mrefresh token has also expired\033[0m")
	}
	return expired
}
