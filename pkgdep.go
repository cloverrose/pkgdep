package pkgdep

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"gopkg.in/yaml.v3"

	"github.com/cloverrose/pkgdep/pkg/checker"
	"github.com/cloverrose/pkgdep/pkg/inspector"
	"github.com/cloverrose/pkgdep/pkg/log"
	"github.com/cloverrose/pkgdep/pkg/orderedmap"
)

const doc = "pkgdep validates if package dependency follows rule"

// Analyzer checks if package dependency follows rule.
var Analyzer = &analysis.Analyzer{
	Name:     "pkgdep",
	Doc:      doc,
	Run:      setupAndRun,
	Requires: []*analysis.Analyzer{},
	Flags:    *flag.NewFlagSet("pkgdep", flag.ExitOnError),
}

// options
var (
	// configFile is file path to pkgdep config file.
	// Allowed file extension is [.yaml, .yml]
	// e.g. ./.pkgdep.yaml
	configFile string

	// log related configuration.
	logConfig = log.Config{
		Level:  "INFO",
		File:   "",
		Format: "json",
	}

	inspectorConfig = inspector.Config{
		File: "",
	}
)

var inspectorInstance *inspector.Inspector

func init() {
	Analyzer.Flags.StringVar(&configFile, "config", "", "config file path.")
	Analyzer.Flags.StringVar(&logConfig.Level, "log.level", logConfig.Level, "logging level. debug, info, warn, error")
	Analyzer.Flags.StringVar(&logConfig.File, "log.file", logConfig.File, "log file path.")
	Analyzer.Flags.StringVar(&logConfig.Format, "log.format", logConfig.Format, "logging format. json or text")
	Analyzer.Flags.StringVar(&inspectorConfig.File, "inspector.file", inspectorConfig.File, "inspector file path")
}

func setupAndRun(pass *analysis.Pass) (any, error) {
	closer, err := log.SetDefault(logConfig)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := closer(); err != nil {
			fmt.Println(err)
		}
	}()
	slog.Debug("Starting pkgdep analyzer...")

	inspectorInstance = inspector.New(inspectorConfig)
	defer func() {
		if err := inspectorInstance.Save(); err != nil {
			slog.Error("fail inspectorInstance.Save", slog.Any("error", err))
		}
	}()

	return run(pass)
}

const (
	modeAllowList = "allow_list"
	modeBlockList = "block_list"
)

type Config struct {
	TargetPackagePrefixList []string              `yaml:"targetPackagePrefixList"`
	IsExcludeTests          bool                  `yaml:"isExcludeTests"`
	EnableRegexp            bool                  `yaml:"enableRegexp"`
	Dependencies            orderedmap.OrderedMap `yaml:"dependencies"`
	GlobalData              map[string]any        `yaml:"globalData"`
	Mode                    string                `yaml:"mode"` // allow_list / block_list (default is allow_list)
}

func (c *Config) isTargetPackage(pkg string) bool {
	for _, prefix := range c.TargetPackagePrefixList {
		if strings.HasPrefix(pkg, prefix) {
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

	if cfg.Mode != modeAllowList && cfg.Mode != modeBlockList {
		// Set default mode.
		cfg.Mode = modeAllowList
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

	checkerInstance := checker.New(cfg.Dependencies, cfg.GlobalData, inspectorInstance)

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

			matched := checkerInstance.CheckDependency(fromPackage, toPackage)
			switch cfg.Mode {
			case modeBlockList:
				if matched {
					pass.Reportf(ip.Pos(), "Dependency from %s to %s is blocked", fromPackage, toPackage)
				}
			default: // modeAllowList
				if !matched {
					pass.Reportf(ip.Pos(), "Dependency from %s to %s is not allowed", fromPackage, toPackage)
				}
			}
		}
	}
	return nil, nil
}
