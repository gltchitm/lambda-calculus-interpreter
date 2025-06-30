package tokenizer

import "slices"

func Tokenize(line string) *TokenStream {
	tokens := []Token{}

	tokenValue := ""

	for _, char := range line {
		char := string(char)

		if char == CommentSymbol {
			break
		}

		if slices.Contains(Symbols, char) {
			if len(tokenValue) > 0 {
				tokens = append(tokens, NewToken(TokenTypeIdentifier, tokenValue))
				tokenValue = ""
			}

			if char != " " {
				tokens = append(tokens, NewToken(TokenTypeSymbol, char))
			}
		} else {
			tokenValue += char
		}
	}

	if len(tokenValue) > 0 {
		tokens = append(tokens, NewToken(TokenTypeIdentifier, tokenValue))
	}

	return NewTokenStream(tokens)
}
