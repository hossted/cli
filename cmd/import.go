/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: `hossted import - Import existing app and send info to hossted API`,
	Long:  `hossted import - Import existing app and send info to hossted API`,
    Example: `hossted import k8s`,
	Run: func(cmd *cobra.Command, args []string) {
		 hossted.Import(ENVIRONMENT)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
