package hossted

import (
	"fmt"
)

func SetMonitoring(env string, flag bool) error {

	config, _ := GetConfig()

	config.Update = flag

	// err := WriteConfigWrapper(config)
	// if err != nil {
	// 	return fmt.Errorf("Can not write to config file. Please check. %w", err)
	// }
	fmt.Println("monitoring set to", flag)

	//send activity log about the command
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}
	fullCommand := "hossted set monitoring " + fmt.Sprint(flag)
	options := `{"monitoring":` + fmt.Sprint(flag) + `}`
	typeActivity := "set_monitoring"
	sendActivityLog(env, uuid, fullCommand, options, typeActivity)

	return nil
}
