/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register your application with the hossted ecosystem",
	Long: `The hossted register commands asks for your email and organization name 
	and registers you with the hossted platfrom

The hossted platform provides secure and hardened docker images and provides 
best practices such as tracking updates , monitoring, centralized logging , backups and much more `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("register called")

		/*
				if len(args) >= 1 && args[0] != "" {
					var Arguments = args[0]
				}

			for _, e := range os.Environ() {
				pair := strings.SplitN(e, "=", 2)
				fmt.Println(pair[0])
			}

				URL := "https://github.com/scraly/gophers/raw/main/" + gopherName + ".png"

				fmt.Println("Try to get '" + gopherName + "' Gopher...")

				// Get the data
				response, err := http.Get(URL)
				if err != nil {
					fmt.Println(err)
				}
				defer response.Body.Close()

				if response.StatusCode == 200 {
					// Create the file
					out, err := os.Create(gopherName + ".png")
					if err != nil {
						fmt.Println(err)
					}
					defer out.Close()

					// Writer the body to file
					_, err = io.Copy(out, response.Body)
					if err != nil {
						fmt.Println(err)
					}

					fmt.Println("Perfect! Just saved in " + out.Name() + "!")
				} else {
					fmt.Println("Error: " + gopherName + " not exists! :-(")
				}
		*/
	},
}

func getHosstedDetails() {
	var hosstedir = viper.GetString("hosstedir") // case-insensitive Setting & Getting
	fmt.Println("hosstedir:", hosstedir)
}

func init() {
	rootCmd.AddCommand(registerCmd)
	// Here you will define your flags and configuration settings.
	getHosstedDetails()

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
