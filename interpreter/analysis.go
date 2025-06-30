package interpreter

import (
	"maps"
	"reflect"

	"github.com/gltchitm/lambda-calculus-interpreter/parser"
)

func unexpectedType(expression parser.Expression) string {
	typeName := "<nil>"
	if expression != nil {
		typeName = reflect.ValueOf(expression).Type().String()
	}

	return "unexpected expression of type: " + typeName
}

func cloneExpression(
	expression parser.Expression,
	variableMap map[*parser.Variable]parser.Expression,
) parser.Expression {
	switch realExpression := expression.(type) {
	case *parser.Application:
		return parser.NewApplication(
			cloneExpression(realExpression.Left, variableMap),
			cloneExpression(realExpression.Right, variableMap),
		)
	case *parser.Function:
		variable := parser.NewVariable(realExpression.Parameter.Name, true)

		newVariableMap := maps.Clone(variableMap)
		newVariableMap[realExpression.Parameter] = variable

		return parser.NewFunction(
			variable,
			cloneExpression(realExpression.Body, newVariableMap),
		)
	case *parser.Variable:
		variable, isBound := variableMap[realExpression]
		if isBound {
			return variable
		} else if realExpression.IsBound {
			// This bound variable is effectively free
			return realExpression
		}

		return parser.NewVariable(realExpression.Name, false)
	default:
		panic(unexpectedType(expression))
	}
}

func (i *Interpreter) findLeftmostRedex(expression *parser.Expression) *parser.Expression {
	switch realExpression := (*expression).(type) {
	case *parser.Application:
		_, isLeftFunction := realExpression.Left.(*parser.Function)
		if isLeftFunction {
			return expression
		}

		left := i.findLeftmostRedex(&realExpression.Left)
		if left != nil {
			return left
		}

		return i.findLeftmostRedex(&realExpression.Right)
	case *parser.Function:
		return i.findLeftmostRedex(&realExpression.Body)
	default:
		return nil
	}
}

func collectAmbiguousBoundVariables(
	expression parser.Expression,
	boundVariables map[string]*parser.Variable,
	ambiguousVariables []*parser.Variable,
) []*parser.Variable {
	switch realExpression := expression.(type) {
	case *parser.Application:
		return append(collectAmbiguousBoundVariables(
			realExpression.Left,
			boundVariables,
			ambiguousVariables,
		), collectAmbiguousBoundVariables(
			realExpression.Right,
			boundVariables,
			ambiguousVariables,
		)...)
	case *parser.Assignment:
		return collectAmbiguousBoundVariables(
			realExpression.Value,
			boundVariables,
			ambiguousVariables,
		)
	case *parser.Function:
		newBoundVariables := maps.Clone(boundVariables)
		newBoundVariables[realExpression.Parameter.Name] = realExpression.Parameter
		return collectAmbiguousBoundVariables(
			realExpression.Body,
			newBoundVariables,
			ambiguousVariables,
		)
	case *parser.Variable:
		for name, boundVariable := range maps.All(boundVariables) {
			if realExpression != boundVariable && realExpression.Name == name {
				ambiguousVariables = append(ambiguousVariables, boundVariable)
				break
			}
		}
	default:
		panic(unexpectedType(expression))
	}

	return ambiguousVariables
}

func collectVariableNames(expression parser.Expression, names []string) []string {
	switch realExpression := expression.(type) {
	case *parser.Application:
		return append(
			collectVariableNames(realExpression.Left, names),
			collectVariableNames(realExpression.Right, names)...,
		)
	case *parser.Assignment:
		return collectVariableNames(realExpression.Value, names)
	case *parser.Function:
		return collectVariableNames(
			realExpression.Body,
			append(names, realExpression.Parameter.Name),
		)
	case *parser.Variable:
		return append(names, realExpression.Name)
	default:
		panic(unexpectedType(expression))
	}
}

func compareExpressions(lhs parser.Expression, rhs parser.Expression) bool {
	if lhs == rhs {
		return true
	}

	switch realLhs := lhs.(type) {
	case *parser.Application:
		realRhs, isRightApplication := rhs.(*parser.Application)
		if !isRightApplication {
			break
		}

		return compareExpressions(realLhs.Left, realRhs.Left) &&
			compareExpressions(realLhs.Right, realRhs.Right)
	case *parser.Function:
		realRhs, isRightFunction := rhs.(*parser.Function)
		if !isRightFunction {
			break
		}

		return compareExpressions(realLhs.Parameter, realRhs.Parameter) &&
			compareExpressions(realLhs.Body, realRhs.Body)
	case *parser.Variable:
		realRhs, isRightVariable := rhs.(*parser.Variable)
		if !isRightVariable {
			break
		}

		return *realLhs == *realRhs
	default:
		panic(unexpectedType(lhs))
	}

	return false
}
