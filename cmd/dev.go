/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"regexp"
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
     - "traefik.http.routers.$PROJECT_NAME.middlewares=tauth"`)

	pattern := []string{
		`.*tauth\.basicauth\.usersfile.*`,
		`.*middlewares=tauth.*`,
	}
	flag := true
	matchOnce := false
	strict := false

	res, err := ToggleCommentLinesByRegex(&s, pattern, flag, matchOnce, strict)
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
func ToggleCommentLinesByRegex(s *string, patterns []string, flag bool, matchOnce bool, strict bool) (string, error) {

	// Compile regex
	var rePatterns []*regexp.Regexp // List of regexp pattern
	for _, p := range patterns {
		pat, err := regexp.Compile(p)
		if err != nil {
			return "", err
		}
		rePatterns = append(rePatterns, pat)
	}

	// For each pattern, match with the string
	for _, re := range rePatterns {
		_ = re
		matches := re.FindAllString(*s, -1)

		if len(matches) == 0 { // No matching pattern
			if strict { // For strict mode, return as soon as we dont find a matching line
				return "", errors.New(fmt.Sprintf("Patten not found - %s", re))
			} else {
				// TODO: handle non-strict mode, only return error if there is one
				return "", errors.New("Not implemented. Please check.")
			}
		} else {
			fmt.Sprintf("matched - %s", re)
		}

		if matchOnce {
			// pass for now
			return "", errors.New("Not implemented. Please check.")
		} else {
			// match all matching lines
			for _, line := range matches {
				newline := toggleComment(line, flag)
				_ = newline
			}
			*s = re.ReplaceAllString(*s, "")
		}

	}
	fmt.Println(*s)

	return "", nil
}

// toggleComment toggles to comment/uncommented a line
// flag true as commenting the line
func toggleComment(line string, flag bool) string {
	var (
		trimmedLine string // For condition checking only
		newline     string // result string
	)
	trimmedLine = strings.TrimSpace(line)
	leadingChar := trimmedLine[0:1]

	if flag { // To Comment

		if leadingChar == "#" {
			return line
		} else {
			// Check first character of original line, whether it's a space
			if line[0:1] == " " {
				newline = "#" + line[1:]
			} else {
				newline = "#" + line
			}

		}

	} else { // To Uncomment

		if leadingChar != "#" {
			return line
		} else {
			newline = strings.Replace(line, "#", " ", 1)
		}
	}

	return newline
}

func init() {
	rootCmd.AddCommand(devCmd)

}
