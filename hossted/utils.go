package hossted

import (
	"io/ioutil"
	"path"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// GetConfigPath gets the pre-defined config path. ~/.hossted/config.yaml
func GetConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	configPath := path.Join(home, ".hossted", "config.yaml")
	return configPath, nil
}

// GetConfigPath gets the config object
func GetConfig() (Config, error) {
	var config Config
	cfgPath, err := GetConfigPath()
	if err != nil {
		return config, err
	}

	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
