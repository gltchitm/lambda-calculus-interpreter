package tokenizer

const CommentSymbol = ";"
const LambdaSymbol = "λ"

var Symbols = []string{" ", "(", ")", ";", ".", "=", `\`, LambdaSymbol, CommentSymbol}

func IsLambdaSymbol(char string) bool {
	return char == LambdaSymbol || char == `\`
}
