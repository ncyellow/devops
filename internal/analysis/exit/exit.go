// Package exit содержит анализатор использования в main функции os.Exit
package exit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "exit",
	Doc:  "check for using os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// игнорируем пакеты не main
		if file.Name.Name != "main" {
			return nil, nil
		}
		// функцией ast.Inspect обходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				// ast.FuncDecl представляет декларацию функции нас интересует только main, если мы находимся не
				// в main то сразу на выход
				if x.Name.Name != "main" {
					return false
				}
			case *ast.CallExpr:
				// Если это вызов функции, то анализируем что же это было
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := s.X.(*ast.Ident); ok {
						// нас интересует только os.Exit
						if ident.Name == "os" && s.Sel.Name == "Exit" {
							pass.Reportf(x.Pos(), "os.Exit was being detected!")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
