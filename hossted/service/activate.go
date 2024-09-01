package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"

	"github.com/fatih/color"
	"github.com/hossted/cli/hossted/service/common"
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
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	hosstedPlatformNamespace   = "hossted-platform"
	hosstedOperatorReleaseName = "hossted-operator"
	trivyOperatorReleaseName   = "trivy-operator"
	grafanaAgentReleaseName    = "hossted-grafana-agent"
)

// ActivateK8s imports Kubernetes clusters.
func ActivateK8s(releaseName string, develMode bool) error {

	// emailsID, err := getEmail()
	// if err != nil {
	// 	return err
	// }

	// getResponse from reading file in .hossted/config.json
	// resp, err := getLoginResponse()
	// if err != nil {
	// 	return err
	// }
	// validate auth token

	// err = validateToken(resp)
	// if err != nil {
	// 	return err
	// }
	// handle usecases for orgID selection
	// orgID, err := useCases(resp)
	// if err != nil {
	// 	return err
	// }

	// prompt user for k8s context
	clusterName, err := promptK8sContext()
	if err != nil {
		return err
	}

	fmt.Println("Your cluster name is ", clusterName)

	isStandby, err := isStandbyMode(releaseName)
	if err != nil {
		return err
	}

	tr, err := common.GetTokenResp()
	if err != nil {
		return err
	}

	orgs, err := common.GetOrgs(tr.AccessToken)
	if err != nil {
		return err
	}

	orgID, err := common.OrgUseCases(orgs)
	if err != nil {
		return err
	}

	fmt.Println(orgID)

	if isStandby {
		fmt.Println("Standby mode detected")
		clientset := getKubeClient()
		fmt.Println("Updating deployment....")
		err := updateDeployment(clientset, hosstedPlatformNamespace, "hossted-operator"+"-controller-manager", "", clusterName, orgID, develMode)
		if err != nil {
			return err
		}

		// config, err := readConfig()
		// if err != nil {
		// 	return err
		// }

		fmt.Println("Updating secret....")
		err = updateSecret(clientset, hosstedPlatformNamespace, "hossted-operator"+"-secret", "AUTH_TOKEN", tr.AccessToken)
		if err != nil {
			return err
		}

		fmt.Println("Updated deployment and secret")

		cveEnabled, monitoringEnabled, loggingEnabled, ingressEnabled, err := askPromptsToInstall()
		if err != nil {
			return err
		}

		if cveEnabled == "true" {
			fmt.Println("Enabled CVE Scan:", cveEnabled)
			installCve()
		}

		if monitoringEnabled == "true" || loggingEnabled == "true" || cveEnabled == "true" || ingressEnabled == "true" {
			fmt.Println("Patching 'hossted-operator-cr' CR")
			err = patchCR(monitoringEnabled, loggingEnabled, cveEnabled, ingressEnabled, releaseName)
			if err != nil {
				return err
			}
		}

		err = patchStopCR(releaseName)
		if err != nil {
			return err
		}

		fmt.Println("Patch'hossted-operator-cr' CR completed")
		return nil
	}

	err = deployOperator(clusterName, "", orgID, tr.AccessToken, develMode)
	if err != nil {
		return err
	}

	return nil
}

