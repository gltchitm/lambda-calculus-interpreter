package tokenizer

import "errors"

const ExpectAnyToken = ""

type TokenStream struct {
	Tokens []Token
}

func NewTokenStream(tokens []Token) *TokenStream {
	return &TokenStream{tokens}
}

func isTokenOfType(token Token, expected string) bool {
	return expected == ExpectAnyToken || token.TokenType == expected
}

func (ts *TokenStream) Read(expected string) (*Token, error) {
	token, err := ts.Peek(expected)
	if err != nil {
		return nil, err
	}

	ts.Tokens = ts.Tokens[1:]

	return token, nil
}

func (ts *TokenStream) Peek(expected string) (*Token, error) {
	if len(ts.Tokens) == 0 {
		return nil, errors.New("unexpected end of input")
	}

	token := ts.Tokens[0]

	if !isTokenOfType(token, expected) {
		return nil, errors.New("expected " + expected + " but found " + token.TokenType)
	}

	return &token, nil
}

func (ts *TokenStream) HasNextToken(expected string) bool {
	return len(ts.Tokens) > 0 && isTokenOfType(ts.Tokens[0], expected)
}
