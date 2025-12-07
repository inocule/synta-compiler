# Synta Lexical Analyzer

> An experimental lexer for an AI-native programming language designed for human–AI collaboration, agentic task execution, and deterministic reasoning.

Synta emphasizes static typing, concurrency primitives, intent-level debugging, and syntax optimized for LLM interpretability. The lexer provides the foundational token stream for Synta's compiler, AI runtime, multi-agent scheduler, and tooling ecosystem.

---

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        SYNTA LEXER SYSTEM                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────┐           ┌──────────────────┐                │
│  │  Frontend (UI)   │  HTTP     │  Backend Server  │                │
│  │   React + TS     │◄─────────►│    Go Runtime    │                │
│  │  localhost:5173  │  JSON     │  localhost:8080  │                │
│  └──────────────────┘           └────────┬─────────┘                │
│           │                               │                         │
│           │                               ▼                         │
│           │                    ┌────────────────────┐               │
│           │                    │  Lexer Core        │               │
│           │                    │  (lexer/lexer.go)  │               │
│           │                    └────────┬───────────┘               │
│           │                             │                           │
│           │                             ▼                           │
│           │                    ┌────────────────────┐               │
│           │                    │  Token Generator   │               │
│           │                    │  (token/token.go)  │               │
│           │                    └────────┬───────────┘               │
│           │                             │                           │
│           ▼                             ▼                           │
│  ┌─────────────────────────────────────────────────┐                │
│  │             TOKEN STREAM OUTPUT                 │                │
│  │  [KEYWORD, IDENT, OPERATOR, NUMBER, STRING...]  │                │
│  └─────────────────────────────────────────────────┘                │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  PARSER PIPELINE (Future):                                          │
│  Token Stream → AST Builder → Type Checker → Agent Scheduler        │
└─────────────────────────────────────────────────────────────────────┘

DATA FLOW:
  Source Code (.synta) 
      ↓
  Lexer Tokenization
      ↓
  Token Metadata (line, col, type)
      ↓
  Parser / AI Runtime
```

---

## ✨ Core Features

### Lexical Analysis
- **Tokens**: Identifiers, numbers (int/float), strings, operators, delimiters
- **Comments**: Single-line `!>` and multi-line `<! ... !>` (unified token capture)
- **Delimiters**: Statement-end token `~` for clear boundaries
- **Tracking**: Precise line/column position for every token
- **Newlines**: Optional newline tokens for structure-aware parsers

### AI-Native Syntax
- **Keywords**: `think`, `reason`, `ask`, `observe`, `intent`, `pipeline`, `agent`, `concurrent`, `sequential`
- **Decorators**: `@agent`, `@task`, `@model`, `@pipeline` for agent definitions
- **Invocation**: Arrow operator `->` for type-safe AI calls
- **Types**: `int`, `float`, `bool`, `char`, `str`, `object`, `list<T>`
- **Operators**: 
  - `:=` for binding (immutable declaration)
  - `=:` for assignment (mutable update)

---

## 📁 Project Structure

```
synta-lexical/
├── backend-server/
│   └── main.go              # Go HTTP server
├── frontend/
│   ├── src/
│   │   ├── Components/
│   │   │   ├── EditorPane.tsx      # Code input (left)
│   │   │   └── OutputTable.tsx     # Token output (right)
│   │   ├── App.tsx
│   │   ├── main.tsx
│   │   └── types.ts
│   ├── package.json
│   └── vite.config.ts
├── lexer/
│   ├── lexer.go             # Core tokenizer
├── token/
│   └── token.go             # Token types & keywords
├── go.mod
└── README.md
```

---

## 🚀 Quick Start

### Prerequisites
```bash
go version    # Requires Go 1.21+
node --version # Requires Node.js 16+
```

### Installation & Running

**1. Start Backend Server**
```bash
cd synta-lexical/cmd/server
go run main.go
# Server running at http://localhost:8080
```

**2. Start Frontend (New Terminal)**
```bash
cd synta-lexical/frontend
npm install           # First time only
npm run dev
# UI available at http://localhost:5173
```

**3. Open Browser**
Navigate to `http://localhost:5173` and start tokenizing!

---

