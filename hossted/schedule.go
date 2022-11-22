package hossted

import (
	"fmt"
	"time"
	"math/rand"
)


func Schedule(env string) error {
	
	fmt.Println("schedule")
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
