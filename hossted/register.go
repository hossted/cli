package hossted

import (
	"fmt"

	"github.com/spf13/viper"
)

// RegisterUsers updates email, organization, etc,.. in the yaml file
func RegisterUsers() error {

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	s := viper.Get("martin")
	fmt.Println(fmt.Sprintf("Getting config file - %s", s))

	fmt.Println("Register User")
	return nil
}
