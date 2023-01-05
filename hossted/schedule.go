package hossted

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func Schedule(env string) error {

	createCronSchedule()

	config, _ := GetConfig()

	if config.Update == true {
		Ping(env) //call hossted ping-send dockers info
	}

	return nil
}

func createCronSchedule() {

	rand.Seed(time.Now().UnixNano())
	hoursRand := rand.Intn(24)

	var command string = "0 " + strconv.Itoa(hoursRand) + " * * * /usr/local/bin/hossted schedule 2>&1 | logger -t mycmd"
	cmd := exec.Command("bash", "-c", `crontab -l | grep -q 'hossted schedule'  && echo 'hossted schedule exists' || (crontab -l; echo "`+command+`") | crontab -`)

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
