package hossted

import (
	"fmt"
)

func SetUpdates(env string, flag bool) error {
	
	config, _ := GetConfig()

	config.Update=flag

	err:= WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}
	
	Schedule(env) //call hossted schedule
	
	return nil
}



