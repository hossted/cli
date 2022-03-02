/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:     "dev",
	Short:   "For Development only.",
	Long:    `For Development only.`,
	Aliases: []string{"x"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("dev called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devCmd)

}
