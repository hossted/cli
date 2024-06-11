package service

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// stand by mode
func InstallOperatorStandbymode() error {
	args := map[string]string{
		"set": "env.EMAIL_ID=" + "" +
			// ",env.HOSSTED_API_URL=https://api.dev.hossted.com/v1/instances" +
			// ",env.HOSSTED_ORG_ID=" + "" +
			// ",secret.HOSSTED_AUTH_TOKEN=" + "" +
			",cve.enable=" + "false" +
			",monitoring.enable=" + "false" +
			",logging.enable=" + "false" +
			",stop=" + "true",
		// ",env.LOKI_URL=" + os.Getenv("LOKI_URL") +
		// ",env.LOKI_USERNAME=" + os.Getenv("LOKI_USERNAME") +
		// ",env.LOKI_PASSWORD=" + os.Getenv("LOKI_PASSWORD") +
		// ",env.MIMIR_URL=" + os.Getenv("MIMIR_URL") +
		// ",env.MIMIR_USERNAME=" + os.Getenv("MIMIR_USERNAME") +
		// ",env.MIMIR_PASSWORD=" + os.Getenv("MIMIR_PASSWORD") +
		// ",env.CONTEXT_NAME=" + "",
	}

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%s Deploying in namespace %s\n", yellow("Hossted Platform Operator in Standby Modee:"), hosstedPlatformNamespace)

	bar := progressbar.DefaultBytes(
		-1,
		"Installing",
	)

	// Simulate installation process with a time delay
	for i := 0; i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
		bar.Add(1)
	}
	RepoAdd("hossted", "https://charts.hossted.com")
	RepoUpdate()
	InstallChart(hosstedOperatorReleaseName, "hossted", "hossted-operator", args)

	return nil
}

// ----------
// "set":      "env.EMAIL_ID=" + "" +
// 			",env.HOSSTED_API_URL=https://api.dev.hossted.com/v1/instances" +
// 			",env.HOSSTED_ORG_ID=" + "" +
// 			",secret.HOSSTED_AUTH_TOKEN=" + "" +
// 			",cve.enable=" + "false" +
// 			",monitoring.enable=" + "false" +
// 			",logging.enable=" + "false" +
// 			",stop=" + "true" +
// 			",env.LOKI_URL=" + "" +
// 			",env.LOKI_USERNAME=" + "" +
// 			",env.LOKI_PASSWORD=" + "" +
// 			",env.MIMIR_URL=" + "" +
// 			",env.MIMIR_USERNAME=" + "" +
// 			",env.MIMIR_PASSWORD=" + "" +
// 			",env.CONTEXT_NAME=" + "",
// -------
