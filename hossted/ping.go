package hossted

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"
	"time"
)

// hossted ping - send docker ,sbom and security infor to hossted API
func Ping(env string) error {

	config, _ := GetConfig() // Ignore error

	//Get uuid, env. Env default to be dev, if env varible
	env = GetHosstedEnv(env)
	uuid, err := GetHosstedUUID(config.UUIDPath)
	if err != nil {
		return err
	}
	//Get dockers
	dockersJson, err := GetDockersInfo()
	if err != nil {
		fmt.Println(err)
		return err
	}
	//Send request
	response, err := PingRequest(dockersJson, uuid, env)
	if err != nil {
		return err
	}
	//fmt.Printf("response: %v\n", response)

	// TODO: Check response status
	message := strings.TrimSpace(response.Message)
	fmt.Printf("response message: %s\n", message)

	return nil
}

// PingRequest sends dockers request
func PingRequest(dockers, uuid, env string) (pingResponse, error) {

	var response pingResponse

	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	fileName := uuid + "-" + fmt.Sprint(umillisec)

	// write dockers to body with multipart (sends as file)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fileName, "file")
	part.Write([]byte(dockers))
	writer.Close()

	// Construct param map for input params
	params := make(map[string]string)
	params["uuid"] = uuid
	params["fileName"] = fileName

	req := HosstedRequest{
		// Endpoint env needs to replace in runtime for url parse to work. Otherwise runtime error.
		//EndPoint:     "https://api.__ENV__hossted.com/v1/instances/dockers",
		EndPoint:     "https://api.hossted.com/v1/instances/dockers", //"https://api.stage.hossted.com/v1/instances/dockers", //"http://localhost:3004/v1/dockers", //, // ,//
		Environment:  env,
		Params:       params,
		BearToken:    "Basic y5TXKDY4kTKbFcFtz9aD1pa2irmzhoziKPnEBcA8",
		SessionToken: "",
		TypeRequest:  "PATCH",
		ContentType:  writer.FormDataContentType(),
		Body:         body,
	}

	fmt.Println("Dockers creation Please wait a second...")
	resp, err := req.SendRequest()
	if err != nil {
		fmt.Println("err", err)
		return response, err
	}

	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return response, fmt.Errorf("Failed to parse JSON. %w", err)
	}

	return response, nil
}
