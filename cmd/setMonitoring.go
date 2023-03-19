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

// monitoringCmd represents the monitoring command
var setMonitoringCmd = &cobra.Command{
	Use: "monitoring",
	Short: `[u] hossted set monitoring true - Allow to send information about dockers , 
	sbom and required security changes to the hossted API`,
	Long: ` [u] Allow to send information about dockers , 
	sbom and required security changes to the hossted API 
	so it can be displayed within the hossted dashboard and recommend the course of action to secure your hossted application.`,
	Aliases: []string{"u"},
	Example: `
  hossted set monitoring true
  hossted set monitoring false
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

		hossted.SetMonitoring(ENVIRONMENT, flag)

	},
}

func init() {
	setCmd.AddCommand(setMonitoringCmd)
}
