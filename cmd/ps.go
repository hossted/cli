/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/hossted/hossted"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "tba",
	Long: `
TBA
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := hossted.ListPS()
		if err != nil {
			return err
		}
		fmt.Println("ps called")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
