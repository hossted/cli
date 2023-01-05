/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

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

		if len(args) < 1 {
			cmd.Help()
			os.Exit(0)
		}

		// Parse input
		var flag bool
		input := strings.ToLower(args[0])
		if input == "true" {
			flag = true
		} else if input == "false" {
			flag = false
		} else {
			return fmt.Errorf("Only true/false is supported. Input - %s\n", input)
		}

		err := hossted.SetRemoteAccess(ENVIRONMENT, flag)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	setCmd.AddCommand(remoteSupportCmd)

}
