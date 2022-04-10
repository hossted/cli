/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
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
  hossted httpopen <app_name> (e.g. hossted httpopen gitbucket)
  hossted h <app_name>
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Check if the user provides the apps name manually
		var input string
		if len(args) >= 1 {
			input = args[0]
		}
		err := hossted.HttpOpen(input)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(httpopenCmd)

}
