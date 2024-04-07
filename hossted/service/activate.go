package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/gofrs/flock"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/client-go/tools/clientcmd"
)

// Response represents the structure of the JSON response.
type response struct {
	Success bool                `json:"success"`
	OrgIDs  []map[string]string `json:"org_ids"`
	Token   string              `json:"token"`
	Message string              `json:"message"`
}

// ActivateK8s imports Kubernetes clusters.
func ActivateK8s() error {
	// Prompt user for email
	emailID, err := promptEmail()
	if err != nil {
		return err
	}

	// getResponse from reading file in .hossted/email@id.json
	resp, err := getResponse(emailID)
	if err != nil {
		return err
	}

	// handle usecases for orgID selection
	err = useCases(resp, emailID)
	if err != nil {
		return err
	}

	// prompt user for k8s context
	clusterName, err := promptK8sContext()
	if err != nil {
		return err
	}

	fmt.Println("Your cluster name is ", clusterName)

	err = deployOperator(clusterName, emailID)
	if err != nil {
		return err
	}

	return nil
}

func promptEmail() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter your email:",
	}

	emailID, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return emailID, nil

}

func getResponse(emailID string) (response, error) {
	//read file
	homeDir, err := os.UserHomeDir()

	folderPath := filepath.Join(homeDir, ".hossted")
	if err != nil {
		return response{}, err
	}

	fileData, err := os.ReadFile(folderPath + "/" + emailID + ".json")
	if err != nil {
		return response{}, fmt.Errorf("User not registered, Please run hossted login to register")
	}

	var resp response
	err = json.Unmarshal(fileData, &resp)
	if err != nil {
		return response{}, err
	}

	return resp, nil
}

func useCases(resp response, emailID string) error {
	if resp.Success {
		if len(resp.OrgIDs) == 0 {
			fmt.Println("We have just sent the confirmation link to", emailID, ". Once you confirm it, you'll be able to continue the activation.")
		} else if len(resp.OrgIDs) == 1 {
			for orgID, email := range resp.OrgIDs[0] {
				prompt := promptui.Select{
					Label: fmt.Sprintf("Are you sure you want to register this cluster with org_name %s", email),
					Items: []string{"Yes", "No"},
				}
				_, value, err := prompt.Run()
				if err != nil {
					return err
				}
				if value == "Yes" {
					fmt.Println("Your orgID is ", orgID)
				} else {
					return nil
				}
			}
		} else if len(resp.OrgIDs) > 1 {
			// Handle cases where len(resp.OrgIDs) > 1
			fmt.Println("There are multiple organization IDs. Handling multiple org IDs logic here.")
		}
	} else {
		return fmt.Errorf("Cluster registration failed to hossted platform")
	}

	return nil
}

// promptK8sContext retrieves Kubernetes contexts from kubeconfig.
func promptK8sContext() (clusterName string, err error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}

	// Load kubeconfig file
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Get current context
	currentContext := config.Contexts
	var contexts []string
	for i := range currentContext {
		contexts = append(contexts, i)
	}

	// // Prompt user to select Kubernetes context
	promptK8s := promptui.Select{
		Label: "Select your Kubernetes context:",
		Items: contexts,
	}
	_, clusterName, err = promptK8s.Run()
	if err != nil {
		return "", err
	}

	// set current context as selected
	config.CurrentContext = clusterName
	err = clientcmd.WriteToFile(*config, kubeconfigPath)
	if err != nil {
		return "", err
	}

	return clusterName, nil
}

func deployOperator(clusterName, emailID string) error {
	operator := promptui.Select{
		Label: fmt.Sprintf("Do you wish to install the operator in %s", clusterName),
		Items: []string{"Yes", "No"},
	}
	_, value, err := operator.Run()
	if err != nil {
		return err
	}

	if value == "Yes" {

		monitoring := promptui.Select{
			Label: fmt.Sprintf("Do you wish to enable monitoring in operator"),
			Items: []string{"Yes", "No"},
		}
		_, monitoringEnable, err := monitoring.Run()
		if err != nil {
			return err
		}

		if monitoringEnable == "Yes" {
			fmt.Println("Enabled Monitoring ", monitoringEnable)
		}

		cve := promptui.Select{
			Label: fmt.Sprintf("Do you wish to enable cve scan in operator"),
			Items: []string{"Yes", "No"},
		}
		_, cveEnable, err := cve.Run()
		if err != nil {
			return err
		}

		if cveEnable == "Yes" {
			fmt.Println("Enabled CVE Scanning ", cveEnable)
			RepoAdd("aqua", "https://aquasecurity.github.io/helm-charts/")
			// Progress bar setup
			fmt.Println("Installing trivy-operator chart...")
			bar := progressbar.DefaultBytes(
				-1,
				"Installing",
			)

			// Simulate installation process with a time delay
			for i := 0; i < 100; i++ {
				time.Sleep(50 * time.Millisecond)
				bar.Add(1)
			}
			InstallChart("trivy-operator", "aqua", "trivy-operator", map[string]string{
				"set": "operator.scannerReportTTL=,operator.scanJobTimeout=30m",
			})

		}

		RepoAdd("hossted", "https://charts.hossted.com")

		// Progress bar setup
		fmt.Println("Installing hossted-operator chart...")
		bar := progressbar.DefaultBytes(
			-1,
			"Installing",
		)

		// Simulate installation process with a time delay
		for i := 0; i < 100; i++ {
			time.Sleep(50 * time.Millisecond)
			bar.Add(1)
		}

		args := map[string]string{
			"set": "env.EMAIL_ID=" + emailID,
		}
		InstallChart("hossted-operator", "hossted", "hossted-operator", args)

	}
	return nil
}

var settings *cli.EnvSettings

// RepoAdd adds repo with given name and url
func RepoAdd(name, url string) {
	settings = cli.New()

	repoFile := settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		log.Fatal(err)
	}

	if f.Has(name) {
		fmt.Printf("repository name (%s) already exists\n", name)
		return
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		log.Fatal(err)
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q has been added to your repositories\n", name)
}

// RepoUpdate updates charts for all helm repos
func RepoUpdate() {
	repoFile := settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		log.Fatal(errors.New("no repositories found. You must add one before updating"))
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, r)
	}

	fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	fmt.Printf("Update Complete. ⎈ Happy Helming!⎈\n")
}

func InstallChart(name, repo, chart string, args map[string]string) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Fatal(err)
	}
	client := action.NewInstall(actionConfig)

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}
	//name, chart, err := client.NameAndChart(args)
	client.ReleaseName = name
	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repo, chart), settings)
	if err != nil {
		log.Fatal(err)
	}

	debug("CHART PATH: %s\n", cp)

	p := getter.All(settings)
	valueOpts := &values.Options{}
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		log.Fatal(err)
	}

	// Add args
	if err := strvals.ParseInto(args["set"], vals); err != nil {
		log.Fatal(errors.Wrap(err, "failed parsing --set data"))
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		log.Fatal(err)
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		log.Fatal(err)
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
	}

	client.Namespace = "hossted-platform"
	client.CreateNamespace = true
	release, err := client.Run(chartRequested, vals)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Release Name: [%s] Status [%s]", release.Name, release.Info.Status)

}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}
