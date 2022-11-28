package hossted

import (
	"fmt"
	//"os/exec"
)

func SetUpdates(flag bool) error {
	
	fmt.Println("updates")
	
	config, _ := GetConfig()
	fmt.Println("config",config)

	config.Update=flag

	err:= WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}

	fmt.Println("----")
	// cmd := exec.Command(`hossted`, `schedule`)
	// dataOutput, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println("ddd",string(dataOutput))
	return nil
}



