// Package main includes the following analyzers:
// exitanalyzer - checks for the presence of os.Exit() in functions and the main file,
// errcheck - checks for error handling,
// analysis - a standard package of linters.
package main

import (
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/nextlag/shortenerURL/cmd/staticlint/analyzer"
)

func main() {
	var checks []*analysis.Analyzer

	// Add all analyzers from the staticcheck package.
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	// Add additional analyzers.
	checks = append(
		checks,
		printf.Analyzer,         // Check for printf-like functions for correct formatting directives.
		shadow.Analyzer,         // Check for shadowed variables.
		structtag.Analyzer,      // Check for struct tags.
		errcheck.Analyzer,       // Check for error handling.
		analyzer.OSExitAnalyzer, // Custom analyzer to check for os.Exit() calls.
	)

	// Run the multichecker with the specified analyzers.
	multichecker.Main(
		checks...,
	)
}
