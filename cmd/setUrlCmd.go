/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setUrlCmdCmd represents the setUrlCmd command
var setUrlCmdCmd = &cobra.Command{
	Use:     "url",
	Short:   "[a] Set authorization of the provided application",
	Long:    `[a] Set authorization of the provided application`,
	Aliases: []string{"u"},
	Example: `
  hossted set url <AppName> example.com
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("setUrlCmd called")
		return nil
	},
}

func init() {
	setCmd.AddCommand(setUrlCmdCmd)
}
