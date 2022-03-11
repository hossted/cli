/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:     "dev",
	Short:   "[x] For Development only",
	Long:    `[x] For Development only`,
	Aliases: []string{"x"},
	RunE: func(cmd *cobra.Command, args []string) error {

		err := hossted.Dev()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(devCmd)

}