func isStandbyMode(releaseName string) (bool, error) {

	isStandby := false
	cr := getDynClient()
	hp, err := cr.Resource(hpGVK).Get(context.TODO(), "hossted-operator-cr", metav1.GetOptions{})
	if err != nil {
		return isStandby, err
	}

	// fmt.Println(hp)
	// cve, _, err := unstructured.NestedMap(hp.Object, "spec", "cve")
	// if err != nil {
	// 	return isStandby, err
	// }
	// cveEnabled, _, err := unstructured.NestedBool(cve, "enable")
	// if err != nil {
	// 	return isStandby, err
	// }

	// logging, _, err := unstructured.NestedMap(hp.Object, "spec", "logging")
	// if err != nil {
	// 	return isStandby, err
	// }
	// loggingEnabled, _, err := unstructured.NestedBool(logging, "enable")
	// if err != nil {
	// 	return isStandby, err
	// }

	// monitoring, _, err := unstructured.NestedMap(hp.Object, "spec", "monitoring")
	// if err != nil {
	// 	return isStandby, err
	// }
	// monitoringEnabled, _, err := unstructured.NestedBool(monitoring, "enable")
	// if err != nil {
	// 	return isStandby, err
	// }
	// ingress, _, err := unstructured.NestedMap(hp.Object, "spec", "ingress")
	// if err != nil {
	// 	return isStandby, err
	// }

	// <<<<<<< edge-standby
	// =======
	// 	monitoring, _, err := unstructured.NestedMap(hp.Object, "spec", "monitoring")
	// 	if err != nil {
	// 		return isStandby, err
	// 	}
	// 	monitoringEnabled, _, err := unstructured.NestedBool(monitoring, "enable")
	// 	if err != nil {
	// 		return isStandby, err
	// 	}
	// 	ingress, _, err := unstructured.NestedMap(hp.Object, "spec", "ingress")
	// 	if err != nil {
	// 		return isStandby, err
	// 	}
	// 	ingressEnabled, _, err := unstructured.NestedBool(ingress, "enable")
	// 	if err != nil {
	// 		return isStandby, err
	// 	}
	// >>>>>>> dev
	stop, _, err := unstructured.NestedBool(hp.Object, "spec", "stop")
	if err != nil {
		return isStandby, err
	}

	if stop {
		isStandby = true
	}

	return isStandby, nil
}

// func getEmail() (string, error) {
// 	config, err := readConfig()
// 	if err != nil {
// 		return "", err
// 	}
// 	return config.Email, nil
// }

// func readConfig() (response, error) {
// 	var config response
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		return config, err
// 	}
// 	folderPath := filepath.Join(homeDir, ".hossted")
// 	fileData, err := os.ReadFile(folderPath + "/" + "config.json")
// 	if err != nil {
// 		return config, err
// 	}

// 	// Parse the JSON data into Config struct
// 	err = json.Unmarshal(fileData, &config)
// 	if err != nil {
// 		return config, err
// 	}
// 	return config, nil
// }

// func getLoginResponse() (response, error) {
// 	//read file
// 	homeDir, err := os.UserHomeDir()

// 	folderPath := filepath.Join(homeDir, ".hossted")
// 	if err != nil {
// 		return response{}, err
// 	}

// 	fileData, err := os.ReadFile(folderPath + "/" + "config.json")
// 	if err != nil {
// 		return response{}, fmt.Errorf("User not registered, Please run hossted login to register")
// 	}

// 	var resp response
// 	err = json.Unmarshal(fileData, &resp)
// 	if err != nil {
// 		return response{}, err
// 	}

// 	return resp, nil
// }

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

