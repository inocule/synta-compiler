// main.go
package main

import (
	"fmt"
	"synta/lexer"
	"synta/token"
)

func main() {
	code := `
	bind x := 10;
	const PI := 3.14;
	
	fn calculate(a, b) {
		if a > b {
			return a + b;
		} else {
			return a - b;
		}
	}
	
	// This is a comment
	for i := 0; i < 5; i++ {
		print("Hello");
	}
	
	// AI keywords
	think "What should I do next?";
	reason about_problem;
	`

	l := lexer.New(code)
	tokens := l.Tokenize()

	fmt.Println("Synta Lexer - Token Output")
	fmt.Println("===========================")
	for _, tok := range tokens {
		if tok.Type == token.NEWLINE {
			continue // Skip printing newlines for cleaner output
		}

		typeName := token.TokenNames[tok.Type]
		if typeName == "" {
			typeName = fmt.Sprintf("TokenType(%d)", tok.Type)
		}

		fmt.Printf("%-15s %-20s Line: %d, Col: %d\n",
			typeName, fmt.Sprintf("'%s'", tok.Lexeme), tok.Line, tok.Column)
	}
}
