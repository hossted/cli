/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
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
	// Activate invokes update true and monitoring true
	hossted.SetUpdates(ENVIRONMENT,true)
	hossted.SetMonitoring(ENVIRONMENT, true)
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
