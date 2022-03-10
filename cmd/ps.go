/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/hossted/hossted"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "docker-compose ps",
	Long: `
docker-compose ps
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := hossted.ListDockerPS()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
