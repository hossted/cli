/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "[r] Register your application with the hossted ecosystem",
	Long: `
The hossted register commands asks for your email and organization name
and registers you with the hossted platfrom

The hossted platform provides secure and hardened docker images and provides
best practices such as tracking updates , monitoring, centralized logging ,
backups and much more.
`,
	Aliases: []string{"r"},
	Example: `
  hossted register
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := hossted.RegisterUsers(ENVIRONMENT)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
