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
	Short: "[a] Activate your application with the hossted ecosystem",
	Long: `
The hossted register commands asks for your email and organization name
and registers you with the hossted platfrom

The hossted platform provides secure and hardened docker images and provides
best practices such as tracking updates , monitoring, centralized logging ,
backups and much more.
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
