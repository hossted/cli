/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setAuthCmd represents the setAuth command
var setAuthCmd = &cobra.Command{
	Use:     "auth",
	Short:   "[a] Set authorization of the provided application",
	Long:    `[a] Set authorization of the provided application`,
	Aliases: []string{"a"},
	Example: `
  hossted set auth <AppName> true
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("setAuth called")
		return nil
	},
}

func init() {
	setCmd.AddCommand(setAuthCmd)
}
