/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
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
		// TOLEARN - move cmd to config
		//curl -k -H "Authorization: Basic FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA" -X POST --compressed "https://app.dev.hossted.com/api/register?uuid=$UUID&email=$EMAIL&organization=$ORGANIZATION"

		//var registerAPI = "https://app.dev.hossted.com/api/register"
		//var hosstedAPIAuth = "Authorization: Basic FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA"
		// TOLEARN - Get from filesystem at /opt/linnovate

		uuid, software := getHosstedDetails()
		fmt.Println("register called, uuid:" + uuid)
		fmt.Println("register called, uuid:" + software)

		/*
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

func rewriteVersions() {
	//_
}

func getHosstedDetails() (string, string) {
	hossteDIR := "/opt/linnovate/"
	uuidFile := "uuid.txt"
	softFile := "software.txt"

	uuid, err := ioutil.ReadFile(hossteDIR + uuidFile)

	if err != nil {
		fmt.Println(err)
	}
	software, err := ioutil.ReadFile(hossteDIR + softFile)

	//fmt.Println("register called, software:" + string(software))
	if err != nil {
		fmt.Println(err)
	}
	//var hosstedir = viper.GetString("hosstedir") // case-insensitive Setting & Getting
	//fmt.Println("hosstedir:", hosstedir)
	return string(uuid), string(software)
}

// Get email and organization from user
//getDetails(){

//}

func init() {
	rootCmd.AddCommand(registerCmd)
	// Here you will define your flags and configuration settings.

	// TOLEARN Get hardcode URL value from the configuration file

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}
*/
