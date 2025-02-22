// Package main содержит примеры использования анализа кода с помощью различных анализаторов, включая
// стандартные анализаторы пакета golang.org/x/tools/go/analysis и кастомные анализаторы.
//
// В этом пакете используется multichecker для объединения нескольких анализаторов в одном процессе
// для проверки исходного кода на различные ошибки и недочеты.
//
// Пример анализаторов:
//   - NoExitInMainAnalyzer: Анализатор, который проверяет, что os.Exit не вызывается в функции main.
//   - Nilfunc: Стандартный анализатор, который проверяет на nil-проверки для функций и методов.
//   - Printf: Стандартный анализатор, который проверяет правильность использования форматов в функциях типа fmt.Printf.
//   - Анализаторы из staticcheck: Статические анализаторы, такие как SA, ST, SX классов из пакета staticcheck.
package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"honnef.co/go/tools/staticcheck"
)

// NoExitInMainAnalyzer — собственный анализатор, который проверяет, что в функции main
// пакета main не вызывается os.Exit напрямую. Это полезно для предотвращения использования os.Exit
// в основной функции программы, что может привести к незавершенному корректному
// завершению программы.
var NoExitInMainAnalyzer = &analysis.Analyzer{
	Name: "noexitinmain",
	Doc:  "checks that os.Exit is not called in main function of the main package",
	Run:  runNoExitInMainAnalyzer,
}

// runNoExitInMainAnalyzer — функция, которая выполняет проверку для обнаружения вызова os.Exit в функции main.
func runNoExitInMainAnalyzer(pass *analysis.Pass) (interface{}, error) {
	// Проходим по всем файлам пакета
	for _, file := range pass.Files {
		// Проверяем только функции main
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем функцию main
			funcDecl, ok := n.(*ast.FuncDecl)
			if ok && funcDecl.Name.Name == "main" {
				// Если это функция main, проверяем на вызов os.Exit
				ast.Inspect(funcDecl, func(n ast.Node) bool {
					// Ищем вызовы os.Exit
					callExpr, ok := n.(*ast.CallExpr)
					if ok {
						if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
							// Проверяем, что вызывается именно os.Exit
							if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "os" && sel.Sel.Name == "Exit" {
								pass.Reportf(callExpr.Pos(), "os.Exit should not be called in the main function")
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

// main — точка входа в программу. Настраивает и запускает multichecker,
// который использует несколько анализаторов, включая стандартные и кастомные анализаторы.
func main() {
	var mychecks []*analysis.Analyzer

	// Добавляем анализатор nilfunc из golang.org/x/tools, который проверяет использование nil
	mychecks = append(mychecks, nilfunc.Analyzer)

	// Добавляем анализатор printf из golang.org/x/tools, который проверяет правильность форматирования строк
	mychecks = append(mychecks, printf.Analyzer)

	// Добавляем анализаторы из staticcheck для каждого класса
	for _, v := range staticcheck.Analyzers {
		// Добавляем анализаторы с префиксом "SA"
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			mychecks = append(mychecks, v.Analyzer)
		}

		// Добавляем анализаторы с префиксом "ST"
		if strings.HasPrefix(v.Analyzer.Name, "ST") {
			mychecks = append(mychecks, v.Analyzer)
		}

		// Добавляем анализаторы с префиксом "SX"
		if strings.HasPrefix(v.Analyzer.Name, "SX") {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// Добавляем собственный анализатор запрещающий использовать прямой вызов os.Exit в функции main пакета main в multichecker
	mychecks = append(mychecks, NoExitInMainAnalyzer)

	multichecker.Main(
		mychecks...,
	)
}
