/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:     "logs",
	Aliases: []string{"l"},
	Short:   "[l] View Applicatin logs",
	Long:    `[l] View Applicatin logs`,
	Example: `
  hossted logs
  hossted logs <app_name> (e.g. hossted ps wikijs)
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Check if the user provides the apps name manually
		var input string
		if len(args) >= 1 {
			input = args[0]
		}

		err := hossted.GetAppLogs(input)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
