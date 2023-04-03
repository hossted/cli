package hossted

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Update(env string) error {
	// Define the GitHub repository information
	repoOwner := "hossted"
	repoName := "cli"

	// Get the latest release information from the GitHub API
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the release information to get the download URL for the asset file
	var releaseInfo struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	err = json.NewDecoder(resp.Body).Decode(&releaseInfo)
	if err != nil {
		return err
	}
	if len(releaseInfo.Assets) == 0 {
		return errors.New("no assets found in release")
	}
	assetURL := ""
	for _, asset := range releaseInfo.Assets {
		if strings.Contains(asset.Name, "linux") {
			assetURL = asset.BrowserDownloadURL
			break
		}
	}
	if assetURL == "" {
		return errors.New("no suitable asset found in release")
	}

	// Download the tarball
	resp, err = http.Get(assetURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Extract the tarball
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		panic(err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if strings.Contains(hdr.Name, "hossted") {

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
			// Create the new hossted binary file with the same permissions as the previous file
			newFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
			if err != nil {
				panic(err)
			}
			defer newFile.Close()

			// Copy the contents of the file to the new file
			_, err = io.Copy(newFile, tr)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("hossted updated successfully")
	return nil
}
