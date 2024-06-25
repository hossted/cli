/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// reconcileComposeCmd represents the reconcileCompose command
var reconcileComposeCmd = &cobra.Command{
	Use:     "reconcile-compose",
	Short:   `[rc] hossted set reconcile compose - Set reconciler for the compose`,
	Long:    ` [rc] hossted set reconcile compose - Set reconciler for the compose`,
	Aliases: []string{"rc"},
	Example: `
  hossted reconcile-compose
	`,
	Run: func(cmd *cobra.Command, args []string) {
		hossted.ReconcileCompose()

	},
}

func init() {
	rootCmd.AddCommand(reconcileComposeCmd)
}
