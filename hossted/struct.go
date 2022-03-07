package hossted

// Config is a struct to parse config.yaml file
type Config struct {
	Email        string `yaml:"email"`
	Organization string `yaml:"organization"`
	UserToken    string `yaml:"userToken"`
	SessionToken string `yaml:"sessionToken"`
	EndPoint     string `yaml:"endPoint"`
	UUIDPath     string `yaml:"uuidPath"`
}

// HosstedRequest is a struct to construct neccessary information to send the request to hossted backend
type HosstedRequest struct {
	EndPoint     string            // Request end point
	Environment  string            // environment, dev or prod
	Params       map[string]string // kv pairs for param
	BearToken    string            // Authorization token
	SessionToken string            // Session token. JWT
}

// RegisterResponse
type RegisterResponse struct {
	StatusCode int    `yaml:"status"`
	Message    string `yaml:"msg"`
	jwt        string `yaml:"jwt"`
	url        string `yaml:"url"`
}
