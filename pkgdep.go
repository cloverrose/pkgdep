package pkgdep

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/go/analysis"
	"gopkg.in/yaml.v3"
)

const doc = "pkgdep validates if package dependency follows rule"

// Analyzer checks if package dependency follows rule.
var Analyzer = &analysis.Analyzer{
	Name:     "pkgdep",
	Doc:      doc,
	Run:      run,
	Requires: []*analysis.Analyzer{},
	Flags:    *flag.NewFlagSet("pkgdep", flag.ExitOnError),
}

// configFile is file path to pkgdep config file.
// Allowed file extension is [.yaml, .yml]
// e.g. ./.pkgdep.yaml
var configFile string

func init() {
	Analyzer.Flags.StringVar(&configFile, "config", "", "config file path.")
}

type Config struct {
	TargetPackagePrefixList []string            `yaml:"targetPackagePrefixList"`
	IsExcludeTests          bool                `yaml:"isExcludeTests"`
	EnableRegexp            bool                `yaml:"enableRegexp"`
	Dependencies            map[string][]string `yaml:"dependencies"`
}

func (c *Config) isTargetPackage(pkg string) bool {
	for _, prefix := range c.TargetPackagePrefixList {
		if strings.HasPrefix(pkg, prefix) {
			return true
		}
	}
	return false
}

func (c *Config) isAllowedDependency(from, to string) bool {
	for fromPattern, toTemplateStrings := range c.Dependencies {
		data, err := matchAndExtract(fromPattern, from)
		if err != nil {
			continue
		}
		for _, toTemplateString := range toTemplateStrings {
			toPattern, err := buildPattern(toTemplateString, data)
			if err != nil {
				continue
			}
			re, err := regexp.Compile(toPattern)
			if err != nil {
				continue
			}
			if re.MatchString(to) {
				return true
			}
		}
	}
	return false
}

func matchAndExtract(pattern, text string) (map[string]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	match := re.FindStringSubmatch(text)
	if match == nil {
		return nil, errors.New("no match")
	}
	names := re.SubexpNames()
	data := make(map[string]string)
	for i, name := range names {
		if i > 0 && i < len(match) {
			data[name] = match[i]
		}
	}
	return data, nil
}

func buildPattern(templateString string, data map[string]string) (string, error) {
	tmpl, err := template.New("example").Parse(templateString)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}

	return result.String(), nil
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)

	// To distinguish the case that config file explicitly set enableRegexp=false,
	// we set it to true by default.
	cfg.EnableRegexp = true

	if strings.HasSuffix(configFile, ".json") {
		return nil, errors.New("JSON configuration file is no longer supported. See breaking_changes.md")
	} else if strings.HasSuffix(configFile, ".yaml") || strings.HasSuffix(configFile, ".yml") {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unsupported file suffix. supported file suffixes are [yaml, yml]")
	}

	if !cfg.EnableRegexp {
		return nil, errors.New("enableRegexp option is obsolete. See breaking_changes.md")
	}

	return cfg, nil
}

func run(pass *analysis.Pass) (any, error) {
	if configFile == "" {
		return nil, nil
	}
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	fromPackage := pass.Pkg.Path()
	if !cfg.isTargetPackage(fromPackage) {
		return nil, nil
	}
	for _, f := range pass.Files {
		if cfg.IsExcludeTests {
			pos := pass.Fset.Position(f.Pos())
			if strings.HasSuffix(pos.Filename, "_test.go") {
				continue
			}
		}

		for _, ip := range f.Imports {
			toPackage, err := strconv.Unquote(ip.Path.Value)
			if err != nil {
				return nil, err
			}
			if !cfg.isTargetPackage(toPackage) {
				continue
			}
			if !cfg.isAllowedDependency(fromPackage, toPackage) {
				pass.Reportf(ip.Pos(), "Dependency from %s to %s is not allowed", fromPackage, toPackage)
			}
		}
	}
	return nil, nil
}
