package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gltchitm/lambda-calculus-interpreter/interpreter"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	i := interpreter.NewInterpreter()

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		line = regexp.MustCompile(`\s`).ReplaceAllString(line, " ")

		line = strings.TrimSpace(line)
		if line == "exit" {
			break
		}

		result := i.ExecuteLine(line)
		if result != nil {
			fmt.Printf("%v\n", *result)
		}
	}
}
