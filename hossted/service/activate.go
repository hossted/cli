package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	grafanaAgentReleaseName    = "hossted-grafana-alloy"
	releaseName                = "hossted-operator-cr"
)

var AUTH_TOKEN = common.HOSSTED_AUTH_TOKEN

const (
	activation_started      = "_HSTD_ACTIVATION_STARTED"
	init_operator           = "_HSTD_INIT_OPERATOR"
	init_operator_success   = "_HSTD_INIT_OPERATOR_SUCCESS"
	init_operator_fail      = "_HSTD_INIT_OPERATOR_FAIL"
	init_cve                = "_HSTD_INIT_CVE"
	init_cve_success        = "_HSTD_INIT_CVE_SUCCESS"
	init_cve_fail           = "_HSTD_INIT_CVE_FAIL"
	init_monitoring         = "_HSTD_INIT_MONITORING"
	init_monitoring_success = "_HSTD_INIT_MONITORING_SUCCESS"
	deployed_operator       = "_HSTD_DEPLOYED_OPERATOR"
	deployed_cve            = "_HSTD_DEPLOYED_CVE"
	deployed_monitoring     = "_HSTD_DEPLOYED_MONITORING"
	activation_complete     = "_HSTD_ACTIVATION_COMPLETE"
)

// ActivateK8s imports Kubernetes clusters.
func ActivateK8s(develMode, verbose bool) error {
	// Prompt user for k8s context
	clusterName, err := promptK8sContext()
	if err != nil {
		return err
	}

	isStandby, _ := isStandbyMode(releaseName)

	tr, err := common.GetTokenResp()
	if err != nil {
		return err
	}

	orgs, userID, err := common.GetOrgs(tr.AccessToken)
	if err != nil {
		return err
	}

	orgID, err := common.OrgUseCases(orgs)
	if err != nil {
		return err
	}

	err = SendEvent("info", activation_started, common.HOSSTED_AUTH_TOKEN, orgID, "", userID, verbose)
	if err != nil {
		fmt.Printf("\033[31mError sending event: %v\033[0m\n", err) // Red for error message
	}

	if isStandby {
		fmt.Println("\033[33mStandby mode detected\033[0m") // Yellow for standby detection
		clientset := getKubeClient()
		fmt.Println("\033[33mUpdating deployment....\033[0m") // Yellow for update status
		err := updateDeployment(clientset, hosstedPlatformNamespace, "hossted-operator"+"-controller-manager", "", clusterName, orgID, userID, develMode)
		if err != nil {
			return err
		}

		fmt.Println("\033[33mUpdating secret....\033[0m") // Yellow for secret update
		err = updateSecret(clientset, hosstedPlatformNamespace, "hossted-operator"+"-secret", "AUTH_TOKEN", tr.AccessToken)
		if err != nil {
			return err
		}

		fmt.Println("\033[32mUpdated deployment and secret\033[0m") // Green for successful update

		// Check if deployment is fully available
		deploymentName := "hossted-operator-controller-manager"
		namespace := hosstedPlatformNamespace
		timeout := 5 * time.Minute        // Set a timeout period
		checkInterval := 10 * time.Second // Poll every 10 seconds

		err = waitForDeploymentAvailability(clientset, namespace, deploymentName, timeout, checkInterval)
		if err != nil {
			fmt.Printf("\033[31mDeployment not available after 5 min: %v\033[0m\n", err) // Red for error message
			return err
		}

		fmt.Println("\033[32mDeployment is fully available. Proceeding...\033[0m") // Green for success

		cveEnabled, monitoringEnabled, loggingEnabled, ingressEnabled, err := askPromptsToInstall()
		if err != nil {
			return err
		}

		if cveEnabled == "true" {
			fmt.Println("\033[32mEnabled CVE Scan:\033[0m", cveEnabled) // Green for CVE scan
			installCve()
		}

		if monitoringEnabled == "true" || loggingEnabled == "true" || cveEnabled == "true" || ingressEnabled == "true" {
			fmt.Println("\033[33mPatching 'hossted-operator-cr' CR\033[0m") // Yellow for patching
			err = patchCR(monitoringEnabled, loggingEnabled, cveEnabled, ingressEnabled, releaseName)
			if err != nil {
				return err
			}
		}

		err = patchStopCR(releaseName)
		if err != nil {
			return err
		}

		fmt.Println("\033[32mPatch 'hossted-operator-cr' CR completed\033[0m") // Green for completion
		return nil
	}

	err = deployOperator(clusterName, "", orgID, tr.AccessToken, userID, develMode, verbose)
	if err != nil {
		_ = SendEvent("info", init_operator_fail, AUTH_TOKEN, orgID, "", userID, verbose)
		return err
	}

	err = SendEvent("info", init_operator_success, AUTH_TOKEN, orgID, "", userID, verbose)
	if err != nil {
		fmt.Printf("\033[31mError sending event: %v\033[0m\n", err) // Red for error message
	}

	err = SendEvent("info", init_monitoring_success, AUTH_TOKEN, orgID, "", userID, verbose)
	if err != nil {
		fmt.Printf("\033[31mError sending event: %v\033[0m\n", err) // Red for error message
	}

	err = SendEvent("info", init_cve_success, AUTH_TOKEN, orgID, "", userID, verbose)
	if err != nil {
		fmt.Println(err)
	}

	err = SendEvent("info", activation_complete, AUTH_TOKEN, orgID, "", userID, verbose)
	if err != nil {
		fmt.Printf("\033[31mError sending event: %v\033[0m\n", err) // Red for error message
	}

	return nil
}

