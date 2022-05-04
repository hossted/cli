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

var (
	input string
)

// psCmd represents the ps command
// hossted ps <appName>
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "[ps] docker-compose ps of the application",
	Long:  "[ps] docker-compose ps of the application",
	Example: `
  hossted ps
  hossted ps <app_name> (e.g. hossted ps wikijs)
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

		err := hossted.ListAppPS(app)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
