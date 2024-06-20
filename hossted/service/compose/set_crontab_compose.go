package compose

import (
	"github.com/hossted/cli/hossted/service/common"
)

func SetCrontabCompose() error {
	command := "/usr/local/bin/hossted reconcile-compose 2>&1 | logger -t mycmd"
	minute := "*/1"
	common.CreateCrontab(command, minute)
	return nil
}
