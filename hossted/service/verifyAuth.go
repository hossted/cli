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
		return err
	}

	// check access_token
	isAccessTokenExpired := checkAccessTokenExpiration(authRes)
	if isAccessTokenExpired {
		// check refresh token expiry
		isRefreshTokenExpired := checkRefreshTokenExpiration(authRes)
		if isRefreshTokenExpired {
			fmt.Println("Doing login again....")
			err := Login(develMode)
			if err != nil {
				return err
			}
			return fmt.Errorf("both access_token and refresh token were expired, please activate again")
		} else {
			// get new access_token using the existing refresh_token
			err := refreshAccessToken(develMode, authRes)
			if err != nil {
				return err
			}
			fmt.Println("Refreshed access token")
		}
	}

	return nil
}

func readAuthRespFile() (authResp, error) {
	var authRespFileData authResp
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return authRespFileData, err
	}
	folderPath := filepath.Join(homeDir, ".hossted")
	fileData, err := os.ReadFile(folderPath + "/" + "authresp.json")
	if err != nil {
		return authRespFileData, err
	}

	err = json.Unmarshal(fileData, &authRespFileData)
	if err != nil {
		return authRespFileData, err
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
		fmt.Println("access token has expired, going to refresh access token")
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
	if currentTimestamp < expirationTime {
		expired = false
	} else {
		fmt.Println("refresh token has also expired")
	}
	return expired
}
