/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var updateCmd = &cobra.Command{
	Use:     "update-cli",
	Short:   `hossted update-cli - update hossted cli`,
	Long:    `hossted update-cli - update hossted cli`,
	Example: `hossted update-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		hossted.Update(ENVIRONMENT)

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
