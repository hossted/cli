/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// statusCmd represents the version command
var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "[v] Print the status of the hossted instance variables",
	Long:    `[v] PRint the status of the hossted instance variables`,
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := hossted.GetConfig()
		if err != nil {
			return err
		}
		fmt.Printf("email :%s\nuuid: %s", config.Email, config.HostUUID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
