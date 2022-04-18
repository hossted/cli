/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// remoteSupportCmd represents the remoteSupport command
var remoteSupportCmd = &cobra.Command{
	Use:     "remote-support",
	Short:   "[r] To enable or disable remote ssh access with our maintanece and support key",
	Long:    "[r] To enable or disable remote ssh access with our maintanece and support key",
	Aliases: []string{"r"},
	Example: `
  hossted set remote-support true
  hossted set remote-support false
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := hossted.SetRemoteAccess(false)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	setCmd.AddCommand(remoteSupportCmd)

}
