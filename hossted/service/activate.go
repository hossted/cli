package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	Email   string              `json:"email"`
	Message string              `json:"message"`
}

// ActivateK8s imports Kubernetes clusters.
func ActivateK8s() error {

	emailsID, err := getEmail()
	if err != nil {
		return err
	}

	// getResponse from reading file in .hossted/email@id.json
	resp, err := getResponse()
	if err != nil {
		return err
	}

	// handle usecases for orgID selection
	orgID, err := useCases(resp)
	if err != nil {
		return err
	}

	// prompt user for k8s context
	clusterName, err := promptK8sContext()
	if err != nil {
		return err
	}

	fmt.Println("Your cluster name is ", clusterName)

	err = deployOperator(clusterName, emailsID, orgID, resp.Token)
	if err != nil {
		return err
	}

	return nil
}

func getEmail() (string, error) {
	homeDir, err := os.UserHomeDir()
	folderPath := filepath.Join(homeDir, ".hossted")
	fileData, err := os.ReadFile(folderPath + "/" + "config.json")
	if err != nil {
		return "", err
	}

	// Parse the JSON data into Config struct
	var config response
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return "", err
	}

	return config.Email, nil

}

func getResponse() (response, error) {
	//read file
	homeDir, err := os.UserHomeDir()

	folderPath := filepath.Join(homeDir, ".hossted")
	if err != nil {
		return response{}, err
	}

	fileData, err := os.ReadFile(folderPath + "/" + "config.json")
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

func useCases(resp response) (orgID string, err error) {
	if resp.Success {
		if len(resp.OrgIDs) == 0 {
			for orgID := range resp.OrgIDs[0] {
				fmt.Println("We have just sent the confirmation link registered emailID", ". Once you confirm it, you'll be able to continue the activation.")
				return orgID, nil
			}
		} else if len(resp.OrgIDs) > 1 {
			fmt.Println("You have multiple organisations to choose from:")

			var items []string
			for i, info := range resp.OrgIDs {
				for _, orgName := range info {
					items = append(items, fmt.Sprintf("%d: %s", i+1, orgName))
				}
			}

			prompt := promptui.Select{
				Label: "Select Your Organisation",
				Items: items,
			}

			_, result, err := prompt.Run()
			if err != nil {
				fmt.Println("Prompt failed:", err)
				return "", err
			}

			userOrgName, err := removePrefix(result)
			if err != nil {
				return "", err
			}
			var selectedOrgID string

			for _, info := range resp.OrgIDs {
				for orgID, orgName := range info {
					if orgName == userOrgName {
						selectedOrgID = orgID
					}
				}
			}

			if selectedOrgID == "" {
				return "", fmt.Errorf("selected organization not found")
			}

			return selectedOrgID, nil

		}
	} else {
		return "", fmt.Errorf("Cluster registration failed to hossted platform")
	}

	return "", nil
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

func deployOperator(clusterName, emailID, orgID, JWT string) error {
	operator := promptui.Select{
		Label: fmt.Sprintf("Do you wish to install the operator in %s", clusterName),
		Items: []string{"Yes", "No"},
	}
	_, value, err := operator.Run()
	if err != nil {
		return err
	}

	if value == "Yes" {

		cveEnabled := "false"
		monitoringEnabled := "false"
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
			monitoringEnabled = "true"
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
			cveEnabled = "true"
		}

		RepoAdd("hossted", "https://charts.hossted.com")
		RepoUpdate()
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
			"set": "env.EMAIL_ID=" + emailID +
				//",env.HOSSTED_API_URL=https://api.dev.hossted.com/v1/instances" +
				",env.HOSSTED_ORG_ID=" + orgID +
				",secret.HOSSTED_AUTH_TOKEN=" + JWT +
				",cve.enable=" + cveEnabled +
				",monitoring.enable=" + monitoringEnabled +
				",env.MIMIR_PASSWORD=" + os.Getenv("MIMIR_PASSWORD") +
				",env.CONTEXT_NAME=" + clusterName,
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

	client.ReleaseName = name
	client.Namespace = "hossted-operator"
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

func removePrefix(text string) (string, error) {
	// Define a regular expression to match a number followed by a colon and a space
	regex := regexp.MustCompile(`^\d+:\s+`)

	match := regex.FindStringSubmatch(text)
	if match != nil {
		// Extract the captured prefix (number and colon)
		prefix := match[0]
		return strings.TrimPrefix(text, prefix), nil
	}

	return text, nil
}
