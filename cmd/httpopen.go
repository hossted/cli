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

// httpopenCmd represents the httpopen command
var httpopenCmd = &cobra.Command{
	Use:     "httpopen",
	Short:   "[h] httpopen appname",
	Long:    "[h] httpopen appname",
	Aliases: []string{"h"},
	Example: `
  hostted httponen
  hossted httpopen <app_name> (e.g. hossted httpopen gitbucket)
  hossted h <app_name>
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
		err := hossted.HttpOpen(app)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(httpopenCmd)

}
