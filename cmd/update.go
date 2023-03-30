/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var updateCmd = &cobra.Command{
	Use:     "update-cli",
	Short:   `hossted update-cli - update the hossted cli to the last version`,
	Long:    `hossted update-cli - update the hossted cli to the last version`,
	Example: `sudo hossted update-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("SUDO_USER") == "" {
			fmt.Println("Please run this command with sudo.")
			os.Exit(0)
		}
		hossted.Update(ENVIRONMENT)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
