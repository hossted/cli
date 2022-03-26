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

// setURLCmd represents the setUrlCmd command
// hossted set url <AppName> example.com
var setURLCmd = &cobra.Command{
	Use:     "url",
	Short:   "[u] Set URL endpoints of the provided application",
	Long:    "[u] Set URL endpoints of the provided application",
	Aliases: []string{"u"},
	Example: `
  hossted set url <AppName> example.com
  hossted set url prometheus example.com
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 2 {
			fmt.Println("Not enough arguments. Please blah blah.")
			cmd.Help()
			os.Exit(0)
		}

		app := args[0]
		url := args[1]

		err := hossted.SetURL(app, url)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	setCmd.AddCommand(setURLCmd)
}