func getKubeClient() *kubernetes.Clientset {
	var kubeconfig string
	path, ok := os.LookupEnv("KUBECONFIG")
	if ok {
		kubeconfig = path
	} else {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func updateDeployment(
	clientset *kubernetes.Clientset,
	namespace, deploymentName, emailID, contextName, hosstedOrgID string,
	develMode bool) error {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	hosstedApiUrl := common.HOSSTED_API_URL
	mimirUrl := common.MIMIR_URL
	lokiUrl := common.LOKI_URL

	if develMode {

		if devUrl := common.HOSSTED_DEV_API_URL; devUrl != "" {
			hosstedApiUrl = devUrl
		}
		if devUrl := common.MIMIR_DEV_URL; devUrl != "" {
			mimirUrl = devUrl
		}
		if devUrl := common.LOKI_DEV_URL; devUrl != "" {
			lokiUrl = devUrl
		}
	}

	// Update environment variables
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == "manager" {
			for j, env := range container.Env {
				if env.Name == "EMAIL_ID" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = emailID
				} else if env.Name == "CONTEXT_NAME" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = contextName
				} else if env.Name == "HOSSTED_ORG_ID" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = hosstedOrgID
				} else if env.Name == "LOKI_PASSWORD" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = common.LOKI_PASSWORD
				} else if env.Name == "LOKI_URL" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = lokiUrl
				} else if env.Name == "LOKI_USERNAME" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = common.LOKI_USERNAME
				} else if env.Name == "MIMIR_PASSWORD" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = common.MIMIR_PASSWORD
				} else if env.Name == "MIMIR_URL" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = mimirUrl
				} else if env.Name == "MIMIR_USERNAME" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = common.MIMIR_USERNAME
				} else if env.Name == "HOSSTED_API_URL" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = hosstedApiUrl

				}
			}
		}
	}

	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func updateSecret(clientset *kubernetes.Clientset, namespace, secretName, secretKey, secretValue string) error {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Update the secret
	secret.Data[secretKey] = []byte(secretValue)

	_, err = clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func patchCR(monitoringEnabled, loggingEnabled, cveEnabled, ingressEnabled, releaseName string) error {
	cr := getDynClient()
	hp, err := cr.Resource(hpGVK).Get(context.TODO(), releaseName+"-cr", metav1.GetOptions{})
	if err != nil {
		return err
	}

	patch := map[string]interface{}{
		"spec": map[string]interface{}{
			"cve": map[string]interface{}{
				"enable": cveEnabled == "true",
			},
			"logging": map[string]interface{}{
				"enable": loggingEnabled == "true",
			},
			"monitoring": map[string]interface{}{
				"enable": monitoringEnabled == "true",
			},
			"ingress": map[string]interface{}{
				"enable": ingressEnabled == "true",
			},
		},
	}

	patchData, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	// Apply the patch
	_, err = cr.Resource(hpGVK).Namespace(hp.GetNamespace()).Patch(context.TODO(), hp.GetName(), types.MergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func patchStopCR(releaseName string) error {
	cr := getDynClient()
	hp, err := cr.Resource(hpGVK).Get(context.TODO(), "hossted-operator"+"-cr", metav1.GetOptions{})
	if err != nil {
		return err
	}

	patch := map[string]interface{}{
		"spec": map[string]bool{
			"stop": false,
		},
	}

	patchData, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	// Apply the patch
	_, err = cr.Resource(hpGVK).Namespace(hp.GetNamespace()).Patch(context.TODO(), hp.GetName(), types.MergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func askPromptsToInstall() (string, string, string, string, error) {
	green := color.New(color.FgGreen).SprintFunc()
	cveEnabled := "false"
	monitoringEnabled := "false"
	loggingEnabled := "false"
	ingressEnabled := "false"
	//------------------------------ Monitoring ----------------------------------
	monitoring := promptui.Select{
		Label: "Do you wish to enable monitoring in hossted platform",
		Items: []string{"Yes", "No"},
	}
	_, monitoringEnable, err := monitoring.Run()
	if err != nil {
		return cveEnabled, monitoringEnabled, loggingEnabled, ingressEnabled, err
	}

	if monitoringEnable == "Yes" {
		fmt.Println("Enabled Monitoring :", green(monitoringEnable))
		monitoringEnabled = "true"
	}

	//------------------------------ Logging ----------------------------------
	logging := promptui.Select{
		Label: "Do you wish to enable logging in hossted-platform",
		Items: []string{"Yes", "No"},
	}
	_, loggingEnable, err := logging.Run()
	if err != nil {
		return cveEnabled, monitoringEnabled, loggingEnabled, ingressEnabled, err
	}

	if loggingEnable == "Yes" {
		fmt.Println("Enabled Logging:", green(loggingEnable))
		loggingEnabled = "true"
	}

	//------------------------------ CVE Scan ----------------------------------
	cve := promptui.Select{
		Label: "Do you wish to enable cve scan in hossted platform",
		Items: []string{"Yes", "No"},
	}
	_, cveEnable, err := cve.Run()
	if err != nil {
		return cveEnabled, monitoringEnabled, loggingEnabled, ingressEnabled, err
	}
	if cveEnable == "Yes" {
		fmt.Println("Enabled CVE:", green(cveEnable))
		cveEnabled = "true"
	}

	//------------------------------ Ingress ----------------------------------
	ingress := promptui.Select{
		Label: "Do you wish to enable ingress in hossted-platform",
		Items: []string{"Yes", "No"},
	}
	_, ingressEnable, err := ingress.Run()
	if err != nil {
		return cveEnabled, monitoringEnabled, loggingEnabled, ingressEnabled, err
	}

	if ingressEnable == "Yes" {
		ingressEnabled = "true"
		fmt.Println("Enabled Ingress:", green(ingressEnabled))
	}
	return cveEnabled, monitoringEnabled, loggingEnabled, ingressEnabled, nil
}

func installCve() {
	yellow := color.New(color.FgYellow).SprintFunc()
	RepoAdd("aqua", "https://aquasecurity.github.io/helm-charts/")

	// Progress bar setup
	fmt.Printf("%s Deploying in namespace %s\n", yellow("Hossted Platform CVE:"), hosstedPlatformNamespace)

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
		"set": "trivy.severity=HIGH\\,CRITICAL,operator.scannerReportTTL=,operator.scanJobTimeout=30m,trivy.command=filesystem,trivyOperator.scanJobPodTemplateContainerSecurityContext.runAsUser=0,operator.scanJobsConcurrentLimit=10",
	})
}

func deployOperator(clusterName, emailID, orgID, JWT string, develMode bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	operator := promptui.Select{
		Label: fmt.Sprintf("Do you wish to install the hossted platform in %s", yellow(clusterName)),
		Items: []string{"Yes", "No"},
	}
	_, value, err := operator.Run()
	if err != nil {
		return err
	}

	if value == "Yes" {
		cveEnabled, monitoringEnabled, loggingEnabled, ingressEnabled, err := askPromptsToInstall()
		if err != nil {
			return err
		}

		if cveEnabled == "true" {
			fmt.Println("Enabled CVE Scan:", green(cveEnabled))
			installCve()
		}

		//------------------------------ Helm Repo Add  ----------------------------------

		RepoAdd("hossted", "https://charts.hossted.com")
		RepoAdd("ingress-nginx", "https://kubernetes.github.io/ingress-nginx")
		RepoUpdate()

		fmt.Println(loggingEnabled)

		//------------------------------ Helm Install Chart ----------------------------------

		hosstedApiUrl := common.HOSSTED_API_URL
		mimirUrl := common.MIMIR_URL
		lokiUrl := common.LOKI_URL

		if develMode {

			if devUrl := common.HOSSTED_DEV_API_URL; devUrl != "" {
				hosstedApiUrl = devUrl
			}
			if devUrl := common.MIMIR_DEV_URL; devUrl != "" {
				mimirUrl = devUrl
			}
			if devUrl := common.LOKI_DEV_URL; devUrl != "" {
				lokiUrl = devUrl
			}
		}

		args := map[string]string{
			"set": "env.EMAIL_ID=" + emailID +
				",env.HOSSTED_ORG_ID=" + orgID +
				",secret.HOSSTED_AUTH_TOKEN=" + JWT +
				",cve.enable=" + cveEnabled +
				",monitoring.enable=" + monitoringEnabled +
				",logging.enable=" + loggingEnabled +
				",ingress.enable=" + ingressEnabled +
				",env.LOKI_URL=" + lokiUrl +
				",env.LOKI_USERNAME=" + common.LOKI_USERNAME +
				",env.LOKI_PASSWORD=" + common.LOKI_PASSWORD +
				",env.MIMIR_URL=" + mimirUrl +
				",env.MIMIR_USERNAME=" + common.MIMIR_USERNAME +
				",env.MIMIR_PASSWORD=" + common.MIMIR_PASSWORD +
				",env.HOSSTED_API_URL=" + hosstedApiUrl +
				",env.CONTEXT_NAME=" + clusterName,
		}

		fmt.Printf("%s Deploying in namespace %s\n", yellow("Hossted Platform Operator:"), hosstedPlatformNamespace)

		bar := progressbar.DefaultBytes(
			-1,
			"Installing",
		)

		// Simulate installation process with a time delay
		for i := 0; i < 100; i++ {
			time.Sleep(50 * time.Millisecond)
			bar.Add(1)
		}

		InstallChart(hosstedOperatorReleaseName, "hossted", "hossted-operator", args)

		//------------------------------ Add Events ----------------------------------
		err = addEvents(JWT)
		if err != nil {
			return err
		}

		fmt.Println(green("Success: "), "Hossted Platfrom components deployed")
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
		//fmt.Printf("repository name (%s) already exists\n", name)
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

	//fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				//	fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				//	fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	//fmt.Printf("Updated Helm Repos Complete. ⎈ Happy Helming!⎈\n")
}

func InstallChart(name, repo, chart string, args map[string]string) {
	actionConfig := new(action.Configuration)
	clientGetter := genericclioptions.NewConfigFlags(false)
	namespace := hosstedPlatformNamespace
	clientGetter.Namespace = &namespace

	if err := actionConfig.Init(clientGetter, hosstedPlatformNamespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
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
	client.Namespace = hosstedPlatformNamespace
	client.CreateNamespace = true
	release, err := client.Run(chartRequested, vals)
	if err != nil {
		log.Fatal(err)
	}
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("Release Name: [%s] Status [%s]\n", green(release.Name), green(release.Info.Status))

}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func debug(format string, v ...interface{}) {
	//format = fmt.Sprintf("[debug] %s\n", format)
	//log.Output(2, fmt.Sprintf(format, v...))
}

func addEvents(token string) error {

	if err := eventOperator(token); err != nil {
		return err
	}
	if err := eventCVE(token); err != nil {
		return err
	}
	if err := eventMonitoring(token); err != nil {
		return err
	}
	return nil
}

func eventMonitoring(token string) error {
	retries := 60
	for i := 0; i < retries; i++ {
		err := checkMonitoringStatus()
		if err == nil {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s Hossted Platform Monitoring started successfully\n", green("Success:"))
			err := SendEvent("success", "Hossted Platform Monitoring started successfully", token)
			if err != nil {
				return err
			}
			return nil
		}

		// If not successful, wait for a short duration before retrying
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Println(yellow("Info:"), "Waiting Hosted Platform Monitoring Agents to get into running state.")
		time.Sleep(3 * time.Second)
	}

	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s Timeout reached. Hossted Platform Monitoring installation failed.\n", red("Error:"))
	return fmt.Errorf("Hossted Platform Monitoring installation failed after %d retries", retries)
}

func checkMonitoringStatus() error {
	releases, err := listReleases()
	if err != nil {
		return err
	}

	for _, release := range releases {
		if release.Name == grafanaAgentReleaseName {
			return nil
		}
	}

	return fmt.Errorf("Grafana Agent release not found")
}

func eventCVE(token string) error {
	releases, err := listReleases()
	if err != nil {
		return err
	}

	timeout := time.After(120 * time.Second)
	for {
		select {
		case <-timeout:
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s Timeout reached. Hossted Platform CVE installation failed.\n", red("Error:"))
			return nil
		default:
			for _, release := range releases {
				if release.Name == trivyOperatorReleaseName {
					green := color.New(color.FgGreen).SprintFunc()
					fmt.Printf("%s Hossted Platform CVE Scan started successfully\n", green("Success:"))
					err := SendEvent("success", "Hossted Platform CVE Scan started successfully", token)
					if err != nil {
						return err
					}
					return nil
				}
			}
			// Sleep for a short duration before checking again
			time.Sleep(1 * time.Second)
		}
	}
}

func eventOperator(token string) error {
	releases, err := listReleases()
	if err != nil {
		return err
	}

	timeout := time.After(120 * time.Second)
	for {
		select {
		case <-timeout:
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s Timeout reached. Hossted Platform operator installation failed.\n", red("Error:"))
			return nil
		default:
			for _, release := range releases {
				if release.Name == hosstedOperatorReleaseName {
					green := color.New(color.FgGreen).SprintFunc()
					fmt.Printf("%s Hossted Platform operator installed successfully\n", green("Success:"))
					err := SendEvent("success", "Hossted Platform operator installed successfully", token)
					if err != nil {
						return err
					}
					return nil
				}
			}
			// Sleep for a short duration before checking again
			time.Sleep(1 * time.Second)
		}
	}
}

func listReleases() ([]*release.Release, error) {
	actionConfig := new(action.Configuration)
	clientGetter := genericclioptions.NewConfigFlags(false)
	namespace := hosstedPlatformNamespace
	clientGetter.Namespace = &namespace

	if err := actionConfig.Init(clientGetter, hosstedPlatformNamespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Fatal(err)
	}

	client := action.NewList(actionConfig)
	client.Deployed = true

	return client.Run()
}

func getDynClient() *dynamic.DynamicClient {

	var conf *rest.Config
	var err error

	// for running locally

	var kubeconfig string
	path, ok := os.LookupEnv("KUBECONFIG")
	if ok {
		kubeconfig = path
	} else {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	conf, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	dynClient, err := dynamic.NewForConfig(conf)
	if err != nil {
		panic(err.Error())
	}

	return dynClient
}

var hpGVK = schema.GroupVersionResource{
	Group:    "hossted.com",
	Version:  "v1",
	Resource: "hosstedprojects",
}

func SendEvent(eventType, message, token string) error {
	url := common.HOSSTED_API_URL + "/statuses"

	type event struct {
		WareType string `json:"ware_type"`
		Type     string `json:"type"`
		UUID     string `json:"uuid"`
		Message  string `json:"message"`
	}

	clusterUUID, err := getClusterUUID()
	if err != nil {
		return err
	}

	newEvent := event{
		WareType: "k8s",
		Type:     eventType,
		UUID:     clusterUUID,
		Message:  message,
	}

	eventByte, err := json.Marshal(newEvent)
	if err != nil {
		return err
	}
	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(eventByte)))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Add Authorization header with Basic authentication
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", []byte(token)))
	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Error sending event, errcode: %d", resp.StatusCode)
	}

	return nil
}

func getClusterUUID() (string, error) {
	var clusterUUID string
	var err error
	yellow := color.New(color.FgYellow).SprintFunc()

	// Retry for 120 seconds
	for i := 0; i < 120; i++ {
		clusterUUID, err = getClusterUUIDFromK8s()
		if err == nil {
			return clusterUUID, nil
		}
		fmt.Println(yellow("Info:"), "Waiting for Hossted Operator to get into running state")
		time.Sleep(4 * time.Second) // Wait for 1 second before retrying
	}

	return "", fmt.Errorf("Failed to get ClusterUUID after 120 seconds: %v", err)
}

func getClusterUUIDFromK8s() (string, error) {
	cs := getDynClient()
	hp, err := cs.Resource(hpGVK).Get(context.TODO(), "hossted-operator-cr", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	clusterUUID, _, err := unstructured.NestedString(hp.Object, "status", "clusterUUID")
	if err != nil || clusterUUID == "" {
		return "", fmt.Errorf("ClusterUUID is nil, func errored")
	}
	return clusterUUID, nil
}

// func validateToken(res response) error {

// 	type validationResp struct {
// 		Success bool   `json:"success"`
// 		Message string `json:"message"`
// 	}

// 	authToken := common.HOSSTED_AUTH_TOKEN
// 	url := common.HOSSTED_API_URL + "/cli/bearer"

// 	payloadStr := fmt.Sprintf(`{"email": "%s", "token": "%s"}`, res.Email, res.Token)
// 	// Create HTTP request
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payloadStr)))
// 	if err != nil {
// 		return err
// 	}

// 	// Set headers
// 	req.Header.Set("Content-Type", "application/json")

// 	// Add Authorization header with Basic authentication
// 	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", []byte(authToken)))

// 	// Perform the request
// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	// check for non 200 status
// 	if resp.StatusCode != 200 {
// 		return fmt.Errorf("Token Validation Failed, Error %s", string(body))
// 	}

// 	var tresp validationResp
// 	err = json.Unmarshal(body, &tresp)
// 	if err != nil {
// 		return err
// 	}
// 	if !tresp.Success {
// 		return fmt.Errorf("Auth token is invalid, Please login again")
// 	}

// 	return nil

// }
