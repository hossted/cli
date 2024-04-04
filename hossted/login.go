package hossted

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Login() error {
	fmt.Println("login")
	email, err := emailPrompt()
	if err != nil {
		return err
	}
	fmt.Println("logging in with email", email)
	postRequest(email)
	return nil
}

func postRequest(email string) {
	authToken := "FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA"
	payloadStr := fmt.Sprintf(`{"email": "%s"}`, email)
	url := "https://api.dev.hossted.com/v1/instances/cli/login"
	payload := []byte(payloadStr)
	fmt.Println(string(payload))

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Add Authorization header with Basic authentication
	authHeader := fmt.Sprintf("Basic %s", []byte(authToken))
	req.Header.Set("Authorization", authHeader)
	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and display the response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}
	response_byte, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Response:", string(response_byte))
	// Save response in file
	saveResponse(email, response_byte)

}

func saveResponse(email string, result []byte) {
	homeDir, err := os.UserHomeDir()
	folderPath := filepath.Join(homeDir, ".hossted")
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	// Check and create .hossted directory in user's home
	setDir(folderPath)
	// Save the response in file.
	updateFile(folderPath, email, result)

}

func updateFile(folderPath string, email string, result []byte) {
	filename := strings.Replace(email, "@", "-", -1)
	filename += ".json"
	filePath := filepath.Join(folderPath, filename)
	// Create the file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating config file:", err)
			return
		}
		file.Write(result)
		file.Close()
	} else if os.IsExist(err) {
		// truncate and update the existing file
		fmt.Println("File exists")
		file, err := os.OpenFile(filePath, os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("Error opening config file:", err)
			return
		}
		err = file.Close()
		if err != nil {
			fmt.Println("Error while saving in config file:", err)
			return
		}
		writeToFile(filePath, result)
	} else {
		fmt.Println("Unable to check the config file", err)
		return
	}
}
func setDir(folderPath string) {
	// Check if the folder exists in the home directory
	// Create the folder if it doesn't exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}
		fmt.Println("Folder created at:", folderPath)
	}

}

func writeToFile(filePath string, result []byte) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	_, err = file.Write(result)
	if err != nil {
		fmt.Println("Error writing in file:", err)
		return
	}
	err = file.Close()
	if err != nil {
		fmt.Println("Error while saving in file:", err)
		return
	}
}
