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
	RunE: func(cmd *cobra.Command, args []string) error {

		err := hossted.GetAppLogs()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
