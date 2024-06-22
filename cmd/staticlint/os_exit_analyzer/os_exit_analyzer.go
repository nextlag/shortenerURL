package os_exit_analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OSExitAnalyzer checks for the presence of os.Exit() in the file and the main function.
var OSExitAnalyzer = &analysis.Analyzer{
	Name: "exitAnalyzer",
	Doc:  "Checks if there any os.Exit implementations in code",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Check that the package is called "main"
		if file.Name.Name != "main" {
			continue
		}

		// Looking for functions and calls os.Exit
		ast.Inspect(file, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					if call, ok := n.(*ast.CallExpr); ok {
						if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
							if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "os" && sel.Sel.Name == "Exit" {
								pass.Reportf(n.Pos(), "found os.Exit")
							}
						}
					}
					return true
				})
			}
			return true
		})
	}
	return nil, nil
}
