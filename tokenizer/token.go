package tokenizer

const (
	TokenTypeSymbol     = "symbol"
	TokenTypeIdentifier = "identifier"
)

type Token struct {
	TokenType string
	Value     string
}

func NewToken(tokenType, value string) Token {
	return Token{tokenType, value}
}

func (t *Token) Is(tokenType, value string) bool {
	return t.TokenType == tokenType && t.Value == value
}
