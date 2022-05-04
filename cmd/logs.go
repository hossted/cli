/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

var followFlag bool

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:     "logs",
	Aliases: []string{"l"},
	Short:   "[l] View Application logs",
	Long:    `[l] View Application logs`,
	Example: `
  hossted logs
  hossted logs <app_name> (e.g. hossted logs wikijs)
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Check if the user provides the apps name manually
		var app string
		if len(args) == 0 {
			config, err := hossted.GetConfig()
			if err != nil {
				return err
			}

			pwd := hossted.GetCurrentDirectory()
			app, _ = config.GetDefaultApp(pwd)

		} else if len(args) == 1 {
			app = args[0]
		}
		if strings.TrimSpace(app) == "" {
			return fmt.Errorf("No input application.")
		}

		err := hossted.GetAppLogs(app, followFlag)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolVarP(&followFlag, "follow", "f", false, "Follow log output.")
}
