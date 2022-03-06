package hossted

// Config is a struct to parse config.yaml file
type Config struct {
	Email        string `yaml:"email"`
	Organization string `yaml:"organization"`
	UserToken    string `yaml:"userToken"`
	SessionToken string `yaml:"sessionToken"`
	EndPoint     string `yaml:"endPoint"`
}

// HosstedRequest is a struct to construct neccessary information to send the request to hossted backend
type HosstedRequest struct {
	EndPoint    string            // Request end point
	Environment string            // environment, dev or prod
	Params      map[string]string // kv pairs for param
	BearToken   string            // Authorization token
	JWT         string            // Session token
}