package anathema

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalysis(t *testing.T) {
	testConfig := &Configuration{
		Packages: Packages{
			Rules: []PackageRule{
				{Path: "pkg/internal/forbidden"},
				{Path: "pkg/internal/old", Replacement: "pkg/internal/new"},
			},
		},
		Symbols: Symbols{
			Rules: []SymbolRule{
				{Package: "external", Name: "External"},
				{Package: "pkg/internal/helpers", Name: "Constant"},
				{Package: "pkg/internal/helpers", Name: "FuncFactory"},
				{Package: "pkg/internal/helpers", Name: "InterfaceType"},
				{Package: "pkg/internal/helpers", Name: "StructFactory"},
				{Package: "pkg/internal/helpers", Name: "StructType"},
				{Package: "pkg/internal/helpers", Name: "Variable"},
				{Package: "pkg/internal/old", Name: "Context", ReplacementPackage: "pkg/internal/new", ReplacementName: "Context"},
			},
		},
	}
	analysistest.Run(t, analysistest.TestData(), Analysis(testConfig), "pkg")
}
