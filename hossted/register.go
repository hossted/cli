package hossted

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/spf13/hossted/utils"
	"github.com/spf13/viper"
)

var (
	//go:embed template/config.gohtml
	configTmpl []byte // config
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

// WriteDummyConfig writes the initial config to the ~/.hossted/config.yaml
func WriteDummyConfig() error {
	var config Config

	json.Unmarshal(configTmpl, &config)
	fmt.Println(utils.PrettyPrint(config))

	return nil
}
