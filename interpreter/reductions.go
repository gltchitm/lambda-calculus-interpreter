package interpreter

import (
	"fmt"

	"github.com/gltchitm/lambda-calculus-interpreter/parser"
)

func (i *Interpreter) applyAlphaReduction(expression parser.Expression) {
	ambiguousVariables := collectAmbiguousBoundVariables(
		expression,
		map[string]*parser.Variable{},
		[]*parser.Variable{},
	)

	foundNames := collectVariableNames(expression, []string{})

	conflictingNames := make(set[string])

	for _, name := range foundNames {
		conflictingNames.Add(name)
	}

	for name := range i.variables {
		conflictingNames.Add(name)
	}

	fixed := make(set[*parser.Variable])

	for _, variable := range ambiguousVariables {
		if fixed.Has(variable) {
			continue
		} else if !variable.IsBound {
			panic("attempted to rename free variable")
		}

		subscript := 1

		for conflictingNames.Has(variable.Name + fmt.Sprint(subscript)) {
			subscript++
		}

		variable.Name += fmt.Sprint(subscript)

		fixed.Add(variable)
		conflictingNames.Add(variable.Name)
	}
}

func (i *Interpreter) applyBetaReduction(
	variable *parser.Variable,
	value parser.Expression,
	expression parser.Expression,
) parser.Expression {
	switch realExpression := expression.(type) {
	case *parser.Application:
		return parser.NewApplication(
			i.applyBetaReduction(variable, value, realExpression.Left),
			i.applyBetaReduction(variable, value, realExpression.Right),
		)
	case *parser.Function:
		return parser.NewFunction(
			realExpression.Parameter,
			i.applyBetaReduction(variable, value, realExpression.Body),
		)
	case *parser.Variable:
		if realExpression == variable {
			return value
		}
		return realExpression
	default:
		panic(unexpectedType(expression))
	}
}
