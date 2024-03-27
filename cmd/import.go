/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/hossted/cli/hossted"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import k8s | docker",
	Short: `hossted import - Import existing app and send info to hossted API`,
	Long:  `hossted import - Import existing app and send info to hossted API`,
    Example: `hossted import k8s`,
	Run: func(cmd *cobra.Command, args []string) {
		if t == "k8s" {
			// err := isValidEmail(args[1])
			// if err == false {
			// 	fmt.Println("Invalid email")
			// }
			// validate := func(input string) error {
			// 	fmt.Println("execution")
			// 	err := strings.Contains(input, "@") && strings.Contains(input, ".")
			// 	if err != false {
			// 		fmt.Println("Invalid email")
			// 		os.Exit(0)
			// 	}
			// 	return nil
			// }
			prompt := promptui.Prompt{
				Label: "email",
			}
			email_id, _ := prompt.Run()
			contexts := hossted.GetK8sContext()

			promptk8s := promptui.Select{
				Label: "Select kubernetes cluster",
				Items: contexts,
			}
			_, cluster, _ := promptk8s.Run()
			// send request to API
			fmt.Println("email is", email_id, "cluster selected is", cluster)
		}

		/*
		Uuid      string `json:"uuid,omitempty"`
	User      string `json:"user,omitempty"`
	Cloud     string `json:"cloud,omitempty"`
	Status    string `json:"status,omitempty"`
	Test_mode string `json:"test_mode,omitempty"`
	Comment   string `json:"comment,omitempty"`
		kluster := hossted.KCluster{
			Uuid:      uuid,
			Cloud:     cloud,
			Status:    status,
			Test_mode: test_mode,
			Comment:   comment,
		}
		authorization = "FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA"
		hossted.Import(ENVIRONMENT, authorization, kluster)

			*/

	},
}

var (
	t string
	Uuid      string 
	User      string 
	Cloud     string 
	Status    string 
	Test_mode string 
	Comment   string 
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&t, "type", "t", "", "supported env type k8s|docker")

	// importCmd.Flags().StringVarP(&uuid, "uuid", "", "", "")
	// importCmd.Flags().StringVarP(&ip, "ip", "", "", "")
	// importCmd.Flags().StringVarP(&user, "user", "", "", "")
	// importCmd.Flags().StringVarP(&url, "url", "", "", "")
	// importCmd.Flags().StringVarP(&product, "product", "", "", "")
	// importCmd.Flags().StringVarP(&cpunum, "cpunum", "", "", "")
	// importCmd.Flags().StringVarP(&mem, "mem", "", "", "")
	// importCmd.Flags().StringVarP(&rootsize, "rootsize", "", "", "")
	// importCmd.Flags().StringVarP(&cloud, "cloud", "", "", "")
	// importCmd.Flags().StringVarP(&image, "image", "", "", "")
	// importCmd.Flags().StringVarP(&status, "status", "", "", "")
	// importCmd.Flags().StringVarP(&authorization, "authorization", "", "", "")
}
