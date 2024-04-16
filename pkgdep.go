package pkgdep

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const doc = "pkgdep checks if package dependency follows rule"

// Analyzer checks if package dependency follows rule.
var Analyzer = &analysis.Analyzer{
	Name:     "pkgdep",
	Doc:      doc,
	Run:      run,
	Requires: []*analysis.Analyzer{},
	Flags:    *flag.NewFlagSet("pkgdep", flag.ExitOnError),
}

var configFile string

func init() {
	Analyzer.Flags.StringVar(&configFile, "config", "", "config json file path. It should be absolute path.")
}

type Config struct {
	TargetPackagePrefixList []string
	IsExcludeTests          bool
	Dependencies            map[string][]string
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
	allows := c.Dependencies[from]
	for _, allow := range allows {
		if allow == "*" { // Wild card
			return true
		}
		if strings.HasSuffix(allow, "/*") {
			newAllow := strings.ReplaceAll(allow, "/*", "/")
			if strings.HasPrefix(to, newAllow) {
				return true
			}
		}
		if allow == to {
			return true
		}
	}
	return false
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("cannot decode JSON config file %s: %v", configFile, err)
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
