/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/hossted/cli/hossted/service"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:     "import k8s | docker",
	Short:   `hossted import - Import existing app and send info to hossted API`,
	Long:    `hossted import - Import existing app and send info to hossted API`,
	Example: `hossted import k8s`,
	Run: func(cmd *cobra.Command, args []string) {
		if import_type == "k8s" {
			err := service.ActivateK8s()
			if err != nil {
				fmt.Println(err)
			}
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
	import_type string
	Uuid        string
	User        string
	Cloud       string
	Status      string
	Test_mode   string
	Comment     string
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&import_type, "type", "t", "", "supported env type k8s|docker")
	importCmd.Flags().StringVarP(&uuid, "uuid", "", "", "")
	importCmd.Flags().StringVarP(&ip, "ip", "", "", "")
	importCmd.Flags().StringVarP(&user, "user", "", "", "")
	importCmd.Flags().StringVarP(&url, "url", "", "", "")
	importCmd.Flags().StringVarP(&product, "product", "", "", "")
	importCmd.Flags().StringVarP(&cpunum, "cpunum", "", "", "")
	importCmd.Flags().StringVarP(&mem, "mem", "", "", "")
	importCmd.Flags().StringVarP(&rootsize, "rootsize", "", "", "")
	importCmd.Flags().StringVarP(&cloud, "cloud", "", "", "")
	importCmd.Flags().StringVarP(&image, "image", "", "", "")
	importCmd.Flags().StringVarP(&status, "status", "", "", "")
	importCmd.Flags().StringVarP(&authorization, "authorization", "", "", "")
}
