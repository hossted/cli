package hossted

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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
		return ca, fmt.Errorf("No Config found for app - %s", in)
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

//////////////////////////////////////////
// Docker compose related struct method
//////////////////////////////////////////

func (d *DockerStruct) Unmarshal(data []byte) error {

	// Set up and define variables
	m := make(map[int][]DockerLine)   // mapping of leading space, with a list of lines
	patternA := "### HOSSTED APP"     // Predefined pattern A
	patternB := "### HOSSTED WRAPPER" // Predefined pattern B
	numA := 0                         // Line num of pattern A
	numB := 0                         // Line num of pattern B

	_ = m
	_ = patternA
	_ = patternB
	_ = numA
	_ = numB

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		_ = i

		trimmedLine := strings.TrimSpace(line)
		ls := countLeadingSpaces(line) // leading space

		s := DockerLine{
			LineNum: i,
			Line:    line,
		}

		_ = trimmedLine
		_ = ls
		_ = s

	}

	return nil
}

// countLeadingSpaces counts the leading space in the sentence.
// Now can not handle lines with tabs
func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
