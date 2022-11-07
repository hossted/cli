package hossted

import (
	"encoding/json"
	"os"
	"strings"
	"fmt"
	//"os/exec"
)



// RegisterUsers updates email, organization in the yaml file, if successful.
// Also it will get the uuid of the machine, and environment from $HOSSTED_ENV.
// It will then send an request to the API server, to get the JWT sessin token.
// TODO: Use the original values as default.
func Ping(env string) error {
	
	fmt.Println(env)
	config, _ := GetConfig() // Ignore error


	// Get uuid, env. Env default to be dev, if env varible
	env = GetHosstedEnv(env)
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}

	//Get dockers
	dockersJson,_ := GetDockersInfo()

	// Send request
	response, err := PingRequest(dockersJson, uuid, env)
	if err != nil {
		return err
	}
    fmt.Printf("response: %v\n", response)

	 // TODO: Check response status
	message := strings.TrimSpace(response.Message)
    fmt.Printf("message: %s\n", message)

	return nil
}

// registerRequest sends register request based on the input, env/email/organization, etc..
// TODO: Set BearToken to env variable
func PingRequest(dockers , uuid, env string) (pingResponse, error) {

	var response pingResponse

	// Construct param map for input params
	params := make(map[string]string)
	params["uuid"] = uuid
	params["dockers"] = dockers

	req := HosstedRequest{
		// Endpoint env needs to replace in runtime for url parse to work. Otherwise runtime error.
		//EndPoint:     "https://api.__ENV__hossted.com/v1/instances/dockers",
		EndPoint:     "http://localhost:3004/v1/dockers",
		Environment:  env,
		Params:       params,
		BearToken:    "Basic y5TXKDY4kTKbFcFtz9aD1pa2irmzhoziKPnEBcA8",
		SessionToken: "",
		TypeRequest:"PATCH",
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

	
if response.Message!=""{
	fmt.Printf("\nresponse.Message: %v\n\n", response.Message)
			os.Exit(0)
}
	return response, nil
}

