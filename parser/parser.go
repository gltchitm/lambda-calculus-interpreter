package parser

import (
	"errors"
	"maps"

	"github.com/gltchitm/lambda-calculus-interpreter/tokenizer"
)

func parse(
	tokens *tokenizer.TokenStream,
	isRoot,
	needsClose,
	isInRun bool,
	boundVariables map[string]*Variable,
) (Expression, error) {
	var expression Expression

	isClosed := false
	skippedClose := false

	for tokens.HasNextToken(tokenizer.ExpectAnyToken) {
		token, err := tokens.Peek(tokenizer.ExpectAnyToken)
		if err != nil {
			return nil, err
		}

		// It's not our responsibility to handle this closing parenthesis,
		// so we return control to the responsible invoker
		if token.Is(tokenizer.TokenTypeSymbol, ")") && !needsClose {
			skippedClose = true
			break
		}

		token, err = tokens.Read(tokenizer.ExpectAnyToken)
		if err != nil {
			return nil, err
		}

		var it Expression

		if token.TokenType == tokenizer.TokenTypeSymbol {
			if token.Value == "(" {
				it, err = parse(tokens, false, true, isInRun, boundVariables)
				if err != nil {
					return nil, err
				}
			} else if token.Value == ")" {
				isClosed = true
				break
			} else if tokenizer.IsLambdaSymbol(token.Value) {
				parameterToken, err := tokens.Read(tokenizer.TokenTypeIdentifier)
				if err != nil {
					return nil, err
				}

				dotToken, err := tokens.Read(tokenizer.TokenTypeSymbol)
				if err != nil {
					return nil, err
				} else if dotToken.Value != "." {
					return nil, errors.New("expected '.' but got '" + dotToken.Value + "'")
				}

				variable := NewVariable(parameterToken.Value, true)

				newBoundVariables := maps.Clone(boundVariables)
				newBoundVariables[parameterToken.Value] = variable

				body, err := parse(tokens, false, false, isInRun, newBoundVariables)
				if err != nil {
					return nil, err
				} else if body == nil {
					return nil, errors.New("expected expression after function header")
				}

				it = NewFunction(variable, body)
			} else if token.Value == "=" && isRoot {
				variable, isVariable := expression.(*Variable)
				if !isVariable {
					return nil, errors.New("invalid left-hand side in assignment")
				}

				value, err := parse(tokens, false, false, isInRun, boundVariables)
				if err != nil {
					return nil, err
				} else if value == nil {
					return nil, errors.New("invalid right-hand side in assignment")
				}

				expression = NewAssignment(variable, value)

				continue
			} else {
				return nil, errors.New("unexpected symbol '" + token.Value + "'")
			}
		} else if token.Is(tokenizer.TokenTypeIdentifier, "run") {
			if isInRun {
				return nil, errors.New("cannot nest run statements")
			}

			body, err := parse(tokens, false, false, true, boundVariables)
			if err != nil {
				return nil, err
			} else if body == nil {
				return nil, errors.New("expected expression after 'run'")
			}

			it = NewRun(body)
		} else if token.Is(tokenizer.TokenTypeIdentifier, "populate") {
			if !isRoot {
				return nil, errors.New("'populate' cannot appear here")
			}

			from, err := tokens.Read(tokenizer.TokenTypeIdentifier)
			if err != nil {
				return nil, err
			}

			fromValue, err := parsePopulateBound("lower", from.Value)
			if err != nil {
				return nil, err
			}

			to, err := tokens.Read(tokenizer.TokenTypeIdentifier)
			if err != nil {
				return nil, err
			}

			toValue, err := parsePopulateBound("upper", to.Value)
			if err != nil {
				return nil, err
			}

			if fromValue > toValue {
				return nil, errors.New("lower bound cannot be less than upper bound")
			}

			it = NewPopulate(fromValue, toValue)

			rest, err := parse(tokens, false, false, isInRun, boundVariables)
			if err != nil {
				return nil, err
			} else if rest != nil {
				return nil, errors.New("unexpected expression(s) after 'populate'")
			}
		} else {
			variable, isBound := boundVariables[token.Value]
			if isBound {
				it = variable
			} else {
				it = NewVariable(token.Value, false)
			}
		}

		if expression == nil {
			expression = it
		} else {
			expression = NewApplication(expression, it)
		}
	}

	if (isClosed && !needsClose) || (skippedClose && isRoot) {
		return nil, errors.New("unexpected closing parenthesis")
	} else if !isClosed && needsClose {
		return nil, errors.New("expected closing parenthesis")
	}

	return expression, nil
}

func Parse(tokens *tokenizer.TokenStream) (Expression, error) {
	return parse(tokens, true, false, false, map[string]*Variable{})
}
