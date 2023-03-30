package hossted

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Update(env string) error {
	filePath := "/usr/local/bin/hossted"

	// Get the file info of the existing hossted binary file
	info, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	// Remove the existing hossted binary file
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	// Make the directories for the new hossted binary file if they don't exist
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Download the new hossted binary file
	resp, err := http.Get("https://github.com/hossted/cli/raw/main/bin/linux/hossted")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Create the new hossted binary file with the same permissions as the previous file
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Copy the contents of the new hossted binary file to the new file with the same permissions as the previous file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
	return nil
}
