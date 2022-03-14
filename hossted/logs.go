package hossted

import "fmt"

// GetAppLogs goes to the app directory, then calls docker-compose logs
func GetAppLogs() error {
	fmt.Println("Get App logs")

	return nil
}
