package hossted

import (
	"fmt"
	"strings"
)

// HttpOpen allows the http access to the frontend
func HttpOpen(input string) error {
	var app ConfigApplication

	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("Please call the command `hossted register` first.\n%w", err)
	}

	app, err = appPrompt(config.Applications, input)
	if err != nil {
		return err
	}

	// Check command
	err = CheckCommands(app.AppName, "httpopen")
	if err != nil {
		return err
	}

	// Get appPath
	appConfig, err := config.GetAppConfig(app.AppName)
	if err != nil {
		return err
	}
	name := appConfig.AppName   // app name
	appDir := appConfig.AppPath // app directory. e.g. /opt/gitbucket
	if err != nil {
		return err
	}

	// sed commands
	commands := []string{
		"sudo sed -i '/tauth.basicauth.usersfile/d' '/opt/gitbucket/docker-compose.yml'",
		"sudo sed -i -e 's/tauth,//g' '/opt/gitbucket/docker-compose.yml'",
		"sudo sed -i '/.middlewares=tauth/d' '/opt/gitbucket/docker-compose.yml'",
		"sudo sed -i '/installation you may remove/d' '/etc/motd'",
	}
	err, _, stderr := Shell(appDir, commands)
	if err != nil {
		return err
	}
	if strings.TrimSpace(stderr) != "" {
		fmt.Println(stderr)
	}

	// Stop Traefik and restart service
	err = stopTraefik(appDir)
	if err != nil {
		return err
	}

	err = dockerUp(appDir)
	if err != nil {
		return err
	}

	fmt.Printf("Service Restarted - %s\n", name)

	return nil
}
