package hossted

// Config is a struct to parse config.yaml file
type Config struct {
	Email        string `yaml:"email"`
	UserToken    string `yaml:"userToken"`
	SessionToken string `yaml:"sessionToken"`
	EndPoint     string `yaml:"endPoint"`
	UUIDPath     string `yaml:"uuidPath"`
	HostUUID     string `yaml:"hostUuid"`
    Update       bool   `yaml:"update"`
	Applications []ConfigApplication `yaml:"applications"`
}

// ConfigApplication is the applications installled in the vm.
// Currently it will look up the values from /opt/hossted/run/software.txt (Preferred) or /opt/linnovate/run/software.txt
// Supporting multiple values for future enhancement.
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
	TypeRequest	 string            // Request type, GET, POST, PUT, DELETE
}

// RegisterResponse is the return response from the register api
type RegisterResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"msg"`
	JWT        string `json:"jwt"`
	URL        string `json:"url"`
}
// pingResponse is the return response from the register api
type pingResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"msg"`
}
// AvailableCommand is the predefined app/command mapping.
// Maintained with the command.go file
type AvailableCommand struct {
	Apps []App `yaml:"apps"`
}

// App is a struct to save the available commands for a particular application
// The length of Commands and Values is expected to be the same, to give additional
// information or examples to users, on to call the command
type App struct {
	App          string   `yaml:"app"`
	CommandGroup string   `yaml:"group"`
	Commands     []string `yaml:"commands"`
	Values       []string `yaml:"values"`
}

// Command is the individual command, with the app and example values information
type Command struct {
	App          string
	CommandGroup string
	Command      string
	Value        string
}

// AvailableCommandMap saves the map for available commands
// e.g. map["app.command"] -> [{app command value} {prometheus url example.com} ]
type AvailableCommandMap map[string]Command

type YamlSetting struct {
	Pattern  string // regex
	NewValue string // value of the new input, should be matching the number of match groups in regex
}
