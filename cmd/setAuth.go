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

// setAuthCmd represents the setAuth command
var setAuthCmd = &cobra.Command{
	Use:     "auth",
	Short:   "[a] Set authorization of the provided application",
	Long:    `[a] Set authorization of the provided application`,
	Aliases: []string{"a"},
	Example: `
  hossted set auth <AppName> false
  hossted set auth <AppName> true
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Check if the user provides the apps name manually
		var (
			app  string
			flag bool
		)

		if len(args) < 1 {
			cmd.Help()
			os.Exit(0)
		}

		if len(args) == 1 { // set auth true
			flag, err := hossted.ConvertBool(args[0])
			if err != nil {
				return err
			}

			config, err := hossted.GetConfig()
			if err != nil {
				return err
			}

			pwd := hossted.GetCurrentDirectory()
			app, _ = config.GetDefaultApp(pwd)

		} else if len(args) == 2 {
			app = args[0]
			flag = args[1]
		}
		if strings.TrimSpace(app) == "" {
			return fmt.Errorf("No input application.")
		}
		err := hossted.HttpOpen(app)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	setCmd.AddCommand(setAuthCmd)
}
