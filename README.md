# Synta Lexical Analyzer

Synta is an experimental lexer for a hypothetical AI-friendly programming language. It focuses on rich keyword coverage (async, agentic, AI verbs, etc.), precise token metadata, and ergonomics for downstream parsers or tooling.

## Features

- Handles identifiers, numbers (int/float), strings, operators, and delimiters
- Supports single-line `//` and block `/~ ~/` comments (captured as tokens)
- Tracks line/column for every token; newline tokens can be included or filtered
- Recognizes AI-specific keywords like `think`, `reason`, `ask`, `observe`
- Supports `@agent` and `@task` decorators for AI agent definitions
- Uses `->` arrow operator for AI invocation (e.g., `AICoder -> "prompt"`)
- Type keywords: `int`, `float`, `char`, `bool`, `str`
- Distinct binding (`:=`) and assignment (`=:`) operators for unambiguous semantics

## Project Layout

```
synta-lexical/
├── main.go              # Demo entry point
├── go.mod
├── lexer/
│   ├── lexer.go         # Core lexer implementation
│   └── lexer_test.go    # Unit tests
└── token/
    └── token.go         # Token types, keywords, helpers
```

## Getting Started

1. Install Go 1.21+ (check with `go version`)
2. From `synta-lexical/`, run the demo:
   ```bash
   go run .
   ```
3. After changes, rebuild with `go build`

## Language Syntax Overview

```synta
// Agent definition
@agent AICoder {
    role: "Chatbot with GitHub MCP access"
    tools: [github_mcp, slm_chatbot]
    model: local/llama-3.1-8b.gguf
    mode: hybrid
    sys_prompt: "You are a helpful coding assistant"
}

// Agent invocation
@task {
    response := AICoder -> "Fix this bug"
    print(response)
}

// Variable binding and assignment
bind int x =: 10      // declare mutable x
const float PI =: 3.14 // declare immutable PI
x =: 20               // reassign x

// Functions
fn calculate(a, b) {
    if a > b {
        return a + b
    } else {
        return a - b
    }
}

// Comments
// Single-line comment
/~ Multi-line
   comment ~/

// Loops
int i =: 0
while i < 5 {
    print("Hello World")
    i++
}
```

## Operator Reference

| Operator | Token | Meaning |
|----------|-------|---------|
| `:=` | BIND_ASSIGN | Declare and bind |
| `=:` | ASSIGN | Assign to existing |
| `->` | ARROW | AI agent invocation |
| `==` | EQ | Equality |
| `!=` | NEQ | Not equal |
| `&&` | AND | Logical and |
| `||` | OR | Logical or |
| `++` | INCREMENT | Increment |
| `--` | DECREMENT | Decrement |

## Using the Lexer in Your Code

```go
src := `bind int x =: 10`
lex := lexer.New(src)
tokens := lex.Tokenize()
for _, tok := range tokens {
    fmt.Printf("%s -> %q (%d,%d)\n",
        token.TokenNames[tok.Type], tok.Lexeme, tok.Line, tok.Column)
}
```

## Extending the Lexer

- **Add keywords**: Update `Keywords` map and `const` block in `token/token.go`
- **Add operators**: Modify the `switch` in `lexer.Tokenize()`
- **Add decorators**: Extend the `@` handling in lexer to recognize new patterns