package anathema

import (
	"fmt"
	"go/ast"
	"strings"

	"dmitri.shuralyov.com/go/generated"
	"golang.org/x/tools/go/analysis"
)

func Analysis(c *Configuration) *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "anathema",
		Doc:  "Flags the use of symbols that have been marked as forbidden.",
		Run:  runner(c),
	}
	a.Flags.StringVar(&configPath, "config", "", "path to the configuration file")
	return a
}

func runner(config *Configuration) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		config = getConfig(config)
		c, err := config.validate()
		if err != nil {
			return nil, err
		}

		for _, file := range pass.Files {
			path := pass.Fset.File(file.Package).Name()
			isGenerated, err := generated.ParseFile(path)
			if err != nil {
				return nil, err
			} else if isGenerated {
				continue
			}

			checkImports(pass, c, file)
			checkSymbols(pass, c, file)
		}

		return nil, nil
	}
}

func checkImports(pass *analysis.Pass, c *configuration, file *ast.File) {
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		repl, ok := c.packages[path]
		if !ok {
			continue
		}

		d := analysis.Diagnostic{
			Pos: imp.Pos(),
			End: imp.End(),
		}
		if repl == "" {
			d.Message = fmt.Sprintf("%s should not be used", path)
		} else {
			d.Message = fmt.Sprintf("%s should be replaced with %s", path, repl)
			d.SuggestedFixes = []analysis.SuggestedFix{
				{
					Message: fmt.Sprintf("Replace import of %s with %s", path, repl),
					TextEdits: []analysis.TextEdit{
						{
							Pos:     imp.Pos(),
							End:     imp.End(),
							NewText: []byte(fmt.Sprintf(`"%s"`, repl)),
						},
					},
				},
			}
		}
		pass.Report(d)
	}
}

func checkSymbols(pass *analysis.Pass, c *configuration, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		se, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if _, ok = se.X.(*ast.Ident); !ok {
			return true
		}

		obj, ok := pass.TypesInfo.Uses[se.Sel]
		if !ok {
			return true
		} else if obj.Pkg() == nil {
			return true
		}

		path := obj.Pkg().Path()
		if idx := strings.LastIndex(obj.Pkg().Path(), "vendor/"); idx > 0 {
			path = path[idx+7:]
		}

		symbol := fmt.Sprintf("%s.%s", path, obj.Name())
		repl, ok := c.symbols[symbol]
		if !ok {
			return true
		}

		d := analysis.Diagnostic{
			Pos: se.Pos(),
			End: se.End(),
		}
		if repl == "" {
			d.Message = fmt.Sprintf("%s should not be used", symbol)
		} else {
			d.Message = fmt.Sprintf("%s should be replaced with %s", symbol, repl)
		}
		pass.Report(d)

		return true
	})
}
