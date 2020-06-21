package anathema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	testcases := map[string]struct {
		config   Configuration
		expected *configuration
	}{
		"PackageStandard": {
			config: Configuration{
				Packages: Packages{
					Rules: []PackageRule{{Path: "go/ast"}},
				},
			},
			expected: &configuration{
				packages: map[string]string{"go/ast": ""},
				symbols:  map[string]string{},
			},
		},
		"PackageRefactored": {
			config: Configuration{
				Packages: Packages{
					Rules: []PackageRule{{Path: "fmt,go/{ast,parser,token},io{,/ioutil},regexp"}},
				},
			},
			expected: &configuration{
				packages: map[string]string{
					"fmt":       "",
					"go/ast":    "",
					"go/parser": "",
					"go/token":  "",
					"io":        "",
					"io/ioutil": "",
					"regexp":    "",
				},
				symbols: map[string]string{},
			},
		},
		"PackageReplacements": {
			config: Configuration{
				Packages: Packages{
					Rules: []PackageRule{{
						Path:        "go/{ast,parser,token}",
						Replacement: "alternative/{ast,parser,token}",
					}},
				},
			},
			expected: &configuration{
				packages: map[string]string{
					"go/ast":    "alternative/ast",
					"go/parser": "alternative/parser",
					"go/token":  "alternative/token",
				},
				symbols: map[string]string{},
			},
		},
		"SymbolStandard": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{
						Package: "fmt",
						Name:    "Print",
					}},
				},
			},
			expected: &configuration{
				packages: map[string]string{},
				symbols: map[string]string{
					"fmt.Print": "",
				},
			},
		},
		"SymbolChangedPackage": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{
						Package:            "fmt",
						Name:               "Print",
						ReplacementPackage: "alternative",
					}},
				},
			},
			expected: &configuration{
				packages: map[string]string{},
				symbols: map[string]string{
					"fmt.Print": "alternative.Print",
				},
			},
		},
		"SymbolChangedName": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{
						Package:         "fmt",
						Name:            "Print",
						ReplacementName: "Println",
					}},
				},
			},
			expected: &configuration{
				packages: map[string]string{},
				symbols: map[string]string{
					"fmt.Print": "fmt.Println",
				},
			},
		},
		"SymbolChangedNameAndPackage": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{
						Package:            "fmt",
						Name:               "Print",
						ReplacementPackage: "myfmt",
						ReplacementName:    "Println",
					}},
				},
			},
			expected: &configuration{
				packages: map[string]string{},
				symbols: map[string]string{
					"fmt.Print": "myfmt.Println",
				},
			},
		},
		"SymbolComplex": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{
						Package:            "fmt",
						Name:               "Print,Printf,Println",
						ReplacementPackage: "myfmt",
						ReplacementName:    "Fprint,Fprintf,Fprintln",
					}},
				},
			},
			expected: &configuration{
				packages: map[string]string{},
				symbols: map[string]string{
					"fmt.Print":   "myfmt.Fprint",
					"fmt.Printf":  "myfmt.Fprintf",
					"fmt.Println": "myfmt.Fprintln",
				},
			},
		},
		"PackageMissingPath": {
			config: Configuration{
				Packages: Packages{
					Rules: []PackageRule{{}},
				},
			},
		},
		"PackageMismatchedReplacement": {
			config: Configuration{
				Packages: Packages{
					Rules: []PackageRule{{
						Path:        "strings,bytes",
						Replacement: "mystrings",
					}},
				},
			},
		},
		"PackageInvalidReplacement": {
			config: Configuration{
				Packages: Packages{
					Rules: []PackageRule{{Path: "foo", Replacement: "bar{"}},
				},
			},
		},
		"SymbolMissingPackage": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{Name: "Print"}},
				},
			},
		},
		"SymbolMissingName": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{Package: "fmt"}},
				},
			},
		},
		"SymbolMismatchedReplacement": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{
						Package:         "fmt",
						Name:            "Print",
						ReplacementName: "Printf,Println",
					}},
				},
			},
		},
		"SymbolMultiplePackages": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{Package: "foo,bar", Name: "Print"}},
				},
			},
		},
		"SymbolMultipleReplacementPackages": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{Package: "foo", Name: "Print", ReplacementPackage: "foo,bar"}},
				},
			},
		},
		"SymbolsInvalidReplacement": {
			config: Configuration{
				Symbols: Symbols{
					Rules: []SymbolRule{{Package: "foo", Name: "Print", ReplacementName: "{Bar,}"}},
				},
			},
		},
	}

	for name := range testcases {
		testcase := testcases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := testcase.config.validate()
			if testcase.expected == nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testcase.expected, result)
		})
	}
}
