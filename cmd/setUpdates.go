/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"strings"
	"fmt"
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// updatesCmd represents the updates command
var setUpdatesCmd = &cobra.Command{
	Use:   "updates",
	Short: `[u] hossted set updates true - Allow to send information about dockers , 
	sbom and required security changes to the hossted API`,
	Long: ` [u] Allow to send information about dockers , 
	sbom and required security changes to the hossted API 
	so it can be displayed within the hossted dashboard and recommend the course of action to secure your hossted application.`,
    Aliases: []string{"u"},
    Example: `
  hossted set updates true
  hossted set updates false
`,

	Run: func(cmd *cobra.Command, args []string) {

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
			fmt.Printf("\033[0;31m Only true/false is supported.")
			fmt.Printf(" Input - %s\n\033[0m", input)
			os.Exit(0)
		}

		hossted.SetUpdates(ENVIRONMENT,flag)
		
	},
}

func init() {
	setCmd.AddCommand(setUpdatesCmd)
}
