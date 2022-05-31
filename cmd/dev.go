/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

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

	// cmd := exec.Command("cat", "/tmp/test.txt")
	// out, err := cmd.Output()
	// if err != nil {
	// 	return err
	// }
	s := string(`
abcde
hello world
     - "traefik.http.middlewares.tauth.basicauth.usersfile=letsencrypt/.htpass"
     - "traefik.http.routers.$PROJECT_NAME.middlewares=tauth"
`)
	pattern := []string{"tauth.basicauth.usersfile"}
	flag := true
	matchOnce := true
	strict := false

	res, err := ToggleCommentLinesByRegex(s, pattern, flag, matchOnce, strict)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

// ToggleCommentLinesByRegex comments/uncomment out lines by a list of regular expressions
// patterns as a list of regular expresson to be checked.
// flag specifies whether it is comment/uncomment of all the lines. true as comment out
// matchOnce specifies the pattern will stop at first match, or match the whole files
// strict mode will returns error as soon as the pattern is not found in the input string
func ToggleCommentLinesByRegex(s string, patterns []string, flag bool, matchOnce bool, strict bool) (string, error) {

	// Split lines
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")

	return "", nil
}

func init() {
	rootCmd.AddCommand(devCmd)

}
