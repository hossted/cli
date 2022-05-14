/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
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
			err  error
		)

		if len(args) == 0 {
			return errors.New("\n  Empty Arguments.\n")

		} else if len(args) == 1 { // set auth true

			flag, err = hossted.ConvertBool(args[0])
			if err != nil {
				return err
			}

			config, err := hossted.GetConfig()
			if err != nil {
				return err
			}

			pwd := hossted.GetCurrentDirectory()
			app, _ = config.GetDefaultApp(pwd)

		} else if len(args) == 2 { // set auth appname true
			app = args[0]
			flag, err = hossted.ConvertBool(args[1])
			if err != nil {
				return err
			}
		}
		if strings.TrimSpace(app) == "" {
			return fmt.Errorf("No input application.")
		}

		err = hossted.SetAuth(app, flag)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	setCmd.AddCommand(setAuthCmd)
}
