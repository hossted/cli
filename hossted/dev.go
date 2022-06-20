package hossted

import "fmt"

func Dev() error {
	path, err := getSoftwarePath()
	if err != nil {
		return err
	}
	fmt.Println(path)

	return nil
}
