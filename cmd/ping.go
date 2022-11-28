/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: `hossted ping - Send docker ,sbom and security infor to hossted API`,
	Long:  `hossted ping - Send docker ,sbom and security infor to hossted API`,
    Example: `hossted ping`,
	Run: func(cmd *cobra.Command, args []string) {
		 hossted.Ping(ENVIRONMENT)
		
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
