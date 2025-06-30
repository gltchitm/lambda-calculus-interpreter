package interpreter

import "github.com/gltchitm/lambda-calculus-interpreter/parser"

func generateChurchNumerals(from, to int) []parser.Expression {
	numerals := []parser.Expression{}

	f := parser.NewVariable("f", true)
	x := parser.NewVariable("x", true)

	var innerApplication parser.Expression = x

	for i := 0; i <= to; i++ {
		// We can't start directly at from so we just avoid returning any
		// of the Church numerals lower than what was requested
		if i >= from {
			numerals = append(numerals, parser.NewFunction(
				f,
				parser.NewFunction(x, innerApplication),
			))
		}
		innerApplication = parser.NewApplication(f, innerApplication)
	}

	return numerals
}
