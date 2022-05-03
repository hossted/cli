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

// setDomainCmd represents the setUrlCmd command
// hossted set url <AppName> example.com
var setDomainCmd = &cobra.Command{
	Use:     "domain",
	Short:   "[d] Set the domain of the provided application",
	Long:    "[d] Set the domain of the provided application",
	Aliases: []string{"d"},
	Example: `
  hossted set domain <AppName> example.com
  hossted set domain prometheus example.com
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var app, domain string

		if len(args) == 1 {
			config, err := hossted.GetConfig()
			if err != nil {
				return err
			}

			pwd := hossted.GetCurrentDirectory()
			app, _ = config.GetDefaultApp(pwd)
			domain = args[0]

		} else if len(args) == 2 {
			app = args[0]
			domain = args[1]

		}

		if strings.TrimSpace(app) == "" {
			fmt.Printf("Not enough arguments. Expecting <AppName> <Domain>\nPlease checking the input params. %v\n\n", args)
			cmd.Help()
			os.Exit(0)
		}

		err := hossted.SetDomain(app, domain)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	setCmd.AddCommand(setDomainCmd)
}
