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
	Packages Packages `yaml:"packages"`
	Symbols  Symbols  `yaml:"symbols"`
}

type Packages struct {
	Whitelist bool          `yaml:"whitelist"`
	Rules     []PackageRule `yaml:"rules"`
}

type PackageRule struct {
	Path        string `yaml:"path"`
	Replacement string `yaml:"replacement"`
}

type Symbols struct {
	Whitelist bool         `yaml:"whitelist"`
	Rules     []SymbolRule `yaml:"rules"`
}

type SymbolRule struct {
	Package            string `yaml:"package"`
	Name               string `yaml:"names"`
	ReplacementPackage string `yaml:"replacement_package"`
	ReplacementName    string `yaml:"replacement_name"`
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

type configuration struct {
	packages          map[string]string
	whitelistPackages bool

	symbols          map[string]string
	whitelistSymbols bool
}

func (c *Configuration) validate() (*configuration, error) {
	var err error
	config := &configuration{
		whitelistPackages: c.Packages.Whitelist,
		whitelistSymbols:  c.Symbols.Whitelist,
	}

	config.packages, err = expandPackageRules(c.Packages.Rules)
	if err != nil {
		return nil, err
	}

	config.symbols, err = expandSymbolRules(c.Symbols.Rules)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func expandPackageRules(rules []PackageRule) (map[string]string, error) {
	expanded := map[string]string{}
	for _, r := range rules {
		packages, err := expandLine(r.Path)
		if err != nil {
			return nil, fmt.Errorf("package rule %+v contained an error in its path: %s", r, err)
		}
		fmt.Println(packages)

		var replacements []string
		if r.Replacement != "" {
			replacements, err = expandLine(r.Replacement)
			if err != nil {
				return nil, fmt.Errorf("package rule %+v contained an error in its replacement: %s", r, err)
			} else if len(replacements) != len(packages) {
				return nil, fmt.Errorf("package rule %+v has a mismatched number of replacement specifications", r)
			}
		}

		for idx := 0; idx < len(packages); idx++ {
			var target string
			if len(replacements) > 0 {
				target = replacements[idx]
			}
			expanded[packages[idx]] = target
		}
	}
	return expanded, nil
}

func expandSymbolRules(rules []SymbolRule) (map[string]string, error) {
	expanded := map[string]string{}
	for _, r := range rules {
		if r.Package == "" {
			return nil, fmt.Errorf("symbol rule %+v is missing a package path", r)
		} else if strings.Count(r.Package, ",") > 0 {
			return nil, fmt.Errorf("symbol rule %+v specifies multiple packages which is not supported", r)
		} else if strings.Count(r.ReplacementPackage, ",") > 0 {
			return nil, fmt.Errorf("symbol rule %+v specifies multiple packages as replacement which is not supported", r)
		}

		symbols, err := expandLine(r.Name)
		if err != nil {
			return nil, fmt.Errorf("symbol rule %+v contained an error in its name: %s", r, err)
		}

		var replacements []string
		if r.ReplacementName != "" {
			replacements, err = expandLine(r.ReplacementName)
			if err != nil {
				return nil, fmt.Errorf("symbol rule %+v contained an error in its replacement: %s", r, err)
			} else if len(replacements) != len(symbols) {
				return nil, fmt.Errorf("symbol rule %+v has a mismatched number of replacement specifications", r)
			}
		}

		var targetPkg string
		if r.ReplacementPackage != "" {
			targetPkg = r.ReplacementPackage
		} else if len(replacements) > 0 {
			targetPkg = r.Package
		}

		for idx := 0; idx < len(symbols); idx++ {
			var target string
			if targetPkg != "" {
				if len(replacements) > 0 {
					target = targetPkg + "." + replacements[idx]
				} else {
					target = targetPkg + "." + symbols[idx]
				}
			}
			expanded[r.Package+"."+symbols[idx]] = target
		}
	}
	return expanded, nil
}

func expandLine(line string) ([]string, error) {
	var curr string
	var specs []string
	for _, sp := range strings.Split(line, ",") {
		if sp == "" {
			return nil, fmt.Errorf("target specification %q contained an empty element", line)
		}

		if (curr != "" || strings.Count(sp, "{") > 0) && strings.Count(sp, "}") == 0 {
			curr += sp + ","
			continue
		}
		curr += sp

		if strings.Count(curr, "{") > 0 {
			expanded, err := expandSpec(curr)
			if err != nil {
				return nil, fmt.Errorf("target specification %q contained an error: %s", line, err)
			}
			specs = append(specs, expanded...)
		} else {
			specs = append(specs, curr)
		}
		curr = ""
	}
	if curr != "" {
		return nil, fmt.Errorf("target specification %q contains an unclosed brace", line)
	}
	return specs, nil
}

var expansionRE = regexp.MustCompile(`^([^{},]*)(?:{([^{}]+)})?([^{},]*)$`)

func expandSpec(spec string) ([]string, error) {
	m := expansionRE.FindStringSubmatch(spec)
	if len(m) == 0 {
		return nil, fmt.Errorf("%q needs to be of form 'foo', 'foo/bar', 'foo.{bar,boo}', etc", spec)
	}

	var specs []string
	for _, s := range strings.Split(m[2], ",") {
		expanded := m[1] + s + m[3]
		if expanded == "" {
			return nil, fmt.Errorf("%q contains an empty element", expanded)
		}
		specs = append(specs, expanded)
	}
	return specs, nil
}
