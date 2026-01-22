// main.go - Unified Synta Parser CLI
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	parser "synta-compiler/syntax-analyzer/synta-parse/parser"
)

func main() {
	// Define all flags
	inputFile := flag.String("input", "tokens.json", "Input token file")
	treeFile := flag.String("tree", "parse-tree.txt", "Output parse tree file")
	astFile := flag.String("ast", "ast.json", "Output AST JSON file")
	errorsFile := flag.String("errors", "parse-errors.txt", "Parse errors file")
	debugFile := flag.String("debug", "parse-debug.txt", "Debug log file")
	format := flag.String("format", "pretty", "Tree format: pretty, compact, or detailed")
	showConsole := flag.Bool("show", false, "Show tree in console")
	skipAST := flag.Bool("skip-ast", false, "Skip AST JSON generation")
	skipDebug := flag.Bool("skip-debug", false, "Skip debug log generation")

	flag.Parse()

	// Print header
	printHeader()

	// Load tokens
	tokens, err := parser.LoadTokens(*inputFile)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		printUsage()
		os.Exit(1)
	}

	fmt.Printf("📄 Loaded %d tokens from %s\n", len(tokens), *inputFile)

	// Create parser and parse
	p := parser.New(tokens)
	program, errors, debugLog := p.Parse()

	// Handle parsing errors
	if len(errors) > 0 {
		fmt.Printf("\n⚠️  Parsing completed with %d error(s)\n\n", len(errors))

		// Write errors to file
		if err := parser.WriteErrors(*errorsFile, errors); err != nil {
			fmt.Printf("Error writing errors file: %v\n", err)
		} else {
			fmt.Printf("📝 Errors written to %s\n", *errorsFile)
		}

		// Show first few errors in console
		fmt.Println("\nFirst errors:")
		for i, e := range errors {
			if i >= 5 {
				fmt.Printf("... and %d more error(s)\n", len(errors)-5)
				break
			}
			fmt.Printf("  %d. %s\n", i+1, e.Error())
		}

		// Still write debug log if available
		if !*skipDebug {
			if err := parser.WriteDebugLog(*debugFile, debugLog); err != nil {
				fmt.Printf("\n⚠️  Could not write debug file: %v\n", err)
			} else {
				fmt.Printf("\n🔍 Debug log written to %s\n", *debugFile)
			}
		}

		fmt.Println("\n❌ Parsing failed. Cannot generate tree with errors present.")
		os.Exit(1)
	}

	// Success! Generate outputs
	fmt.Printf("\n✅ Parsing completed successfully!\n")
	fmt.Printf("📊 Total statements parsed: %d\n\n", len(program.Statements))

	// Generate and write parse tree
	var treeContent string
	switch *format {
	case "compact":
		treeContent = parser.GenerateCompactTree(program)
	case "detailed":
		treeContent = parser.GenerateDetailedTree(program)
	default:
		treeContent = parser.GeneratePrettyTree(program)
	}

	if err := os.WriteFile(*treeFile, []byte(treeContent), 0644); err != nil {
		fmt.Printf("❌ Error writing tree file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("🌳 Parse tree written to %s (format: %s)\n", *treeFile, *format)

	// Write AST JSON (unless skipped)
	if !*skipAST {
		astJSON, err := json.MarshalIndent(program, "", "  ")
		if err != nil {
			fmt.Printf("⚠️  Error marshaling AST: %v\n", err)
		} else {
			if err := os.WriteFile(*astFile, astJSON, 0644); err != nil {
				fmt.Printf("⚠️  Error writing AST file: %v\n", err)
			} else {
				fmt.Printf("📦 AST JSON written to %s\n", *astFile)
			}
		}
	}

	// Write debug log (unless skipped)
	if !*skipDebug {
		if err := parser.WriteDebugLog(*debugFile, debugLog); err != nil {
			fmt.Printf("⚠️  Could not write debug file: %v\n", err)
		} else {
			fmt.Printf("🔍 Debug log written to %s\n", *debugFile)
		}
	}

	// Write empty errors file to indicate success
	if err := os.WriteFile(*errorsFile, []byte(""), 0644); err != nil {
		fmt.Printf("⚠️  Could not write errors file: %v\n", err)
	} else {
		fmt.Printf("✓  No errors (empty file: %s)\n", *errorsFile)
	}

	// Show tree in console if requested
	if *showConsole {
		fmt.Println("\n" + strings.Repeat("=", 70))
		fmt.Println("PARSE TREE")
		fmt.Println(strings.Repeat("=", 70))
		fmt.Println(treeContent)
		fmt.Println(strings.Repeat("=", 70))
	}

	// Print summary
	printSummary(program, *format)
}

func printHeader() {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("  Synta Syntax Analyzer")
	fmt.Println("  Unified Parser with Full CLI Features")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()
}

func printUsage() {
	fmt.Println("\nUsage:")
	fmt.Println("  synta-parse [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -input string")
	fmt.Println("        Input token file (default: tokens.json)")
	fmt.Println("  -tree string")
	fmt.Println("        Output parse tree file (default: parse-tree.txt)")
	fmt.Println("  -ast string")
	fmt.Println("        Output AST JSON file (default: ast.json)")
	fmt.Println("  -errors string")
	fmt.Println("        Parse errors file (default: parse-errors.txt)")
	fmt.Println("  -debug string")
	fmt.Println("        Debug log file (default: parse-debug.txt)")
	fmt.Println("  -format string")
	fmt.Println("        Tree format: pretty, compact, or detailed (default: pretty)")
	fmt.Println("  -show")
	fmt.Println("        Show tree in console")
	fmt.Println("  -skip-ast")
	fmt.Println("        Skip AST JSON generation")
	fmt.Println("  -skip-debug")
	fmt.Println("        Skip debug log generation")
	fmt.Println("\nExamples:")
	fmt.Println("  synta-parse")
	fmt.Println("  synta-parse -input my_tokens.json -format compact -show")
	fmt.Println("  synta-parse -skip-ast -skip-debug")
}

func printSummary(program *parser.Program, format string) {
	fmt.Println("\n" + strings.Repeat("-", 70))
	fmt.Println("SUMMARY")
	fmt.Println(strings.Repeat("-", 70))

	// Count statement types
	statementCounts := make(map[string]int)
	for _, stmt := range program.Statements {
		typeName := fmt.Sprintf("%T", stmt)
		// Clean up type name
		typeName = strings.TrimPrefix(typeName, "*parser.")
		statementCounts[typeName]++
	}

	fmt.Printf("Total Statements: %d\n", len(program.Statements))
	fmt.Println("\nStatement Breakdown:")
	for typeName, count := range statementCounts {
		fmt.Printf("  %-25s %d\n", typeName+":", count)
	}

	fmt.Printf("\nTree Format: %s\n", format)
	fmt.Println(strings.Repeat("-", 70))
}
