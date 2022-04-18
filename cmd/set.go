/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:     "set",
	Short:   "[s] Change application settings",
	Long:    `[s] Change application settings`,
	Aliases: []string{"s"},
	Example: `
  hossted set list
  hossted set auth <AppName> true
  hossted set domain <AppName> example.com
  hossted set ssl <AppName> sign
  hossted set remote-support true
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
