/*
Copyright © 2022 Lior Kesos lior@hossted.com

*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "hossted",
		Short: "A brief description of your application.",
		Long: `
A brief description of your application
`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var err error
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hossted.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cfgFile, err = checkConfigFilePath()
	if err != nil {
		fmt.Println(err)
	}
}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hossted")

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	return
}

// checkConfigFilePath checks if the ~/.hossted/config.yaml is created under home folder
// Create it if it doesnt exist. Will create folder recursively
// TODO: put upder .hossted dir
func checkConfigFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	filePath := path.Join(home, ".hossted", "config.yaml")
	folder := path.Dir(filePath)

	if _, err := os.Stat(filePath); err != nil {

		// Create directory if not exists
		if _, err := os.Stat(folder); err != nil {
			os.MkdirAll(folder, os.ModePerm)
		}

		// Create file
		file, err := os.Create(filePath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		// Create file
		fmt.Printf("\nNo existing config file. \nNew config file is created  - %s \n\n", filePath)

		return "", err
	} else {
		// Normal case
		// fmt.Printf("\nUsing TODO file  - %s \n\n", path)
	}

	return filePath, nil
}
