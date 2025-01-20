package hossted

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// SetDomain set the domain for different apps
// TODO: check whether the function is generic for different apps. Divide to different cases if not.
// TODO: check error for sed command
func SetDomain(env, app, domain string) error {
	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("unable to get activateType from config, error:%w", err)
	}

	//check ActivateType in config
	if config.ActivateType == "k8s" {
		var kubeconfig string
		path, ok := os.LookupEnv("KUBECONFIG")
		if ok {
			kubeconfig = path
		} else {
			kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}

		conf, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// Create the clientset
		clientset, err := kubernetes.NewForConfig(conf)
		if err != nil {
			panic(err.Error())
		}

		namespace := app

		// List all Ingresses in the namespace
		ingressList, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing Ingresses: %v", err)
		}

		// Assuming there's only one Ingress, pick the first one
		ingress := &ingressList.Items[0]

		if len(ingress.Spec.Rules) > 0 {
			if domain == ingress.Spec.Rules[0].Host {
				fmt.Println("domain is same as the existing one, hence no update. Exiting.")
				return nil
			}
			ingress.Spec.Rules[0].Host = domain

		} else {
			fmt.Println("No rules found in Ingress spec")
		}

		// Update the Ingress
		updatedIngress, err := clientset.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("unable to set domain for %s ingress, error:%w", app, err)
		}
		fmt.Printf("domain %s set for %s ingress\n", updatedIngress.Spec.Rules[0].Host, app)

		// err = service.SendEvent("success", "Hossted Platform Domain set successfully", common.HOSSTED_AUTH_TOKEN)
		// if err != nil {
		// 	return err
		// }
		fmt.Println("Hossted Platform Domain set successfully")

	} else {

		err = CheckHosstedAuthFiles()
		if err != nil {
			fmt.Println("Please run hossted activate -t compose, to activate the vm")
			os.Exit(1)
		}
		if !HasContainerRunning() {
			fmt.Println("The application still in configuration")
			os.Exit(0)
		}

		command := "domain"

		err = CheckCommands(app, command)
		if err != nil {
			return fmt.Errorf("\n\n%w", err)
		}

		err = ChangeMOTD(domain)
		if err != nil {
			return err
		}

		check := verifyInputFormat(domain, "domain")
		if !check {
			return fmt.Errorf("invalid domain input. Expecting domain name (e.g. example.com). input - %s", domain)
		}

		appDir := "/opt/" + app
		envPath, err := getAppFilePath(appDir, ".env")
		if err != nil {
			return err
		}

		// Use sed to change the domain
		// TODO: check if the line really exists in the file first
		fmt.Println("Changing settings...")
		text := fmt.Sprintf("s/(PROJECT_BASE_URL=)(.*)/\\1%s/", domain)
		cmd := exec.Command("sudo", "sed", "-i", "-E", text, envPath)
		_, err = cmd.Output()
		if err != nil {
			return err
		}

		// Try command
		fmt.Println("Stopping traefik...")
		err = stopTraefik(appDir)
		if err != nil {
			return err
		}

		err = dockerUp(appDir)
		if err != nil {
			return err
		}

		fmt.Printf("Service Restarted - %s\n", app)

		//send activity log about the command
		uuid, err := GetHosstedUUID(config.UUIDPath)
		if err != nil {
			return err
		}
		fullCommand := "hossted set domain " + fmt.Sprint(domain)
		options := `{"domain":"` + fmt.Sprint(domain) + `"}`
		typeActivity := "set_domain"

		sendActivityLog(env, uuid, fullCommand, options, typeActivity)
		return nil
	}
	return nil
}

func ChangeMOTD(domain string) error {
	filepath := "/etc/motd"

	// Read the file
	b, err := readProtected(filepath)
	if err != nil {
		return fmt.Errorf("unable to read the /etc/motd file. Please check %s and contact administrator: %w", filepath, err)
	}
	content := string(b)

	// Match and update any URL that starts with https:// followed by a domain
	re := regexp.MustCompile(`https://[\w\.\-]+\.\w+`)
	updatedContent := re.ReplaceAllString(content, fmt.Sprintf("https://%s", domain))

	if updatedContent == content {
		return errors.New("no matching pattern found in /etc/motd. Please ensure the content is formatted correctly")
	}

	// Write the updated content back to the file
	err = writeProtected(filepath, []byte(updatedContent))
	if err != nil {
		return fmt.Errorf("failed to write to the /etc/motd file: %w", err)
	}

	return nil
}

// CheckHosstedAuthFiles checks if the files ~/.hossted/auth.json and ~/.hossted/authresp.json exist.
func CheckHosstedAuthFiles() error {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Define the file paths
	authFilePath := filepath.Join(homeDir, ".hossted", "auth.json")
	authRespFilePath := filepath.Join(homeDir, ".hossted", "authresp.json")

	// Check if auth.json exists
	if _, err := os.Stat(authFilePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", authFilePath)
	}

	// Check if authresp.json exists
	if _, err := os.Stat(authRespFilePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", authRespFilePath)
	}

	// Both files exist
	return nil
}
