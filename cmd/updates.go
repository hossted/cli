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
var updatesCmd = &cobra.Command{
	Use:   "updates",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
			fmt.Printf("Only true/false is supported. Input - %s\n", input)
			os.Exit(0)
		}

		hossted.Updates(ENVIRONMENT, flag)
		
	},
}

func init() {
	rootCmd.AddCommand(updatesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
