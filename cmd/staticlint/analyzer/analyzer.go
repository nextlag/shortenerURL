package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OSExitAnalyzer checks for the presence of os.Exit() in the file and the main function.
var OSExitAnalyzer = &analysis.Analyzer{
	Name: "osExitAnalyzer",
	Doc:  "Checks for os.Exit calls in the main function of package main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			ast.Inspect(file, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
					ast.Inspect(fn.Body, func(n ast.Node) bool {
						if call, ok := n.(*ast.CallExpr); ok {
							if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
								if pkg, ok := fun.X.(*ast.Ident); ok && pkg.Name == "os" && fun.Sel.Name == "Exit" {
									pass.Reportf(call.Pos(), "os.Exit call is prohibited in main function")
								}
							}
						}
						return true
					})
				}
				return true
			})
		}
	}
	return nil, nil
}
