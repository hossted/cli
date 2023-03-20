package hossted

import (
	"fmt"
	"os/exec"
)

func SetMonitoring(env string, flag bool) error {

	config, _ := GetConfig()

	config.Monitoring = flag

	err := WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}
	fmt.Println("monitoring set to", flag)

	apps, err := GetAppInfo()
	if err != nil {
		return err
	}

	cmdArgs := []string{"--profile", "monitoring"}
	if flag {
		cmdArgs = append(cmdArgs, "up", "-d")
	} else {
		cmdArgs = append(cmdArgs, "down")
	}
	cmd := exec.Command("docker-compose", cmdArgs...)
	cmd.Dir = apps[0].AppPath

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	fmt.Println(string(output))

	if !flag {
		cmd = exec.Command("docker-compose", "up", "-d")
		cmd.Dir = apps[0].AppPath
		output, err = cmd.Output()
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		fmt.Println(string(output))
	}

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
