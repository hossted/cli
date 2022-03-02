/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/hossted/hossted"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register your application with the hossted ecosystem",
	Long: `The hossted register commands asks for your email and organization name
	and registers you with the hossted platfrom

The hossted platform provides secure and hardened docker images and provides
best practices such as tracking updates , monitoring, centralized logging , backups and much more.`,
	Aliases: []string{"r"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := hossted.RegisterUsers()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
