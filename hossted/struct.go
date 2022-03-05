package hossted

type Config struct {
	Email        string `yaml:"email"`
	Organization string `yaml:"organization"`
	UserToken    string `yaml:"userToken"`
}
