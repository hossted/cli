package cmd

import (
	"fmt"
	"os"

	"github.com/hossted/cli/hossted"
	"github.com/hossted/cli/hossted/service"
	"github.com/hossted/cli/hossted/service/compose"

	"github.com/spf13/cobra"
)

var (
	activate_type, composeFilePath string
	develMode, verbose             bool
)

// registerCmd represents the register command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "[a] Activate your application with the hossted platform",
	Long: `
Hossted activate connects your instance to the hossted platform and sends instance health information so it can be managed in the dashboard.
	`,
	Aliases: []string{"a"},
	Example: `
hossted activate
	`,
	Run: func(cmd *cobra.Command, args []string) {

		validTypes := map[string]bool{
			"k8s":     true,
			"compose": true,
			"standby": true,
		}

		if !validTypes[activate_type] {
			fmt.Printf("\033[31mInvalid type: %s. Valid types are: k8s, compose, standby\033[0m\n", activate_type)
			os.Exit(1)
		}

		// Write activate_type to config file
		config, _ := hossted.GetConfig() // Ignore error
		config.ActivateType = activate_type

		err := hossted.WriteConfigWrapper(config)
		if err != nil {
			fmt.Printf("\033[31mCannot write ActivateType to config file. Please check: %v\033[0m\n", err)
			return
		}

		if activate_type == "k8s" {
			err = service.VerifyAuth(develMode)
			if err != nil {
				fmt.Printf("\033[31mAuth verification failed: %v\033[0m\n", err)
				return
			}

			err = service.ActivateK8s(develMode, verbose)
			if err != nil {
				fmt.Printf("\033[31mKubernetes activation failed: %v\033[0m\n", err)
			} else {
				fmt.Printf("\033[32mKubernetes activated successfully!\033[0m\n")
			}
			return
		} else if activate_type == "compose" {
			err = service.VerifyAuth(develMode)
			if err != nil {
				fmt.Printf("\033[31mAuth verification failed: %v\033[0m\n", err)
				return
			}

			if composeFilePath == "" {
				dir, err := os.Getwd()
				if err != nil {
					fmt.Printf("\033[31mError getting current working directory: %v\033[0m\n", err)
					return
				}
				composeFilePath = dir
			}
			err = compose.ActivateCompose(composeFilePath, develMode)
			if err != nil {
				fmt.Printf("\033[31mCompose activation failed: %v\033[0m\n", err)
				return
			}

			err = compose.SetCrontabCompose()
			if err != nil {
				fmt.Printf("\033[31mError setting crontab for compose: %v\033[0m\n", err)
			} else {
				fmt.Printf("\033[32mCompose activated successfully with crontab set!\033[0m\n")
			}
		} else if activate_type == "standby" {
			err := service.InstallOperatorStandbymode()
			if err != nil {
				fmt.Printf("\033[31mStandby mode activation failed: %v\033[0m\n", err)
			} else {
				fmt.Printf("\033[32mStandby mode activated successfully!\033[0m\n")
			}
			return
		} else {
			hossted.SetUpdates(ENVIRONMENT, true)
			hossted.SetMonitoring(ENVIRONMENT, true)
			fmt.Printf("\033[32mUpdates and monitoring set successfully!\033[0m\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
	activateCmd.Flags().StringVarP(&activate_type, "type", "t", "", "Supported env type: k8s|compose|standby")
	activateCmd.Flags().StringVarP(&composeFilePath, "compose_filepath", "f", "", "Compose filepath (optional)")
	activateCmd.Flags().BoolVar(&develMode, "d", false, "Toggle development mode")
	activateCmd.Flags().BoolVarP(&verbose, "verbose", "V", false, "Enable verbose mode to send activation event")
}
