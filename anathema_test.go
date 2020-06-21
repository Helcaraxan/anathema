package anathema

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalysis(t *testing.T) {
	testConfig := &Configuration{
		Packages: []string{
			"pkg/internal/forbidden",
			"pkg/internal/old=pkg/internal/new",
		},
		Symbols: []string{
			"external.External",
			"pkg/internal/helpers.Constant",
			"pkg/internal/helpers.FuncFactory",
			"pkg/internal/helpers.InterfaceType",
			"pkg/internal/helpers.StructFactory",
			"pkg/internal/helpers.StructType",
			"pkg/internal/helpers.Variable",
			"pkg/internal/old.Context=pkg/internal/new.Context",
		},
	}
	analysistest.Run(t, analysistest.TestData(), Analysis(testConfig), "pkg")
}
