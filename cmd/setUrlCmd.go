/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// setURLCmdCmd represents the setUrlCmd command
var setURLCmdCmd = &cobra.Command{
	Use:     "url",
	Short:   "[a] Set authorization of the provided application",
	Long:    `[a] Set authorization of the provided application`,
	Aliases: []string{"u"},
	Example: `
  hossted set url <AppName> example.com
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := hossted.SetURL()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	setCmd.AddCommand(setURLCmdCmd)
}
