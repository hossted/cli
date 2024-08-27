/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
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
	activate_type, releaseName, composeFilePath, org_id string
	develMode                                           bool
)

// registerCmd represents the register command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "[a] Activate your application with the hossted platform",
	Long: `
Hossted activate connects you're instance to the hossted platform and sends instance health informtation so it can be mnaged in the dashboard.
	`,
	Aliases: []string{"a"},
	Example: `
hossted activate
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// write activate_type to config file
		config, _ := hossted.GetConfig() // Ignore error
		// Assign back to config object
		config.ActivateType = activate_type

		// Write back to file
		err := hossted.WriteConfigWrapper(config)
		if err != nil {
			fmt.Println("Can not write ActivateType to config file. Please check.", err)
			return
		}

		if activate_type == "k8s" {
			err = service.VerifyAuth(develMode)
			if err != nil {
				fmt.Println("Auth verification is failed, error:", err)
				return
			}

			err := service.ActivateK8s(releaseName, develMode)
			if err != nil {
				fmt.Println(err)
			}
			return
		} else if activate_type == "compose" {
			err = service.VerifyAuth(develMode)
			if err != nil {
				fmt.Println("Auth verification is failed, error:", err)
				return
			}

			if composeFilePath == "" {
				dir, err := os.Getwd()
				if err != nil {
					fmt.Println("Error getting current working directory:", err)
				}
				composeFilePath = dir
			}
			err := compose.ActivateCompose(composeFilePath, develMode)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = compose.SetCrontabCompose()
			if err != nil {
				fmt.Println("error in setting crontab for compose: ", err)
				return
			}
		} else if activate_type == "standby" {
			err := service.InstallOperatorStandbymode()
			if err != nil {
				fmt.Println(err)
			}
			return
		} else {
			hossted.SetUpdates(ENVIRONMENT, true)
			hossted.SetMonitoring(ENVIRONMENT, true)
		}

	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
	activateCmd.Flags().StringVarP(&activate_type, "type", "t", "", "supported env type k8s|compose")
	activateCmd.Flags().StringVar(&releaseName, "release_name", "", "release name (optional)")
	activateCmd.Flags().StringVarP(&composeFilePath, "compose_filepath", "f", "", "compose filepath (optional)")
	activateCmd.Flags().BoolVar(&develMode, "d", false, "Toggle development mode")
}
