package hossted

import "fmt"

// SetRemoteAccess set the remote access by comment/uncomment the hossted public key in the ~/.ssh/authorized_keys file
func SetRemoteAccess(b bool) error {
	fmt.Println("Set remote access")

	return nil
}
