package hossted

// Config is a struct to parse config.yaml file
type Config struct {
	Email        string              `yaml:"email"`
	UserToken    string              `yaml:"userToken"`
	SessionToken string              `yaml:"sessionToken"`
	EndPoint     string              `yaml:"endPoint"`
	UUIDPath     string              `yaml:"uuidPath"`
	Applications []ConfigApplication `yaml:"applications"`
}

type ConfigApplication struct {
	AppName string `yaml:"appName"`
	AppPath string `yaml:"appPath"`
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
	StatusCode int    `json:"status"`
	Message    string `json:"msg"`
	JWT        string `json:"jwt"`
	URL        string `json:"url"`
}

type AvailableCommand struct {
	Commands []struct {
		App      string   `yaml:"app"`
		Commands []string `yaml:"commands"`
	} `yaml:"apps"`
}
