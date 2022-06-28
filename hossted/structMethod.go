package hossted

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

// Unmarshal blah blah
// TODO: Add test to check error logic
func (d *DockerStruct) Unmarshal(data []byte) error {

	// Set up and define variables
	SPACING := 2                        // hardcoded constant for parsing DockerApp
	m := make(map[int][]DockerLine)     // mapping of leading space, with a list of lines
	patternA := "# HOSSTED APP"         // Predefined pattern A
	patternB := "# HOSSTED WRAPPER"     // Predefined pattern B
	patternC := "# HOSSTED WRAPPER END" // Predefined pattern C
	numA := 0                           // Line num of pattern A
	numB := 0                           // Line num of pattern B
	numC := 0                           // Line num of pattern B
	var version int                     // Parse docker compose version
	var err error

	_ = m
	_ = patternA
	_ = patternB
	_ = numA
	_ = numB
	_ = version

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		_ = i

		trimmedLine := strings.TrimSpace(line)
		ls := countLeadingSpaces(line) // leading space

		// Parse version line
		if strings.HasPrefix(strings.ToLower(line), "version:") {
			// matching the line with pattern version: '2'
			v := strings.ReplaceAll(strings.ReplaceAll(line, " ", ""), "version:", "")
			version, err = strconv.Atoi(v)
			if err != nil {
				version = 0
			}
		}

		s := DockerLine{
			LineNum: i,
			Line:    line,
		}

		// Capture Pattern Lines
		if (numA == 0) && strings.Contains(trimmedLine, patternA) {
			numA = i
		}
		if (numB == 0) && strings.Contains(trimmedLine, patternB) {
			numB = i
		}
		if (numC == 0) && strings.Contains(trimmedLine, patternC) {
			numC = i
		}

		// Build mapping, with keys as no of leading space, values as list of lines
		if val, ok := m[ls]; ok {
			m[ls] = append(val, s)
		} else {
			m[ls] = []DockerLine{s}
		}
	}

	// Return error if any of the patterns are missing from the docker compose file
	// TODO: Need add test for checking error logics
	if (numA == 0) || (numB == 0) || (numC == 0) {
		return errors.New("The specific hossted docker file pattern lines are not available in the docker file.\nPlease check with administrator.\n")
	}

	// The pattern should be in sequence. patternA (First) > patternB (Second) > patternC (Last)
	// TODO: Need add test for checking error logics
	if (numA >= numB) || (numB >= numC) || (numA >= numC) {
		return errors.New("The specific hossted docker file pattern lines are not in specific orders in the docker file.\nPlease check with administrator.\n")
	}

	// Parse apps
	var (
		apps  []DockerApp // Normal apps
		wapps []DockerApp // Wrapped apps

	)

	secondSpacing := m[SPACING] // mapping with 2 leading spaces
	nApps := len(secondSpacing)
	nLine := len(lines)
	_ = nApps
	_ = nLine
	_ = apps
	_ = wapps

	for i := 0; i < nApps; i++ {
		var (
			start int
			end   int
		)

		fmt.Println("------")

		if i < nApps-1 {
			start = secondSpacing[i].LineNum
			end = secondSpacing[i+1].LineNum - 1
			fmt.Printf("If: %d - (%d, %d)\n", i, start, end)

		} else {
			start = secondSpacing[i].LineNum
			end = nLine
			fmt.Printf("Else: %d - (%d, %d)\n", i, start, end)
		}

		appName := strings.TrimSpace(lines[start]) // app name
		content := lines[start:end]                // app content
		app := DockerApp{
			Name:    appName,
			Content: content,
		}
		apps = append(apps, app)
		fmt.Println(PrettyPrint(app))

	}

	// fmt.Println(PrettyPrint(m))
	return nil
}

// countLeadingSpaces counts the leading space in the sentence.
// Now can not handle lines with tabs
func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
