package hossted

import (
	"fmt"
	"time"
	"math/rand"
	"os"
	"os/exec"
)

func Schedule(env string) error {
	
	fmt.Println("schedule")
	createCronSchedule()
	
	config, _ := GetConfig()

    currentTime := time.Now()
    yyyy, mm, dd := currentTime.Date()
    tomorrow := time.Date(yyyy, mm, dd+1, 0, 0, 0, 0, currentTime.Location())
    fmt.Println("tomorrow",tomorrow)
	
	duration := tomorrow.Sub(currentTime)
	hours:=int(duration.Hours())
	fmt.Printf("difference %d hours\n",hours)

	hoursRand:=rand.Intn(hours)
	fmt.Println("hoursRand",hoursRand)

	
	time.Sleep(time.Duration(hoursRand)*time.Second)

	fmt.Println("config.Update",config.Update)
	if config.Update==true{
		fmt.Println("update")
		Ping(env)
	}

	return nil
}

func createCronSchedule() {
    fmt.Println("createCronSchedule")
	
	err := os.WriteFile("cronSchedule.sh", []byte("/home/linnovate/devel/hossted/cli/bin/linux/hossted schedule 2>&1 | logger -t mycmd"), 0755)
    if err != nil {
        fmt.Printf("Unable to write file: %v", err)
		return
    }

	cmd:= exec.Command(`sudo`,`mv`,`-u`,`cronSchedule.sh`, `/etc/cron.hourly`)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
}

