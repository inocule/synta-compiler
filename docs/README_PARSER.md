A complete lexical and syntax analyzer for the Synta programming language.

## Project Structure

```
synta-lexical/
├── cmd/
│   ├── synta-lex/      # Lexical analyzer CLI
│   │   └── main.go
│   └── synta-parse/    # Syntax analyzer CLI
│       └── main.go
├── lexer/
│   └── lexer.go        # Tokenizer implementation
├── token/
│   └── token.go        # Token definitions
├── parser/
│   ├── parser.go       # Recursive descent parser
│   ├── ast.go          # AST node definitions
│   └── errors.go       # Error handling
├── bin/                # Compiled executables (generated)
├── Makefile            # Build automation
└── go.mod
```

## Installation

```bash
# Initialize Go module (if not done)
go mod init synta-lexical

# Build both tools
make build
```

This creates:
- `bin/synta-lex` - Lexical analyzer
- `bin/synta-parse` - Syntax analyzer

## Usage

### Option 1: Run Complete Pipeline (Recommended)

```bash
make all FILE=yourcode.synta
```

This runs both lexer and parser in sequence.

### Option 2: Run Steps Individually

**Step 1: Lexical Analysis**
```bash
./bin/synta-lex -input code.synta -output tokens.json -errors lex-errors.txt
```

**Step 2: Syntax Analysis**
```bash
./bin/synta-parse -input tokens.json -tree parse-tree.txt -ast ast.json -errors parse-errors.txt
```

### Option 3: Using Make Targets

```bash
# Build only
make build

# Run lexer
make lex FILE=code.synta

# Run parser
make parse

# Clean output files
make clean

# Run quick test
make test
```

## Output Files

### From Lexer (`synta-lex`):
1. **tokens.json** - All tokens in JSON format
2. **lex-errors.txt** - Lexical errors (if any)

### From Parser (`synta-parse`):
1. **parse-tree.txt** - Human-readable parse tree
2. **ast.json** - Abstract Syntax Tree in JSON
3. **parse-errors.txt** - Syntax errors (if any)
4. **parse-debug.txt** - Detailed parsing steps

## Examples

### Example 1: Simple Variable Binding

**Input (test.synta):**
```synta
bind x := 42;
bind y := x + 10;
print(x, y);
```

**Run:**
```bash
make all FILE=test.synta
```

**Output (parse-tree.txt):**
```
PROGRAM
  BIND_STMT (x) [1:6]
    LITERAL (42) [1:11]
  BIND_STMT (y) [2:6]
    BINARY_EXPR (+) [2:11]
      IDENTIFIER (x) [2:11]
      LITERAL (10) [2:15]
  PRINT_STMT (print) [3:1]
    IDENTIFIER (x) [3:7]
    IDENTIFIER (y) [3:10]
```

### Example 2: Agent Declaration

**Input (agent.synta):**
```synta
@agent DataProcessor {
    fn process(data) {
        think("Processing data...");
        return data;
    }
}
```

**Run:**
```bash
./bin/synta-lex -input agent.synta
./bin/synta-parse
```

### Example 3: Control Flow

**Input (control.synta):**
```synta
bind count := 0;

while (count < 10) {
    if (count % 2 == 0) {
        print("Even:", count);
    }
    count := count + 1;
}
```

## Error Handling

The parser provides detailed error messages:

**Example Error:**
```
Line 5:12 - Parse error: Expected ')' after if condition (at 'then')
```

**Debug Output (parse-debug.txt):**
```
Parsing program
Parsing declaration at token: bind
Parsing bind statement
Parsing expression statement
Parsing if statement
ERROR: Line 5:12 - Parse error: Expected ')' after if condition (at 'then')
```

## CLI Flags

### synta-lex
```bash
-input string    Input .synta source file (required)
-output string   Output token file (default: "tokens.json")
-errors string   Lexical errors file (default: "lex-errors.txt")
```

### synta-parse
```bash
-input string    Input token file (default: "tokens.json")
-tree string     Output parse tree file (default: "parse-tree.txt")
-ast string      Output AST JSON file (default: "ast.json")
-errors string   Parse errors file (default: "parse-errors.txt")
-debug string    Debug log file (default: "parse-debug.txt")
```

## Development

### Adding New Language Features

1. **Update token.go** - Add new token types
2. **Update lexer.go** - Add tokenization rules
3. **Update parser.go** - Add parsing methods
4. **Update ast.go** - Add new AST node types

### Testing

```bash
# Quick test with sample code
make test

# Test specific file
make all FILE=examples/your-test.synta

# View parse tree
cat parse-tree.txt

# View errors
cat parse-errors.txt lex-errors.txt
```

## Troubleshooting

**Problem:** `command not found: synta-lex`
```bash
# Make sure you ran: make build
# Or use full path: ./bin/synta-lex
```

**Problem:** `Error reading token file`
```bash
# Make sure lexer ran first:
make lex FILE=yourfile.synta
# Then run parser:
make parse
```

**Problem:** Parse errors but unsure why
```bash
# Check debug log for detailed parsing steps:
cat parse-debug.txt
```