## 📝 Language Syntax Examples

### Agent Definition
```synta
@agent AICoder {
    role: "GitHub-integrated coding assistant",
    tools: [github_mcp, slm_chatbot, pdf_scanner],
    model: "local/llama-3.1-8b.gguf",
    mode: "hybrid",
    sys_prompt: "Up-to-date assistant for repositories"
}
```

### Task Execution
```synta
@task {
    response:str =: AICoder -> "Fix syntax errors"
    print(response)
}
```

### Variables & Types
```synta
!> Binding (immutable)
bind x:int := 10;
const PI := 3.14;

!> Assignment (mutable)
x =: 20;
```

### Functions
```synta
fn calculate(a:int, b:int) -> int do {
    if a > b {
        return a + b
    } else {
        return a - b
    }
}
```

### Intent Blocks
```synta
intent {
    goal: "Analyze Q4 sales trends";
    context: "Processed dataset from pipeline";
    reason: "Preprocessing completed successfully";
}
```

### Concurrency
```synta
for task in tasks concurrent {
    process(task);
}

pipeline analysis_flow {
    start analyze_dataset;
    then generate_report concurrent;
}
```

### Model Fine-Tuning
```synta
model, tokenizer =: unsloth.FastLanguageModel.from_pretrained(
    model_name: "unsloth/Phi-4-mini-instruct",
    max_seq_length: 2048,
    dtype: None,
    load_in_4bit: True
);
```

---

## 🎯 Language Design Principles

### 1. Static & Strong Typing
- Variables use `name:type` annotation
- Assignment `=:` vs binding `:=` distinction
- Compile-time type checking for AI invocations

### 2. Intent-Level Debugging
- First-class `intent` blocks for AI chain-of-thought
- Context window replenishment support
- Deterministic introspection for agentic runtimes

### 3. Concurrency as Core Syntax
- `concurrent`, `parallel`, `sequential` as reserved words
- Runtime scheduler integration
- DAG-based task execution graphs

### 4. AI-Friendly Parsing
- Low-ambiguity operators (`=:`, `:=`, `->`)
- Mandatory braces for blocks
- Explicit type annotations
- Machine-readable pipeline syntax

### 5. Memory Safety
- Scoped lifetimes with `own` and `borrow` semantics
- Deterministic cleanup events
- Built-in garbage collection

### 6. Native LLM/SLM Integration
- Vector database identifiers
- RAG pipeline keywords
- RLHF-loop verbs
- Embedded model descriptors

### 7. Agentic Extensions
- Multi-agent orchestration keywords
- `depends_on`, `emit`, `listen` primitives
- Deterministic execution DAGs

---

## 🔧 Extending the Lexer

### Add Keywords
```go
// In token/token.go
const (
    NEWKEYWORD = "NEWKEYWORD"
)

var Keywords = map[string]TokenType{
    "newkeyword": NEWKEYWORD,
}
```

### Add Operators
```go
// In lexer/lexer.go Tokenize() switch
case '★':
    tokens = append(tokens, Token{Type: STAR_OP, Literal: "★"})
```

### Add Decorators
```go
// In lexer/lexer.go @ handling
case '@':
    if isLetter(peekChar()) {
        decorator := readDecorator()
        // Add new decorator logic
    }
```

---

## 📊 Token Types Reference

| Category | Examples |
|----------|----------|
| **Keywords** | `think`, `reason`, `intent`, `pipeline`, `concurrent` |
| **Decorators** | `@agent`, `@task`, `@model`, `@pipeline` |
| **Operators** | `:=` (bind), `=:` (assign), `->` (invoke) |
| **Types** | `int`, `float`, `bool`, `str`, `list<T>` |
| **Delimiters** | `~` (statement end), `{`, `}`, `(`, `)` |
| **Comments** | `!>` (single), `<! !>` (multi) |

---

## 🧪 Testing

```bash
cd lexer
go test -v
```

---

## 📜 License

Experimental - Educational Use

---

## 🤝 Contributing

This is an experimental language design. Contributions welcome for:
- Additional AI-native keywords
- Enhanced token metadata
- Performance optimizations
- Parser integration examples

---

**Built for the future of human–AI collaborative programming** 🚀