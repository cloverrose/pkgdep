package pkgdep

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/cloverrose/pkgdep/pkg/log"
)

func init() {
	RegisterPlugin()
}

func RegisterPlugin() {
	// https://golangci-lint.run/plugins/module-plugins/
	register.Plugin("pkgdep", newPlugin)
}

func newPlugin(conf any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[settings](conf)
	if err != nil {
		return nil, err
	}

	return &plugin{settings: &s}, nil
}

type settings struct {
	Config string
	Log    log.Config
}

type plugin struct {
	settings *settings
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	if p.settings.Config != "" {
		configFile = p.settings.Config
	}
	if p.settings.Log.Level != "" {
		logConfig.Level = p.settings.Log.Level
	}
	if p.settings.Log.File != "" {
		logConfig.File = p.settings.Log.File
	}
	if p.settings.Log.Format != "" {
		logConfig.Format = p.settings.Log.Format
	}
	return []*analysis.Analyzer{
		Analyzer,
	}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

var _ register.LinterPlugin = &plugin{}
