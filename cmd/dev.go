/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := dev()
		if err != nil {
			return err
		}

		return nil
	},
}

func dev() error {

	cmd := exec.Command("cat", "/tmp/test.txt")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return nil
}

// ToggleCommentLinesByRegex comments/uncomment out lines by a list of regular expressions
// patterns as a list of regular expresson to be checked.
// toggle specifies whether it is comment/uncomment of all the lines
// matchOnce specifies the pattern will stop at first match, or match the whole files
// strict mode will returns error as soon as the pattern is not found in the input string
func ToggleCommentLinesByRegex(s string, patterns []string, toggle string, matchOnce bool, strict bool) (string, error) {

	return "", nil
}

func init() {
	rootCmd.AddCommand(devCmd)

}
