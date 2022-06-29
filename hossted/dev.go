package hossted

import (
	"fmt"
	"io/ioutil"
)

func Dev() error {

	fmt.Println("Dev")

	var ds DockerStruct // docker struct

	b, err := ioutil.ReadFile("docker-compose.yaml")
	if err != nil {
		return err
	}

	err = ds.Unmarshal([]byte(b))
	if err != nil {
		return err
	}

	// fmt.Println(PrettyPrint(ds))
	s, err := ds.Write()
	if err != nil {
		return err
	}
	fmt.Println(s)

	return nil
}
