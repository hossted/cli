package hossted

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// SendRequest sends a request to hossted API server with parameters
// TODO: This is insecure; use only in dev environments.
// TODO: Add timeout context
// TODO: Check all params is not null
// TODO: Check response status
func (h *HosstedRequest) SendRequest() (string, error) {

	// Set http client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Parse url params. "https://app.dev.hossted.com/api/register?uuid=$UUID&email=$EMAIL&organization=$ORGANIZATION"
	raw := h.EndPoint
	u, _ := url.Parse(raw)
	q, _ := url.ParseQuery(u.RawQuery)

	for k, v := range h.Params {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	endpoint := u.String()
	endpoint = updateEndpointEnv(endpoint, h.Environment)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", h.BearToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("Status code - %d. Message - %s\n", resp.StatusCode, string(body))

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP Status is not 200. %s", string(body))
	}

	return string(body), nil
}
