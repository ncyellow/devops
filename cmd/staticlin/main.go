// Package main содержит пользовательский multichecker.
// Где подключены:
// 	1."golang.org/x/tools/go/analysis/passes/printf" - проверка на корректность различных fmt.Print функций
//	2."golang.org/x/tools/go/analysis/passes/shadow" - проверка на создание "затененных" переменных
//	3."golang.org/x/tools/go/analysis/passes/structtag" - проверка тегов для для полей структур
//  4. из пакета statickcheck все проверки вида SA, ST. QF, S1
//  5. пользовательский анализатор на использование os.Exit в пакете main в функции main
//  6. Публичный анализатор sqlrows - проверяем корректное использование sqlrow
//  7. Публичный анализатор forcetypeassert - проверка принудительных проверок типов
// Стандартное использование аналогично go vet .\...
// main .\...
package main

import (
	"github.com/gostaticanalysis/forcetypeassert"
	"github.com/gostaticanalysis/sqlrows/passes/sqlrows"
	"github.com/ncyellow/devops/internal/analysis/exit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {

	myChecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		exit.Analyzer,
		sqlrows.Analyzer,
		forcetypeassert.Analyzer,
	}

	// Добавляем все проверки типа SA
	for _, v := range staticcheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}
	// Добавляем все проверки типа ST
	for _, v := range stylecheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}
	// Добавляем все проверки типа ST100
	for _, v := range simple.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}
	// Добавляем все проверки типа QF1001
	for _, v := range quickfix.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}

	multichecker.Main(
		myChecks...,
	)
}
