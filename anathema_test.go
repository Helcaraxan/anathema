package anathema

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestBlacklist(t *testing.T) {
	t.Parallel()

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
	analysistest.Run(t, analysistest.TestData(), Analysis(testConfig), "pkg/blacklist")
}

func TestWhitelist(t *testing.T) {
	t.Parallel()

	testConfig := &Configuration{
		Packages: Packages{
			Whitelist: true,
			Rules: []PackageRule{
				{Path: "pkg/internal/helpers"},
				{Path: "pkg/internal/new"},
			},
		},
		Symbols: Symbols{
			Whitelist: true,
			Rules: []SymbolRule{
				{Package: "pkg/internal/new", Name: "Background"},
			},
		},
	}
	analysistest.Run(t, analysistest.TestData(), Analysis(testConfig), "pkg/whitelist")
}
