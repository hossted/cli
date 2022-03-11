/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	COMMITHASH = ""
	BUILDTIME  = ""
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "[v] Get the version of the hossted CLI program",
	Long:    `[v] Get the version of the hossted CLI program`,
	Aliases: []string{"v"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("hossted version %s.\nBuilt on %s (%s)\n", VERSION, BUILDTIME, COMMITHASH)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
