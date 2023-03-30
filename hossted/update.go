package hossted

import (
	"fmt"
	"os/exec"
)

func Update(env string) error {

	cmd := "sudo curl -L 'https://github.com/hossted/cli/raw/main/bin/linux/hossted' -o /usr/local/bin/hossted"

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Error executing command %s: %v\n", cmd, err)
	}
	fmt.Printf("updated %s\n", out)
	return nil
}
