/*
Copyright © 2022 Lior Kesos lior@hossted.com

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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

	cfgFile, err = checkFilePath()
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

// checkFilePath checks if the ~/.hossted.yaml is created under root folder
// Create it if it doesnt exist
// TODO: put upder .hossted dir
func checkFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, ".hossted.yaml")

	if _, err := os.Stat(path); os.IsNotExist(err) {

		file, err := os.Create(path)
		if err != nil {
			return "", err
		}
		defer file.Close()

		// Create file
		fmt.Printf("\nNo existing config file. \nNew config file is created  - %s \n\n", path)
		return "", err
	} else {
		// Normal case
		// fmt.Printf("\nUsing TODO file  - %s \n\n", path)
	}

	return path, nil
}
