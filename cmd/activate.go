/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/hossted/cli/hossted"
	"github.com/hossted/cli/hossted/service"
	"github.com/hossted/cli/hossted/service/compose"

	"github.com/spf13/cobra"
)

var (
	activate_type, releaseName, composeFilePath, token, org_id string
	develMode bool
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

		if activate_type == "k8s" {
			err := service.ActivateK8s(releaseName, token, org_id)
			if err != nil {
				fmt.Println(err)
			}
			return
		} else if activate_type == "compose" {
			// hossted.SetUpdates(ENVIRONMENT, true)
			// hossted.SetMonitoring(ENVIRONMENT, true)
			err := compose.ActivateCompose(composeFilePath, token, org_id, develMode)
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
	activateCmd.Flags().StringVarP(&token, "token", "", "", "token for orgID")
	activateCmd.Flags().StringVarP(&org_id, "org_id", "", "", "orgID")
	activateCmd.Flags().StringVar(&releaseName, "release_name", "", "release name (optional)")
	activateCmd.Flags().StringVar(&composeFilePath, "compose_filepath", "", "compose filepath (optional)")
	activateCmd.Flags().BoolVar(&develMode, "d", false, "Toggle development mode")
}
