package hossted

import (
	"bufio"
	"embed"

	"fmt"
	"html/template"
	"io"

	"github.com/spf13/viper"
)

var (
	//go:embed templates
	templates embed.FS
)

// RegisterUsers updates email, organization, etc,.. in the yaml file
func RegisterUsers() error {

	viper.ConfigFileUsed()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("here")
		return err
	}

	s := viper.Get("email")
	fmt.Println(fmt.Sprintf("Getting config file - %s", s))

	fmt.Println("Register User")
	return nil
}

// WriteConfig writes the config to the config file (~/.hossted/config.yaml)
func WriteConfig(w io.Writer, config Config) error {

	// Read Template
	t, err := template.ParseFS(templates, "templates/config.tmpl")
	if err != nil {
		return err
	}

	// Write to template
	err = t.Execute(w, config)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(w)
	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
