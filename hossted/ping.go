package hossted

import (
	"encoding/json"
	"os"
	"strings"
	"fmt"
)



// RegisterUsers updates email, organization in the yaml file, if successful.
// Also it will get the uuid of the machine, and environment from $HOSSTED_ENV.
// It will then send an request to the API server, to get the JWT sessin token.
// TODO: Use the original values as default.
func Ping(env string) error {
	var docker string
	
	fmt.Println(env)
	config, _ := GetConfig() // Ignore error


	// Get uuid, env. Env default to be dev, if env varible
	env = GetHosstedEnv(env)
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}

	// Collect docker info
	//ocker ps --format '{"docker_id":"{{ .ID }}", "image_id": "{{ .Image }}", "name":"{{ .Names }}","size":"{{ .Size}}", "tag":"latest"}'
	// dockerJson = GetDockerInfo()
	//cmd := exec.Command("sudo", "docker-compose", "ps", "--format" ,'{"docker_id":"{{ .ID }}", "image_id": "{{ .Image }}", "name":"{{ .Names }}","size":"{{ .Size}}", "tag":"latest"}')
	//cmd.Dir = app.AppPath

	/*out, err := cmd.Output()
	if err != nil {
		return err
	}
	*/
	//fmt.Println(out)
	// Send request
	docker = "{'docker_id':'c5b75b06b05a','image_id': 'hossted-api_consumer', 'name':'api_consumer','size':'2.02MB (virtual 1.13GB)', 'tag':'latest'}"
	response, err := PingRequest( docker, uuid, env)
	if err != nil {
		return err
	}

	// TODO: Check response status
	jwt := strings.TrimSpace(response.JWT)
	endpoint := strings.TrimSpace(response.URL)

	// Assign back to config object
	config.SessionToken = jwt
	config.EndPoint = endpoint

	// Write back to file
	err = WriteConfigWrapper(config)
	if err != nil {
		return fmt.Errorf("Can not write to config file. Please check. %w", err)
	}
	dashboardUrl := endpoint + "?token=" + config.SessionToken
	fmt.Println(fmt.Sprintf("Please visit the dashboard link - %s\n", dashboardUrl))
	return nil
}

// registerRequest sends register request based on the input, env/email/organization, etc..
// TODO: Set BearToken to env variable
func PingRequest(docker , uuid, env string) (RegisterResponse, error) {

	var response RegisterResponse

	// Construct param map for input params
	params := make(map[string]string)
	params["uuid"] = uuid

	req := HosstedRequest{
		// Endpoint env needs to replace in runtime for url parse to work. Otherwise runtime error.
		// Old EndPoint:     "https://app.__ENV__hossted.com/api/register",
		EndPoint:     "https://api.__ENV__hossted.com/v1/instances/dockers",
		Environment:  env,
		Params:       params,
		BearToken:    "Basic FrTc3TlygOaFDQOGmteaQ7LRwKOx8XNIGfmLa5NA",
		SessionToken: "",
	}

	fmt.Println("Docker creation Please wait a second...")
	resp, err := req.SendRequest()
	if err != nil {
		return response, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return response, fmt.Errorf("Failed to parse JSON. %w", err)
	}

	// Check if the sessionToken is null
	if strings.TrimSpace(response.JWT) == "" {
		if strings.Contains(response.Message, "already been registered") {

			fmt.Printf("\nInstance is already registered. Please contact support@hossted.com if you'd like to add a user.\n\n")
			os.Exit(0)

		} else {
			fmt.Printf("Empty Session Token. Some kind of error occoured - %s.\n", resp)
			os.Exit(0)
		}
	}

	return response, nil
}

