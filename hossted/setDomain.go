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
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hossted/cli/hossted/service"
	"github.com/hossted/cli/hossted/service/common"
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

		tr, err := common.GetTokenResp()
		if err != nil {
			return err
		}
		err = service.SendEvent("success", "Hossted Platform Domain set successfully", tr.AccessToken)
		if err != nil {
			return err
		}

	} else {
		if !HasContainerRunning() {
			fmt.Println("The application still in configuration")
			os.Exit(0)
		}

		command := "domain"

		err = CheckCommands(app, command)
		if err != nil {
			return fmt.Errorf("\n\n%w", err)
		}

		err = AddDomainToMotd(domain)
		if err != nil {
			return err
		}

		check := verifyInputFormat(domain, "domain")
		if !check {
			return fmt.Errorf("Invalid domain input. Expecting domain name (e.g. example.com).\nInput - %s\n", domain)
		}

		// Get .env file and appDir
		appConfig, err := config.GetAppConfig(app)
		if err != nil {
			return err
		}
		appDir := appConfig.AppPath
		envPath, err := getAppFilePath(appConfig.AppPath, ".env")
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

// ChangeMOTD changes the content of the MOTD file, to match the set domain changes
// TODO: print status
// TODO: Allow domain to be something other than .com by changing the regex patten
func ChangeMOTD(domain string) error {

	filepath := "/etc/motd"
	b, err := readProtected(filepath)
	if err != nil {
		return fmt.Errorf("Can't read the /etc/motd file. Please check - %s and contact administrator.\n%w\n", filepath, err)
	}
	content := string(b)

	// Currently only .com is supported. Looking for line like
	// Your ^[[01;32mgitbucket^[[0m is available under ^[[01;34m http://3.215.23.221.c.hossted.com ^[[0m
	re, err := regexp.Compile(`(.*available under\s*.*https?:\/\/)(.*\.com)(.*)`)
	if err != nil {
		return err
	}

	matches := re.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 {
		if len(matches[0]) == 4 {
			new := matches[0][1] + domain + matches[0][3]
			content = strings.Replace(content, matches[0][0], new, 1) // Replace the containing new with new string
		}
	} else {
		return errors.New("No matching pattern in /etc/motd. Please check.\n")
	}

	// Write back to file
	err = writeProtected(filepath, []byte(content))
	if err != nil {
		return err
	}

	return nil
}
