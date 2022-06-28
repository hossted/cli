package hossted

import (
	"fmt"
	"io/ioutil"
)

type Dummy struct {
	Head     []string
	Version  string
	Services []Service
	Bottom   []string
}

type Line struct {
	LineNum int
	Line    string
}

type Service struct {
	Name    string
	Content string
}

func Dev() error {

	fmt.Println("Dev")
	b, err := ioutil.ReadFile("docker-compose.yaml")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
