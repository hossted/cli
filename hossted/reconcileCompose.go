package hossted

import (
	"fmt"

	"github.com/hossted/cli/hossted/service/compose"
)

func ReconcileCompose() error {
	// get OsInfo
	osInfo, err := compose.GetClusterInfo()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = compose.ReconcileCompose(osInfo, "false")
	if err != nil {
		return err
	}

	return nil
}
