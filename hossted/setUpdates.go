package hossted

import (
	"fmt"
)

func SetUpdates(env string, flag bool) error {

	config, _ := GetConfig()

	config.Update = flag

	err := WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}
	fmt.Println("updates set to", flag)

	Schedule(env) //call hossted schedule

	//send activity log about the command
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}
	fullCommand := "hossted set updates " + fmt.Sprint(flag)
	options := `{"updates":` + fmt.Sprint(flag) + `}`
	typeActivity := "set_updates"
	sendActivityLog(env, uuid, fullCommand, options, typeActivity)

	return nil
}
