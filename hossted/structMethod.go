package hossted

import "fmt"

// GetAppConfig finds the specific Application Config with the application name
// e.g. Find the config for prometheus -> ConfigApplication{prometheus /opt/prometheus}
func (c *Config) GetAppConfig(in string) (ConfigApplication, error) {
	var ca ConfigApplication
	for _, app := range c.Applications {
		if app.AppName == in {
			ca = app
			break
		}
	}
	// Check if any matched
	if ca.AppName == "" {
		return ca, fmt.Errorf("No Config found for app - %s", in)
	}
	return ca, nil
}
