package anathema

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Packages []string `yaml:"packages"`
	Symbols  []string `yaml:"symbols"`
}

var configPath string

func getConfig(c *Configuration) *Configuration {
	if c != nil {
		return c
	} else if configPath == "" {
		log.Fatal("Need to specify a configuration.")
	}

	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Unable to read the specified configuration at %q: %v", configPath, err)
	}

	c = &Configuration{}
	if err = yaml.Unmarshal(raw, c); err != nil {
		log.Fatalf("The configuration in %q could not be parsed: %v", configPath, err)
	}
	return c
}

func (c *Configuration) validate() (*configuration, error) {
	config := &configuration{
		packages: map[string]string{},
		symbols:  map[string]string{},
	}

	for _, l := range c.Packages {
		packages, err := unfoldLine(l)
		if err != nil {
			return nil, err
		}

		for _, p := range packages {
			switch sp := strings.Split(p, "="); len(sp) {
			case 1:
				config.packages[p] = ""
			case 2:
				config.packages[sp[0]] = sp[1]
			default:
				return nil, fmt.Errorf("invalid package specification %q should be of form '<import-path>[=<import-path>]'", p)
			}
		}
	}

	symbolSpecIsValid := func(spec string) bool {
		return strings.Count(spec[strings.LastIndex(spec, ".")+1:], "/") == 0 &&
			strings.Count(spec[strings.LastIndex(spec, "/")+1:], ".") == 1
	}

	for _, l := range c.Symbols {
		symbols, err := unfoldLine(l)
		if err != nil {
			return nil, err
		}

		for _, s := range symbols {
			switch sp := strings.Split(s, "="); len(sp) {
			case 1:
				if !symbolSpecIsValid(s) {
					return nil, fmt.Errorf("invalid symbol specification %q should be of form '<import-path>.<symbol>'", s)
				}
				config.symbols[s] = ""
			case 2:
				if !symbolSpecIsValid(sp[0]) {
					return nil, fmt.Errorf("invalid symbol specification %q should be of form '<import-path>.<symbol>'", sp[0])
				} else if !symbolSpecIsValid(sp[1]) {
					return nil, fmt.Errorf("invalid symbol specification %q should be of form '<import-path>.<symbol>'", sp[1])
				}
				config.symbols[sp[0]] = sp[1]
			default:
				return nil, fmt.Errorf("invalid symbol specification %q should be of form '<import-path>.<symbol>[=<import-path>.<symbol>]'", s)
			}
		}
	}
	return config, nil
}

type configuration struct {
	packages map[string]string
	symbols  map[string]string
}

func unfoldLine(line string) ([]string, error) {
	switch specs := strings.Split(line, "="); len(specs) {
	case 1:
		return unfoldSpec(specs[0])
	case 2:
		sources, err := unfoldSpec(specs[0])
		if err != nil {
			return nil, err
		}
		targets, err := unfoldSpec(specs[1])
		if err != nil {
			return nil, err
		}
		if len(sources) != len(targets) {
			return nil, fmt.Errorf("%q has a mismatched number of source and target specifications", line)
		}

		var results []string
		for idx := 0; idx < len(sources); idx++ {
			results = append(results, fmt.Sprintf("%s=%s", sources[idx], targets[idx]))
		}
		return results, nil
	default:
		return nil, fmt.Errorf("%q has an unexpected number of '=' characters", line)
	}
}

var specRE = regexp.MustCompile(`^([^{}]*)(?:{([^{}]+)})?([^{}]*)$`)

func unfoldSpec(spec string) ([]string, error) {
	m := specRE.FindStringSubmatch(spec)
	if len(m) == 0 {
		return nil, fmt.Errorf("%q is an invalid specification", spec)
	}

	var specs []string
	for _, s := range strings.Split(m[2], ",") {
		specs = append(specs, m[1]+s+m[3])
	}
	return specs, nil
}
