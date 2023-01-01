package hossted

import (
	"encoding/json"
	"fmt"
	"strings"
)

// hossted init -  Send instance hossted API
func Init(env, authorization, image string, instance Instance) error {

	env = GetHosstedEnv(env)
	instancejson, _ := json.Marshal(instance)

	//fmt.Println("instancejson:", string(instancejson))
	//Send request
	response, err := InitRequest(env, authorization, image, string(instancejson))
	if err != nil {
		return err
	}

	// TODO: Check response status
	message := strings.TrimSpace(response.Message)
	fmt.Printf("response: %s\n", message)

	return nil
}

// PingRequest sends dockers request
func InitRequest(env, authorization, image, instancejson string) (initResponse, error) {

	var response initResponse

	// Construct param map for input params
	var data map[string]string
	err := json.Unmarshal([]byte(instancejson), &data)
	if err != nil {
		panic(err)
	}

	params := make(map[string]string)
	for k, v := range data {
		params[k] = v
	}
	req := HosstedRequest{
		// Endpoint env needs to replace in runtime for url parse to work. Otherwise runtime error.
		//EndPoint:     "https://api.__ENV__hossted.com/v1/instances/dockers",
		EndPoint:     "https://api.hossted.com/v1/instances/registry", //"https://api.stage.hossted.com/v1/instances/registry", // "https://api.dev.hossted.com/v1/instances/registry", //,
		Environment:  env,
		Params:       params,
		BearToken:    "Basic " + authorization,
		SessionToken: "",
		TypeRequest:  "POST",
	}

	resp, err := req.SendRequest()
	if err != nil {
		//fmt.Println("\033[31m", "Error:", "\033[0m")
		fmt.Printf("\033[0;31mError: \033[0m")
		fmt.Printf("%s\n\033[0m", err)
		return response, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return response, fmt.Errorf("Failed to parse JSON. %w", err)
	}

	return response, nil
}
