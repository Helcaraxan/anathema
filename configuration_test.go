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
		"ImportPathStandard": {
			config: Configuration{
				Symbols: []string{"go/ast.File"},
			},
			expected: &configuration{
				packages: map[string]string{},
				symbols:  map[string]string{"go/ast.File": ""},
			},
		},
		"ImportPathWithDots": {
			config: Configuration{
				Symbols: []string{"golang.org/x/net/context.Context"},
			},
			expected: &configuration{
				packages: map[string]string{},
				symbols:  map[string]string{"golang.org/x/net/context.Context": ""},
			},
		},
		"SymbolsRefactored": {
			config: Configuration{
				Packages: []string{
					"fmt.{Print,Printf,Println}",
				},
			},
			expected: &configuration{
				packages: map[string]string{
					"fmt.Print":   "",
					"fmt.Printf":  "",
					"fmt.Println": "",
				},
				symbols: map[string]string{},
			},
		},
		"Replacements": {
			config: Configuration{
				Symbols: []string{"source/pkg.Symbol=target/pkg.Symbol"},
			},
			expected: &configuration{
				packages: map[string]string{},
				symbols:  map[string]string{"source/pkg.Symbol": "target/pkg.Symbol"},
			},
		},
		"ReplacementsRefactored": {
			config: Configuration{
				Symbols: []string{"source/pkg.{Print,Printf,Println}=source/pkg.{Sprint,Sprintf,Sprintln}"},
			},
			expected: &configuration{
				packages: map[string]string{},
				symbols: map[string]string{
					"source/pkg.Print":   "source/pkg.Sprint",
					"source/pkg.Printf":  "source/pkg.Sprintf",
					"source/pkg.Println": "source/pkg.Sprintln",
				},
			},
		},
		"MissingSymbol": {
			config: Configuration{
				Symbols: []string{"go/ast"},
			},
		},
		"InvalidSourceSymbol": {
			config: Configuration{
				Symbols: []string{"go/ast=replacement/pkg.Type"},
			},
		},
		"InvalidReplacementSymbol": {
			config: Configuration{
				Symbols: []string{"go/ast.Type=replacement/pkg"},
			},
		},
		"TooManySelectors": {
			config: Configuration{
				Symbols: []string{"go/ast.Foo.Bar"},
			},
		},
		"TooManyReplacements": {
			config: Configuration{
				Symbols: []string{"source/pkg.Symbol=target/pkg.Symbol=intruder/pkg.Symbol"},
			},
		},
		"MismatchedRefactor": {
			config: Configuration{
				Symbols: []string{"source/pkg.{Print,Printf,Println}=source/pkg.{Sprint,Sprintf}"},
			},
		},
	}

	for name := range testcases {
		testcase := testcases[name]
		t.Run(name, func(t *testing.T) {
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
