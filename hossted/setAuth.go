package hossted

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// SetAuth sets the authorization of the application
func SetAuth(env, app string, flag bool) error {

	if !HasContainerRunning() {
		fmt.Println("The application still in configuration")
		os.Exit(0)
	}

	if flag == true {
		return errors.New("\n  Not Implemented for the command set auth true.\n")
	}

	config, err := GetConfig()
	if err != nil {
		return fmt.Errorf("Something is wrong with get config.\n%w", err)
	}

	// Check command
	err = CheckCommands(app, "auth")
	if err != nil {
		return err
	}

	// Get appPath
	appConfig, err := config.GetAppConfig(app)
	if err != nil {
		return err
	}
	name := appConfig.AppName   // app name
	appDir := appConfig.AppPath // app directory. e.g. /opt/gitbucket
	if err != nil {
		return err
	}

	// sed commands
	filepath := fmt.Sprintf("/opt/%s/docker-compose.yml", app) // docker compsoe file path
	commands := []string{
		fmt.Sprintf("sudo sed -i '/tauth.basicauth.usersfile/d' '%s'", filepath),
		fmt.Sprintf("sudo sed -i '/.middlewares=tauth,hauth/d' '%s'", filepath),
		fmt.Sprintf("sudo sed -i -e 's/tauth,//g' '%s'", filepath),
		fmt.Sprintf("sudo sed -i '/.middlewares=tauth/d' '%s'", filepath),
		"sudo sed -i '/installation you may remove/d' '/etc/motd'",
	}

	fmt.Println("Change settings...")
	err, _, stderr := Shell(appDir, commands)
	if err != nil {
		return err
	}
	if strings.TrimSpace(stderr) != "" {
		fmt.Println(stderr)
	}

	// Remove letsencrypt
	htpassPath := fmt.Sprintf("/opt/%s/hossted/letsencrypt/.htpass", name)
	rmCommands := []string{
		fmt.Sprintf("sudo rm '%s'", htpassPath),
	}
	fmt.Println(fmt.Sprintf("Removed %s", htpassPath))
	err, _, stderr = Shell(appDir, rmCommands)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s does not exists.\nProbably being removed earlier.\n", htpassPath), err.Error())
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

	//send activity log about the command
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}
	fullCommand := "hossted set auth " + fmt.Sprint(flag)
	options := `{"auth":` + fmt.Sprint(flag) + `}`
	typeActivity := "set_auth"
	sendActivityLog(env, uuid, fullCommand, options, typeActivity)

	return nil
}
