/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/hossted/cli/hossted/service"
	"github.com/hossted/cli/hossted/service/compose"

	"github.com/spf13/cobra"
)

var activate_type string

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
			err := service.ActivateK8s()
			if err != nil {
				fmt.Println(err)
			}
			return
		} else if activate_type == "compose" {
			// hossted.SetUpdates(ENVIRONMENT, true)
			// hossted.SetMonitoring(ENVIRONMENT, true)
			err := compose.ActivateCompose()
			if err != nil {
				fmt.Println(err)
			}
			return

		} else if activate_type == "standby" {
			err := service.InstallOperatorStandbymode()
			if err != nil {
				fmt.Println(err)
			}
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
	activateCmd.Flags().StringVarP(&activate_type, "type", "t", "", "supported env type k8s|docker")
}
