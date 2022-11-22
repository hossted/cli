package hossted

import (
	"fmt"
)

func Updates(env string,flag bool) error {
	
	fmt.Println("updates")
	
	config, _ := GetConfig()
	fmt.Println("config",config)

	config.Update=flag

	err:= WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}
	return nil
}



