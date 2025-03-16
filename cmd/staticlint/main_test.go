package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestNoExitInMainAnalyzer проверяет, что анализатор NoExitInMainAnalyzer корректно
// обнаруживает вызовы os.Exit в функции main.
func TestNoExitInMainAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, NoExitInMainAnalyzer, "a")
}
