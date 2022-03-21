/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

var (
	input string
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "docker-compose ps",
	Long: `
docker-compose ps of the application
`,
	Example: `
  hossted ps
  hossted ps <app_name> (e.g. hossted ps wikijs)
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Check if the user provides the apps name manually
		var input string
		if len(args) >= 1 {
			input = args[0]
		}

		err := hossted.ListAppPS(input)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
