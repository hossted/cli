package hossted

import "fmt"

// SetAuth sets the authorization of the application
func SetAuth(app string, flag bool) error {
	fmt.Printf("app: %+v\n", app)
	fmt.Printf("flag: %+v\n", flag)
	return nil
}
