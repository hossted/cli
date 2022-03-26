package hossted

import (
	"fmt"
	"os"
)

func SetURL(app, url string) error {
	fmt.Println("Set URL")

	config, err := GetConfig()
	if err != nil {
		fmt.Printf("Please call the command `hossted register` first.\n%w", err)
		os.Exit(0)
	}
	_ = config

	return nil

}
