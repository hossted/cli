package hossted

import (
	"fmt"
	"os"
)

// ListPS
func ListPS() error {

	// Test
	_, err := GetConfig()
	if err != nil {
		fmt.Println("Please call the command `hossted register` first.\n")
		os.Exit(0)
	}

	return nil
}
