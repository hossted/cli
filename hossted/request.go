package hossted

import (
	"crypto/tls"
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
	q.Add("uuid", h.Params["uuid"])
	q.Add("email", h.Params["email"])
	q.Add("organization", h.Params["organization"])
	u.RawQuery = q.Encode()
	endpoint := u.RawQuery

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

	return string(body), nil
}
