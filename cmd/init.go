package cmd

import (
	"fmt"
	"os"

	"github.com/hossted/cli/hossted"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   `hossted init - Send instance hossted API`,
	Long:    `hossted init -  Send instance hossted API`,
	Example: `hossted init --uuid deefr455844-555-ttt6rr`,
	Run: func(cmd *cobra.Command, args []string) {

		if uuid == "" {
			fmt.Printf("\033[0;31m Error:\033[0m uuid is required.\n")
			os.Exit(0)
		}
		instance := hossted.Instance{
			Uuid:      uuid,
			Ip:        ip,
			User:      user,
			Url:       url,
			Product:   product,
			Cpunum:    cpunum,
			Mem:       mem,
			Rootsize:  rootsize,
			Cloud:     cloud,
			Status:    status,
			Test_mode: test_mode,
			Comment:   comment,
		}

		hossted.Init(ENVIRONMENT, authorization, image, instance)
	},
}

var (
	uuid          string
	ip            string
	user          string
	url           string
	product       string
	cpunum        string
	mem           string
	rootsize      string
	cloud         string
	image         string
	authorization string
	status        string
	test_mode     string
	comment       string
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&uuid, "uuid", "", "", "")
	initCmd.Flags().StringVarP(&ip, "ip", "", "", "")
	initCmd.Flags().StringVarP(&user, "user", "", "", "")
	initCmd.Flags().StringVarP(&url, "url", "", "", "")
	initCmd.Flags().StringVarP(&product, "product", "", "", "")
	initCmd.Flags().StringVarP(&cpunum, "cpunum", "", "", "")
	initCmd.Flags().StringVarP(&mem, "mem", "", "", "")
	initCmd.Flags().StringVarP(&rootsize, "rootsize", "", "", "")
	initCmd.Flags().StringVarP(&cloud, "cloud", "", "", "")
	initCmd.Flags().StringVarP(&image, "image", "", "", "")
	initCmd.Flags().StringVarP(&status, "status", "", "", "")
	initCmd.Flags().StringVarP(&authorization, "authorization", "", "", "")

}
