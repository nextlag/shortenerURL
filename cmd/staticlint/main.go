// Пакет multichecker включает в себя следующие анализаторы:
// exitanalyzer (проверка наличия os.Exit() в функции и файле main),
// errcheck (проверка обработки ошибок),
// analysis (стандартный пакет линтеров).
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

	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	checks = append(
		checks,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		errcheck.Analyzer,
		analyzer.OSExitAnalyzer,
	)

	multichecker.Main(
		checks...,
	)
}
