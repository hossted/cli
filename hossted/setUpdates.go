package hossted

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func SetUpdates(env string, flag bool) error {

	if !HasContainerRunning() {
		fmt.Println("The application still in configuration")
		os.Exit(0)
	}

	config, _ := GetConfig()

	config.Update = flag

	err := WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}
	fmt.Println("updates set to", flag)

	//send activity log about the command
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}
	fullCommand := "hossted set updates " + fmt.Sprint(flag)
	options := `{"updates":` + fmt.Sprint(flag) + `}`
	typeActivity := "set_updates"
	sendActivityLog(env, uuid, fullCommand, options, typeActivity)

	if flag == true {
		createCronSchedule()

	} else {
		stopCronSchedule()
	}

	return nil
}
func createCronSchedule() {
	now := time.Now().Add(1 * time.Minute)
	minute := now.Minute()
	if now.Second() > 30 {
		minute++
	}
	hour := now.Hour()
	if minute > 59 {
		minute = 0
		hour++
	}
	if hour > 23 {
		hour = 0
	}
	cronTime := fmt.Sprintf("%d %d * * *", minute, hour)
	var command string = "" + cronTime + " /usr/local/bin/hossted schedule 2>&1 | logger -t mycmd"
	cmd := exec.Command("bash", "-c", `crontab -l | grep -q 'hossted schedule'  && echo 'hossted updates already running' || (crontab -l; echo "`+command+`") | crontab -  && echo 'Hossted will now send package, security and monitoring information to the hossted api and will appear on the hossted dashboard.'`)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	err1 := cmd.Wait()
	if err1 != nil {
		fmt.Println(err1)
	}
}

func stopCronSchedule() {
	// create a new crontab without 'hossted schedule'
	cmd := exec.Command("bash", "-c", "crontab -l | grep -v 'hossted schedule' | crontab -")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error stopping cron job:", err)
		return
	}
}
