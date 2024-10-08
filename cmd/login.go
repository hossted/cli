/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hossted/cli/hossted/service"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "[l] To login in hossted platform",
	Long: `
Login enables user to login into hossted platform using email address`,
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prompt user for input
		err := service.Login(develMode)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().BoolVar(&develMode, "d", false, "Toggle development mode")
}