func isStandbyMode(releaseName string) (bool, error) {

	isStandby := false
	cr := getDynClient()
	hp, err := cr.Resource(hpGVK).Get(context.TODO(), releaseName, metav1.GetOptions{})
	if err != nil {
		return isStandby, err
	}

	stop, _, err := unstructured.NestedBool(hp.Object, "spec", "stop")
	if err != nil {
		return isStandby, err
	}

	if stop {
		isStandby = true
	}

	return isStandby, nil
}

// waitForDeploymentAvailability checks if the deployment is fully available
func waitForDeploymentAvailability(clientset *kubernetes.Clientset, namespace, deploymentName string, timeout, interval time.Duration) error {
	startTime := time.Now()

	for {
		fmt.Println("\033[33mWaiting for deployments to be active....\033[0m") // Yellow message
		time.Sleep(15 * time.Second)
		// Get the latest version of the deployment
		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("\033[31mfailed to get deployment: %v\033[0m", err) // Red error message
		}

		// Check if the deployment is fully available
		if deployment.Status.AvailableReplicas == *deployment.Spec.Replicas {
			fmt.Println("\033[32mhossted operator controller manager deployment is fully available\033[0m") // Green success message
			return nil                                                                                      // Deployment is fully available
		}

		// Check if the timeout has been reached
		if time.Since(startTime) > timeout {
			return fmt.Errorf("\033[31mtimeout waiting for deployment to become available\033[0m") // Red error message
		}

		// Wait for the next check
		time.Sleep(interval)
	}
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
		fmt.Printf("\033[31mError loading kubeconfig: %v\033[0m\n", err)
		os.Exit(1)
	}

	// Get current contexts
	currentContext := config.Contexts
	var contexts []string
	for i := range currentContext {
		contexts = append(contexts, i)
	}

	// Prompt user to select Kubernetes context
	promptK8s := promptui.Select{
		Label: "\033[32mSelect your Kubernetes context:\033[0m",
		Items: contexts,
		Templates: &promptui.SelectTemplates{
			Active:   "\033[33m▸ {{ . }}\033[0m", // Yellow arrow and context name for active selection
			Inactive: "{{ . }}",                  // Default color for inactive items
			Selected: "\033[32mKubernetes context '{{ . }}' selected successfully.\033[0m",
		},
	}

	_, clusterName, err = promptK8s.Run()
	if err != nil {
		return "", err
	}

	// Set current context as selected
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
	namespace, deploymentName, emailID, contextName, hosstedOrgID, userID string,
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
				} else if env.Name == "HOSSTED_TOKEN" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = common.HOSSTED_AUTH_TOKEN
				} else if env.Name == "HOSSTED_USER_ID" {
					deployment.Spec.Template.Spec.Containers[i].Env[j].Value = userID
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
	hp, err := cr.Resource(hpGVK).Get(context.TODO(), "hossted-operator-cr", metav1.GetOptions{})
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
	hp, err := cr.Resource(hpGVK).Get(context.TODO(), releaseName, metav1.GetOptions{})
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
		fmt.Println("Enabled Monitoring:", green(monitoringEnable))
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

func deployOperator(clusterName, emailID, orgID, JWT, userID string, develMode, verbose bool) error {
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
			err = SendEvent("info", init_cve, AUTH_TOKEN, orgID, "", userID, verbose)
			if err != nil {
				fmt.Println(err)
			}
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
				",env.CONTEXT_NAME=" + clusterName +
				",env.HOSSTED_TOKEN=" + common.HOSSTED_AUTH_TOKEN +
				",env.HOSSTED_USER_ID=" + userID +
				",serviceAccount.create=" + "true",
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

		err = SendEvent("info", init_operator, AUTH_TOKEN, orgID, "", userID, verbose)
		if err != nil {
			fmt.Printf("\033[31mError sending event: %v\033[0m\n", err) // Red for error message
		}

		err = SendEvent("info", init_monitoring, AUTH_TOKEN, orgID, "", userID, verbose)
		if err != nil {
			fmt.Printf("\033[31mError sending event: %v\033[0m\n", err) // Red for error message
		}

		InstallChart(hosstedOperatorReleaseName, "hossted", "hossted-operator", args)

		//------------------------------ Add Events ----------------------------------
		// clusterUUID, err := getClusterUUIDPolling()
		// if err != nil {
		// 	return err
		// }

		// err = addEvents(AUTH_TOKEN, orgID, clusterUUID, userID, verbose)
		// if err != nil {
		// 	return err
		// }

		fmt.Println(green("Success: "), "Hossted Platform components deployed")
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

	b, err := os.ReadFile(repoFile)
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

func addEvents(token, orgID, clusterUUID, userID string, verbose bool) error {

	if err := eventOperator(token, orgID, clusterUUID, userID, verbose); err != nil {
		return err
	}
	if err := eventCVE(token, orgID, clusterUUID, userID, verbose); err != nil {
		return err
	}
	if err := eventMonitoring(token, orgID, clusterUUID, userID, verbose); err != nil {
		return err
	}
	return nil
}

func eventMonitoring(token, orgID, clusterUUID, userID string, verbose bool) error {
	retries := 60
	for i := 0; i < retries; i++ {
		err := checkMonitoringStatus()
		if err == nil {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s Hossted Platform Monitoring started successfully\n", green("Success:"))
			err := SendEvent("info", deployed_monitoring, token, orgID, clusterUUID, userID, verbose)
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
	fmt.Printf("%s timeout reached. Hossted Platform Monitoring installation failed.\n", red("Error:"))
	return fmt.Errorf("hossted Platform Monitoring installation failed after %d retries", retries)
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

	return fmt.Errorf("grafana Agent release not found")
}

func eventCVE(token, orgID, clusterUUID, userID string, verbose bool) error {
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
					err := SendEvent("info", deployed_cve, token, orgID, clusterUUID, userID, verbose)
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

func eventOperator(token, orgID, clusterUUID, userID string, verbose bool) error {
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
					err := SendEvent("info", deployed_operator, token, orgID, clusterUUID, userID, verbose)
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

func SendEvent(eventType, message, token, orgID, clusterUUID, userID string, verbose bool) error {
	url := common.HOSSTED_API_URL + "/statuses"

	type event struct {
		WareType string `json:"ware_type"`
		Type     string `json:"type"`
		UUID     string `json:"uuid,omitempty"`
		UserID   string `json:"user_id"`
		OrgID    string `json:"org_id"`
		Message  string `json:"message"`
	}

	newEvent := event{
		WareType: "k8s",
		Type:     eventType,
		UUID:     clusterUUID,
		UserID:   userID,
		OrgID:    orgID,
		Message:  message,
	}
	eventByte, err := json.MarshalIndent(newEvent, "", "  ")
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
		return fmt.Errorf("error sending event, errcode: %d", resp.StatusCode)
	}

	if verbose {
		fmt.Printf("\033[32mSuccess: Event '%s' sent successfully! Message: %s\033[0m\n", eventType, message)
	}

	return nil
}

func getClusterUUIDPolling() (string, error) {
	var clusterUUID string
	var err error
	yellow := color.New(color.FgYellow).SprintFunc()

	//Retry for 120 seconds
	for i := 0; i < 120; i++ {
		clusterUUID, err = getClusterUUIDFromK8s()
		if err == nil {
			return clusterUUID, nil
		}
		fmt.Println(yellow("Info:"), "Waiting for Hossted Operator to get into running state")
		time.Sleep(4 * time.Second) // Wait for 1 second before retrying
	}

	return "", fmt.Errorf("failed to get ClusterUUID after 120 seconds: %v", err)
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
