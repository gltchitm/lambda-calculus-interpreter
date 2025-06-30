package interpreter

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gltchitm/lambda-calculus-interpreter/parser"
	"github.com/gltchitm/lambda-calculus-interpreter/tokenizer"
)

type Interpreter struct {
	variables map[string]parser.Expression
}

func NewInterpreter() *Interpreter {
	return &Interpreter{map[string]parser.Expression{}}
}

func (i *Interpreter) ExecuteLine(line string) *string {
	var result string

	tokens := tokenizer.Tokenize(line)

	expression, err := parser.Parse(tokens)
	if err != nil {
		result = color.HiRedString(err.Error())
		goto done
	}

	if expression != nil {
		populate, isPopulate := expression.(*parser.Populate)
		if isPopulate {
			numerals := generateChurchNumerals(populate.From, populate.To)

			notice := ""

			for index, numeral := range numerals {
				key := fmt.Sprint(index + populate.From)

				_, isAlreadyDefined := i.variables[key]
				if isAlreadyDefined {
					notice = "\n(Skipped one or more numbers that were already defined)"
					continue
				}

				i.variables[fmt.Sprint(index+populate.From)] = numeral
			}

			result = fmt.Sprintf(
				"Populated numbers from %d to %d.%s",
				populate.From,
				populate.To,
				notice,
			)

			goto done
		}

		_, isRootVariable := expression.(*parser.Variable)

		expression = i.resolveFreeVariables(expression)
		expression = i.resolveRuns(expression)
		i.resolveNameAmbiguities(expression)

		assignment, isAssignment := expression.(*parser.Assignment)
		if isAssignment {
			_, alreadyAssigned := i.variables[assignment.To.Name]
			if alreadyAssigned {
				result = fmt.Sprintf("%v is already defined", assignment.To)
			} else {
				i.variables[assignment.To.Name] = assignment.Value

				result = fmt.Sprintf(
					"Added %v as %v",
					assignment.Value,
					assignment.To,
				)
			}
		} else if expression != nil {
			// Do not attempt to resolve the resulting expression to a free
			// variable if the user explicitly entered the variable name
			if !isRootVariable {
				for name, value := range i.variables {
					if compareExpressions(expression, value) {
						expression = parser.NewVariable(name, false)
						break
					}
				}
			}

			result = fmt.Sprintf("%v", expression)
		}
	}

done:
	if result == "" {
		return nil
	} else {
		return &result
	}
}
