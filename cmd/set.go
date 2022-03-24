/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:     "set",
	Short:   "[s] Change application settings",
	Long:    `[s] Change application settings`,
	Aliases: []string{"s"},
	Example: `
  hossted set auth <AppName> true
  hossted set url <AppName> linovate
  hossted set ssl <AppName> sign
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("set called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
