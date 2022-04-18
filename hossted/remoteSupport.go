package hossted

import (
	"fmt"
	"strings"
)

// SetRemoteAccess set the remote access by comment/uncomment the hossted public key in the ~/.ssh/authorized_keys file
func SetRemoteAccess(flag bool) error {
	lines, err := changeRemoteAccess(false)
	if err != nil {
		return err
	}
	fmt.Println(lines)

	return nil
}

func changeRemoteAccess(flag bool) ([]string, error) {
	filepath := "/root/.ssh/authorized_keys"

	// Public key of hossted
	PUBLICKEY := "AAAAB3NzaC1yc2EAAAADAQABAAABAQCXfKWaimQKigm7A8mxqoEr2e2OOBxpBMvYu8BmMuP2+GEU3FP2CHnVUAPlQ7ByniW+Qg7xlDLQa9kb+5X3r2bwnN7FkEzD3fF/qbfJHjdFlPGN+tdkkoSZzfzsa/z3F4H7wI5pz48cZryph8/dKlrG92GF05womLrFUQyHPbuyMYqGoXTMIaXYaWQm84BndMu+cvqas25YOLyYQXT50BrcWNF24+aqo6H9upmRMN+R1KN0pHxV+5jPDOyZj8PiAAVWCAM574mSm/R8jKZ0qdyMcUL0uThIBtDbMXhgfJn3NdD/3x/DxNnXy84WIHPpY5OPzU3miuVMkNpD8QZRLTFX"
	hosstedKey := fmt.Sprintf("ssh-rsa %s linnovate", PUBLICKEY)

	// Read from /root/.ssh/authorized_keys, and split to lines
	b, err := readProtected(filepath)
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")

	var content []string
	for _, l := range lines {
		if strings.Contains(l, PUBLICKEY) {
			if flag {
				content = append(content, hosstedKey)
			} else {
				hosstedKey = fmt.Sprintf("# %s", hosstedKey)
				content = append(content, hosstedKey)
			}
		}
		content = append(content, l)
	}
	return content, nil
}
