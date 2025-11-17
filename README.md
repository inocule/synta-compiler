## Synta Lexical Analyzer

Synta is an experimental lexer for a hypothetical AI-friendly programming language. It focuses on rich keyword coverage (async, agentic, AI verbs, etc.), precise token metadata, and ergonomics for downstream parsers or tooling.

### Features
- Handles identifiers, numbers (int/float), strings, operators, and delimiters.
- Supports single-line `//` and block `/~ ~ /` comments.
- Tracks line/column for every token; newline tokens can be included or filtered.
- Recognizes AI-specific keywords like `think`, `reason`, `ask`, `observe`, alongside conventional control-flow, declaration, concurrency, and I/O terms.

### Project Layout
- `main.go` – demo entry point that feeds sample Synta code through the lexer and prints token diagnostics.
- `lexer/lexer.go` – core lexer implementation (`Lexer` struct, tokenization logic).
- `token/token.go` – token type definitions, keyword lookup table, and helper utilities.
- `lexer/lexer_test.go` – unit tests that assert tokenization behavior.

### Getting Started
1. Install Go 1.21+.
2. From `lexical/`, run the demo:
   ```
   go run .
   ```
   You will see a formatted table of tokens for the sample program defined in `main.go`.
3. Run tests:
   ```
   go test ./...
   ```

### Using the Lexer in Your Code
```go
src := `bind x := 10`
lex := lexer.New(src)
tokens := lex.Tokenize()
for _, tok := range tokens {
    fmt.Printf("%s -> %q (%d,%d)\n",
        token.TokenNames[tok.Type], tok.Lexeme, tok.Line, tok.Column)
}
```

### Extending the Lexer
- Add new keywords by updating the `Keywords` map and `const` block in `token/token.go`.
- Add operators or punctuation inside the `switch` in `lexer.Tokenize`.
- Adjust whitespace/comment handling via `skipWhitespace` or `skipComment`.

### License
Specify your license terms here (default is "All rights reserved" if unspecified).

