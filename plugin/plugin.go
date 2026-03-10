package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	projectanalyzer "github.com/rTexty/logsLinter/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

const pluginName = "logslinter"

func init() {
	register.Plugin(pluginName, New)
}

type Linter struct {
	config projectanalyzer.Config
}

func New(settings any) (register.LinterPlugin, error) {
	decodedSettings, err := register.DecodeSettings[projectanalyzer.Settings](settings)
	if err != nil {
		return nil, err
	}

	return &Linter{config: decodedSettings.ToConfig()}, nil
}

func (l *Linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{projectanalyzer.NewAnalyzer(l.config)}, nil
}

func (l *Linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
