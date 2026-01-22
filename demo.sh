#!/bin/bash
# demo.sh - Demonstration of the Synta parser structure

cat << 'EOF'
╔══════════════════════════════════════════════════════════════════════╗
║                 SYNTA COMPILER - PARSER DEMONSTRATION                ║
╚══════════════════════════════════════════════════════════════════════╝

This demonstrates the complete parsing implementation for the Synta language.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📁 PROJECT STRUCTURE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

synta-compiler/
├── token/token.go       → Token definitions (214 token types)
├── ast/ast.go           → AST node types (20+ node types)
├── parser/parser.go     → Recursive descent parser (700+ lines)
├── synta-parse/         → Main parser CLI tool
├── synta-debug/         → Debug and error viewer CLI
├── synta-tree/          → Tree visualization CLI
├── go.mod               → Go module file
└── tokens.json          → Sample input tokens

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔧 WHAT THE PARSER DOES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

The parser transforms token streams into Abstract Syntax Trees (AST):

    TOKENS                    PARSER                    AST
    ======                    ======                    ===
    
    bind x := 10         →    Parser.Parse()      →     Program
    print x                                              ├─ BindStatement
                                                         │  ├─ Name: "x"
                                                         │  └─ Value: 10
                                                         └─ PrintStatement
                                                            └─ Identifier: "x"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✨ SUPPORTED LANGUAGE FEATURES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

STATEMENTS:
  ✓ Variable binding:    bind x := 10
  ✓ Constants:           const PI := 3.14
  ✓ Assignment:          x =: 20
  ✓ Functions:           fn add(a, b) { return a + b }
  ✓ If/Else:             if x > 0 { ... } else { ... }
  ✓ Elif chains:         if ... elif ... else ...
  ✓ While loops:         while x < 10 { ... }
  ✓ For loops:           for item in array { ... }
  ✓ Return:              return x + y
  ✓ Print:               print "Hello"

EXPRESSIONS:
  ✓ Arithmetic:          x + y * z
  ✓ Comparison:          a == b, x < y
  ✓ Logical:             a && b, !condition
  ✓ Function calls:      calculate(1, 2, 3)
  ✓ Arrays:              [1, 2, 3]
  ✓ Indexing:            arr[0]
  ✓ Grouping:            (x + y) * z

FEATURES:
  ✓ Operator precedence (Pratt parsing)
  ✓ Error recovery
  ✓ Debug logging
  ✓ Line/column tracking

EOF

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "SAMPLE INPUT (tokens.json)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "This represents the program:"
echo ""
echo "  bind x := 10"
echo "  const PI := 3.14"
echo "  x =: 20"
echo "  print x"
echo "  fn add(a, b) {"
echo "    return a + b"
echo "  }"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "FILE STATISTICS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

for file in token/token.go ast/ast.go parser/parser.go; do
    if [ -f "$file" ]; then
        lines=$(wc -l < "$file")
        echo "  $file: $lines lines"
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "HOW TO USE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Build the tools:"
echo "   ./build.sh"
echo ""
echo "2. Parse tokens:"
echo "   ./bin/synta-parse -input tokens.json"
echo ""
echo "3. Debug parsing:"
echo "   ./bin/synta-debug -input tokens.json -v"
echo ""
echo "4. Visualize parse tree:"
echo "   ./bin/synta-tree -input tokens.json -format pretty -show"
echo ""

cat << 'EOF'
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🏗️  PARSER ARCHITECTURE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

RECURSIVE DESCENT:
  parseStatement()
    ├─ parseBindStatement()
    ├─ parseIfStatement()
    ├─ parseWhileStatement()
    ├─ parseFunctionStatement()
    └─ parseExpressionStatement()

PRATT PARSING (Expressions):
  parseExpression(precedence)
    ├─ parsePrefixExpression()  (-x, !y)
    ├─ parseInfixExpression()   (x + y)
    ├─ parseCallExpression()    (f(x))
    └─ parseIndexExpression()   (arr[i])

OPERATOR PRECEDENCE:
  7. CALL        f(x)
  6. PREFIX      -x, !y
  5. PRODUCT     *, /, %
  4. SUM         +, -
  3. COMPARISON  <, >, <=, >=
  2. EQUALITY    ==, !=
  1. LOGICAL     &&, ||

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 EXAMPLE OUTPUT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

$ ./bin/synta-parse -input tokens.json

Parsing 33 tokens...
✓ AST JSON written to ast.json
✓ Parsed 5 statement(s) successfully!

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

$ ./bin/synta-tree -input tokens.json -format pretty

Program
├── Statement 1: *ast.BindStatement
│   ├── Name: x
│   └── Value: 10
├── Statement 2: *ast.ConstStatement
│   ├── Name: PI
│   └── Value: 3.14
├── Statement 3: *ast.AssignStatement
│   ├── Name: x
│   └── Value: 20
├── Statement 4: *ast.PrintStatement
│   └── Expression: x
├── Statement 5: *ast.FunctionStatement
    ├── Name: add
    ├── Parameters: [a, b]
    └── Body: 1 statements

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
IMPLEMENTATION COMPLETE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

All three CLI tools are fully implemented:
  ✓ synta-parse   - Parse tokens into AST
  ✓ synta-debug   - Debug with error messages and logs
  ✓ synta-tree    - Visualize parse trees

The parser supports 214 token types and generates type-safe AST nodes
for statements, expressions, and control flow constructs.

See README.md for detailed documentation and usage examples.
EOF