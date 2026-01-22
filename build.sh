#!/bin/bash
# build.sh - Build all Synta compiler tools

set -e

echo "🔨 Building Synta Compiler Tools..."
echo ""

# Create bin directory
mkdir -p bin

# Build synta-parse
echo "Building synta-parse..."
cd synta-parse
go build -o ../bin/synta-parse
cd ..
echo "✓ synta-parse built"

# Build synta-debug
echo "Building synta-debug..."
cd synta-debug
go build -o ../bin/synta-debug
cd ..
echo "✓ synta-debug built"

# Build synta-tree
echo "Building synta-tree..."
cd synta-tree
go build -o ../bin/synta-tree
cd ..
echo "✓ synta-tree built"

echo ""
echo "✅ All tools built successfully!"
echo ""
echo "Tools available in ./bin/:"
echo "  - synta-parse  : Parse tokens and generate AST"
echo "  - synta-debug  : Debug parser with detailed logs"
echo "  - synta-tree   : Generate visual parse trees"
echo ""
echo "Example usage:"
echo "  ./bin/synta-parse -input tokens.json"
echo "  ./bin/synta-debug -input tokens.json -v"
echo "  ./bin/synta-tree -input tokens.json -show"