package hossted

import (
	"fmt"
   	"time"
	"math/rand"
	"os"
	"os/exec"

)

func Schedule(env string) error {
	
	createCronSchedule()
	
	config, _ := GetConfig()

    currentTime := time.Now()
    yyyy, mm, dd := currentTime.Date()
    tomorrow := time.Date(yyyy, mm, dd+1, 0, 0, 0, 0, currentTime.Location())
	duration := tomorrow.Sub(currentTime)
	hours:=int(duration.Hours())

	hoursRand:=rand.Intn(hours)
	
	time.Sleep(time.Duration(hoursRand)*time.Hour) 
	if config.Update==true{ 
		Ping(env) //call hossted ping-send dockers info
	}

	return nil
}


func createCronSchedule() {

	var command=`crontab -l | grep -q 'hossted schedule'  && echo 'hossted schedule exists' || (crontab -l; echo "0 0 * * * /usr/local/bin/hossted schedule 2>&1 | logger -t mycmd") | crontab -`
	
	cmd:= exec.Command("bash", "-c", command)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err:= cmd.Start()
    if err != nil {
      fmt.Println(err)
    }
    err1 := cmd.Wait()
    if err1 != nil {
      fmt.Println(err1)
    }
}