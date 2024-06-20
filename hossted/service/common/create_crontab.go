package common

import (
	"fmt"
	"os"
	"os/exec"
)

func CreateCrontab(command string, minute string) {
	cronTime := fmt.Sprintf("%s * * * *", minute)
	var cronJob string = "" + cronTime + " " + command

	cmd := exec.Command("bash", "-c", `crontab -l | grep -q '`+command+`'  && echo '`+command+` already running' || ((crontab -l; echo "`+cronJob+`") | crontab -  && echo '`+command+` configured.')`)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
