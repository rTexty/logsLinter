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

type Linter struct{}

func New(any) (register.LinterPlugin, error) {
	return &Linter{}, nil
}

func (l *Linter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{projectanalyzer.Analyzer}, nil
}

func (l *Linter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}