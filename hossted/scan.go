package hossted

import (
	"fmt"
	
	"os/exec"
	"strings"
)

var (
	images string
)

// scanApps uses trivy to scan the application for vulnerabilities listing their images and checking each one with trivy
func ScanApps(scanType string) error  {
	fmt.Printf("Scanning for %s\n", scanType)
	switch scanType {
	case "images":

		images := listDockerImages()
		// for each image run trivy
		for _, image := range images {
			fmt.Println("Scanning image: " + image)
			scanImage(image)
		}
	}
	return nil
}

func scanImage(image string) error {
	cmd := exec.Command("sudo", "trivy", "image", image)
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

// function that returns a list of docker images
func listDockerImages() []string {

	cmd := exec.Command("sudo", "docker-compose", "images", "-q" )
	cmd.Dir = GetCurrentDirectory()

	// return a list of docker images
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	return strings.Split(string(out), "\n")
}
