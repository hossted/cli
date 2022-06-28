package hossted

import (
	"fmt"
	"io/ioutil"
)

func Dev() error {

	fmt.Println("Dev")
	b, err := ioutil.ReadFile("docker-compose.yaml")
	if err != nil {
		return err
	}
	fmt.Println(len(b))

	return nil
}
