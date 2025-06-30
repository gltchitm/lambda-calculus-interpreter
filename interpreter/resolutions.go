package interpreter

import "github.com/gltchitm/lambda-calculus-interpreter/parser"

func (i *Interpreter) resolveFreeVariables(expression parser.Expression) parser.Expression {
	switch realExpression := expression.(type) {
	case *parser.Application:
		realExpression.Left = i.resolveFreeVariables(realExpression.Left)
		realExpression.Right = i.resolveFreeVariables(realExpression.Right)
		return realExpression
	case *parser.Assignment:
		realExpression.Value = i.resolveFreeVariables(realExpression.Value)
		return realExpression
	case *parser.Function:
		realExpression.Body = i.resolveFreeVariables(realExpression.Body)
		return realExpression
	case *parser.Run:
		realExpression.Body = i.resolveFreeVariables(realExpression.Body)
		return realExpression
	case *parser.Variable:
		// We must not attempt to resolve a bound variable
		if realExpression.IsBound {
			return expression
		}

		value, ok := i.variables[realExpression.Name]
		if !ok {
			return realExpression
		}

		return cloneExpression(value, map[*parser.Variable]parser.Expression{})
	default:
		panic(unexpectedType(expression))
	}
}

func (i *Interpreter) resolveRuns(expression parser.Expression) parser.Expression {
	switch realExpression := expression.(type) {
	case *parser.Application:
		realExpression.Left = i.resolveRuns(realExpression.Left)
		realExpression.Right = i.resolveRuns(realExpression.Right)
	case *parser.Assignment:
		realExpression.Value = i.resolveRuns(realExpression.Value)
	case *parser.Function:
		realExpression.Body = i.resolveRuns(realExpression.Body)
	case *parser.Run:
		i.resolveRedexes(&realExpression.Body)
		return realExpression.Body
	case *parser.Variable:
		break
	default:
		panic(unexpectedType(expression))
	}

	return expression
}

func (i *Interpreter) resolveRedexes(expression *parser.Expression) {
	for {
		redex := i.findLeftmostRedex(expression)
		if redex == nil {
			break
		}

		application := (*redex).(*parser.Application)
		function := application.Left.(*parser.Function)

		*redex = i.applyBetaReduction(
			function.Parameter,
			cloneExpression(application.Right, map[*parser.Variable]parser.Expression{}),
			function.Body,
		)
	}
}

func (i *Interpreter) resolveNameAmbiguities(expression parser.Expression) {
	i.applyAlphaReduction(expression)
}
