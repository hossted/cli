/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// scanCmd represents the ps command
// hossted ps <appName>
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "[scan] docker-compose ps of the application",
	Long:  "[scan] docker-compose ps of the application",
	Example: `
  hossted scan
  hossted scan images
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// check number of parameters
		if len(args) > 1 {
			return fmt.Errorf("Too many parameters.")
		} else {
			hossted.ScanApps("images")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
