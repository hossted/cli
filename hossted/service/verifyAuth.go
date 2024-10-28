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
	// read authresp.json file
	authRes, err := readAuthRespFile()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("\033[33mExisting auth token not found, proceeding with login...\033[0m")
			err := Login(develMode)
			if err != nil {
				return fmt.Errorf("\033[31m%v\033[0m", err)
			}
			return nil
		} else {
			return fmt.Errorf("\033[31m%v\033[0m", err)
		}
	}

	// check access_token
	isAccessTokenExpired := checkAccessTokenExpiration(authRes)
	if isAccessTokenExpired {
		// check refresh token expiry
		isRefreshTokenExpired := checkRefreshTokenExpiration(authRes)
		if isRefreshTokenExpired {
			fmt.Println("\033[31mDoing login again....\033[0m")
			err := Login(develMode)
			if err != nil {
				return fmt.Errorf("\033[31m%v\033[0m", err)
			}
			return fmt.Errorf("\033[31mboth access_token and refresh token were expired, please activate again\033[0m")
		} else {
			// get new access_token using the existing refresh_token
			err := refreshAccessToken(develMode, authRes)
			if err != nil {
				return fmt.Errorf("\033[31m%v\033[0m", err)
			}
			fmt.Println("\033[32mRefreshed access token\033[0m")
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
