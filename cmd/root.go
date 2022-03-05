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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cfgFile, err = checkConfigFilePath()
	_ = cfgFile

	if err != nil {
		fmt.Println(err)
	}
}

func initConfig() {

	// Not allowed to change cfg path anyway
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	folder := fmt.Sprintf("%s/%s", home, ".hossted")
	viper.AddConfigPath(folder)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	return
}

// checkConfigFilePath checks if the ~/.hossted/config.yaml is created under home folder
// Create it if it doesnt exist. Will create folder recursively. Also it will init the config file yaml.
// TODO: init yaml file
func checkConfigFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	// Derive config path, and .hossted folder. Under user home
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

		// Write init config from template

		// Create file
		fmt.Printf("\nNo existing config file. \nNew config file is created  - %s \n\n", filePath)

		return "", err
	} else {
		// Normal case
		// Do nothing. config.yaml exists
	}

	return filePath, nil
}
