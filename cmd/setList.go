/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// setListCmd represents the setList command
var setListCmd = &cobra.Command{
	Use:     "list",
	Short:   "[l] List all the commands of the available applications",
	Long:    "[l] List all the commands of the available applications",
	Aliases: []string{"l"},
	Example: `
  hossted set list
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := hossted.ListCommands()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	setCmd.AddCommand(setListCmd)

}
