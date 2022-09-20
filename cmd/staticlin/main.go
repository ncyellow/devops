// Package main содержит пользовательский статический анализатор. Где подключены
// 	1."golang.org/x/tools/go/analysis/passes/printf"
//	2."golang.org/x/tools/go/analysis/passes/shadow"
//	3."golang.org/x/tools/go/analysis/passes/structtag"
//  4. из пакета statickcheck все проверки вида SA, ST. QF, S1
//  5. пользовательский анализатор на использование os.Exit в main
//  6. Публичный анализатор sqlrows + forcetypeassert
// Стандартное использование main .\...

package main

import (
	"github.com/gostaticanalysis/forcetypeassert"
	"github.com/gostaticanalysis/sqlrows/passes/sqlrows"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	analysis2 "github.com/ncyellow/devops/internal/analysis"
)

func main() {

	myChecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		analysis2.ExitCheckAnalyzer,
		sqlrows.Analyzer,
		forcetypeassert.Analyzer,
	}

	// Как я проверил словарь Analyzers содержит только SA проверки, добавляем все что в нем есть
	for _, v := range staticcheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}

	// Просят добавить несколько добавить - добавим все
	for _, v := range stylecheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}

	for _, v := range simple.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}

	// Просят добавить несколько добавить - добавим все
	for _, v := range quickfix.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}

	multichecker.Main(
		myChecks...,
	)
}
