package hossted

import (
	"fmt"

	"github.com/spf13/viper"
)

// RegisterUsers updates email, organization, etc,.. in the yaml file
func RegisterUsers() error {
	s := viper.Get("cfgFile")
	fmt.Println(s)

	fmt.Println("Register User")
	return nil
}
