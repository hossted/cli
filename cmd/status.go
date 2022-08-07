/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// statusCmd represents the version command
var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "[v] Print the status of the hossted instance variables",
	Long:    `[v] PRint the status of the hossted instance variables`,
	Aliases: []string{"st"},
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := hossted.GetConfigPath()
		if err != nil {
			return err
		}

		dat, err := os.ReadFile(configPath)

		fmt.Print(string(dat))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
