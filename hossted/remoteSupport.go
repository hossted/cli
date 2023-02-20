package hossted

import (
	"fmt"
	"strings"
)

var AUTHORIZED_KEY_PATH = "/root/.ssh/authorized_keys"

// SetRemoteAccess set the remote access by comment/uncomment the hossted public key in the ~/.ssh/authorized_keys file
func SetRemoteAccess(env string, flag bool) error {
	filepath := AUTHORIZED_KEY_PATH
	_ = filepath

	lines, err := changeRemoteAccess(flag)
	if err != nil {
		return err
	}

	// Write file content back to file
	content := strings.Join(lines, "\n")
	err = writeProtected(filepath, []byte(content))
	if err != nil {
		return err
	}
	fmt.Printf("Updated authorized key - %s\n", filepath)

	// send activity log about the command
	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("Something is wrong with get config.\n%w", err)
	}
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}
	fullCommand := "hossted set remote-support " + fmt.Sprint(flag)
	options := `{"remote-support":` + fmt.Sprint(flag) + `}`
	typeActivity := "set_remote"
	sendActivityLog(env, uuid, fullCommand, options, typeActivity)
	return nil
}

// changeRemoteAccess comment/uncomment the hossted key line
// return list of string as each line of the original file
func changeRemoteAccess(flag bool) ([]string, error) {
	filepath := AUTHORIZED_KEY_PATH

	// Public key of hossted
	PUBLICKEY := "AAAAB3NzaC1yc2EAAAADAQABAAABAQCXfKWaimQKigm7A8mxqoEr2e2OOBxpBMvYu8BmMuP2+GEU3FP2CHnVUAPlQ7ByniW+Qg7xlDLQa9kb+5X3r2bwnN7FkEzD3fF/qbfJHjdFlPGN+tdkkoSZzfzsa/z3F4H7wI5pz48cZryph8/dKlrG92GF05womLrFUQyHPbuyMYqGoXTMIaXYaWQm84BndMu+cvqas25YOLyYQXT50BrcWNF24+aqo6H9upmRMN+R1KN0pHxV+5jPDOyZj8PiAAVWCAM574mSm/R8jKZ0qdyMcUL0uThIBtDbMXhgfJn3NdD/3x/DxNnXy84WIHPpY5OPzU3miuVMkNpD8QZRLTFX"
	hosstedKey := fmt.Sprintf("ssh-rsa %s linnovate", PUBLICKEY)

	// Read from /root/.ssh/authorized_keys, and split to lines
	b, err := readProtected(filepath)
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")

	// append the original content with the comment/uncomment hossted key
	var content []string
	for _, l := range lines {
		if strings.Contains(l, PUBLICKEY) {
			if flag {
				content = append(content, hosstedKey)
			} else {
				hosstedKey = fmt.Sprintf("# %s", hosstedKey)
				content = append(content, hosstedKey)
			}
		} else {
			// only append non empty line
			if strings.TrimSpace(l) != "" {
				content = append(content, l)
			}
		}

	}
	return content, nil
}
