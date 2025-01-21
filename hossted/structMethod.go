package hossted

import (
	"errors"
	"fmt"
	"regexp"
)

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
		return ca, fmt.Errorf("no config found for app - %s", in)
	}
	return ca, nil
}

// GetDefaultApp gets the default app from the config file. If there is only an single app
// from the config, return it as default app. otherwise return error
func (c *Config) GetDefaultApp(pwd string) (string, error) {
	var app string

	pattern := `\/opt\/(.*)(\/)?` // e.g. /opt/gitbucket
	re := regexp.MustCompile(pattern)
	matched, err := regexp.MatchString(pattern, pwd)
	if err != nil {
		return "", fmt.Errorf("Something is wrong with the regex to check default app. %w/n", err)
	}
	if matched {
		matches := re.FindStringSubmatch(pwd)
		if len(matches) >= 2 { // [/opt/gitbucket gitbucket ]
			app = matches[1]
		}
	}
	if app != "" {
		return app, nil
	}

	// Check from config applications
	if len(c.Applications) == 1 {
		app = c.Applications[0].AppName
	} else {
		return "", errors.New("No default apps.")
	}

	return app, nil
}
