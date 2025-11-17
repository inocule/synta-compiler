// token/token.go
package token

type TokenType int

const (
	// Literals
	IDENTIFIER TokenType = iota
	INTEGER
	FLOAT
	STRING

	// Keywords - Control Flow
	IF
	ELIF
	ELSE
	FOR
	WHILE
	MATCH
	RETURN
	AWAIT
	BREAK
	CONTINUE

	// Keywords - Declarations
	BIND
	CONST
	CRAFT
	USE
	AS
	FROM
	FN
	STRUCT

	// Keywords - Error Handling
	TRY
	CATCH
	RAISE

	// Keywords - Type System
	TYPE
	CAST
	ANY
	NONE
	TRAIT

	// Keywords - Concurrency
	ASYNC
	EMIT
	LISTEN
	DISPATCH
	MERGE

	// Keywords - Special Constructs
	WITH
	THEN
	DEFER
	PIPE
	PASS

	// Keywords - AI Integration
	THINK
	ASK
	PROMPT
	ADAPT
	CALL_API
	TRAIN
	EVALUATE
	REASON
	OBSERVE

	// Keywords - I/O
	READ
	WRITE
	PRINT
	LOG
	SAVE
	FLOW
	CONTEXT
	MEMORY

	// Reserved Words
	AGENT
	CORE
	MODEL
	TOOLS
	ROLE
	MODE
	OWN
	MOVE
	DROP
	LET
	PUB
	PRIV
	GLOBAL
	UNSAFE
	RAW
	FUTURE
	MACRO
	DELEGATE
	ROUTE
	COMPOSE
	INSPECT

	// Noise Words
	PLEASE
	MAYBE
	DO

	// Operators
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	MODULO
	ASSIGN
	BIND_ASSIGN
	PLUS_ASSIGN
	MINUS_ASSIGN
	MULT_ASSIGN
	DIV_ASSIGN
	MOD_ASSIGN
	EQ
	NEQ
	LT
	GT
	LTE
	GTE
	AND
	OR
	NOT
	BITWISE_XOR
	AMPERSAND
	INCREMENT
	DECREMENT

	// Delimiters
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	SEMICOLON
	COMMA
	COLON
	DOT
	ARROW

	// Special
	COMMENT
	NEWLINE
	EOF
	ILLEGAL
)

var TokenNames = map[TokenType]string{
	IDENTIFIER: "IDENTIFIER", INTEGER: "INTEGER", FLOAT: "FLOAT", STRING: "STRING",
	IF: "IF", ELIF: "ELIF", ELSE: "ELSE", FOR: "FOR", WHILE: "WHILE",
	MATCH: "MATCH", RETURN: "RETURN", AWAIT: "AWAIT", BREAK: "BREAK", CONTINUE: "CONTINUE",
	BIND: "BIND", CONST: "CONST", CRAFT: "CRAFT", USE: "USE", AS: "AS", FROM: "FROM",
	FN: "FN", STRUCT: "STRUCT", TRY: "TRY", CATCH: "CATCH", RAISE: "RAISE",
	TYPE: "TYPE", CAST: "CAST", ANY: "ANY", NONE: "NONE", TRAIT: "TRAIT",
	ASYNC: "ASYNC", EMIT: "EMIT", LISTEN: "LISTEN", DISPATCH: "DISPATCH", MERGE: "MERGE",
	WITH: "WITH", THEN: "THEN", DEFER: "DEFER", PIPE: "PIPE", PASS: "PASS",
	THINK: "THINK", ASK: "ASK", PROMPT: "PROMPT", ADAPT: "ADAPT", CALL_API: "CALL_API",
	TRAIN: "TRAIN", EVALUATE: "EVALUATE", REASON: "REASON", OBSERVE: "OBSERVE",
	READ: "READ", WRITE: "WRITE", PRINT: "PRINT", LOG: "LOG", SAVE: "SAVE",
	FLOW: "FLOW", CONTEXT: "CONTEXT", MEMORY: "MEMORY",
	AGENT: "AGENT", CORE: "CORE", MODEL: "MODEL", TOOLS: "TOOLS", ROLE: "ROLE",
	MODE: "MODE", OWN: "OWN", MOVE: "MOVE", DROP: "DROP", LET: "LET",
	PUB: "PUB", PRIV: "PRIV", GLOBAL: "GLOBAL", UNSAFE: "UNSAFE", RAW: "RAW",
	FUTURE: "FUTURE", MACRO: "MACRO", DELEGATE: "DELEGATE", ROUTE: "ROUTE",
	COMPOSE: "COMPOSE", INSPECT: "INSPECT", DO: "DO", PLEASE: "PLEASE", MAYBE: "MAYBE",
	PLUS: "PLUS", MINUS: "MINUS", MULTIPLY: "MULTIPLY", DIVIDE: "DIVIDE", MODULO: "MODULO",
	ASSIGN: "ASSIGN", BIND_ASSIGN: "BIND_ASSIGN", PLUS_ASSIGN: "PLUS_ASSIGN",
	MINUS_ASSIGN: "MINUS_ASSIGN", MULT_ASSIGN: "MULT_ASSIGN", DIV_ASSIGN: "DIV_ASSIGN",
	MOD_ASSIGN: "MOD_ASSIGN", EQ: "EQ", NEQ: "NEQ", LT: "LT", GT: "GT", LTE: "LTE", GTE: "GTE",
	AND: "AND", OR: "OR", NOT: "NOT", BITWISE_XOR: "BITWISE_XOR", AMPERSAND: "AMPERSAND",
	INCREMENT: "INCREMENT", DECREMENT: "DECREMENT",
	LPAREN: "LPAREN", RPAREN: "RPAREN", LBRACKET: "LBRACKET", RBRACKET: "RBRACKET",
	LBRACE: "LBRACE", RBRACE: "RBRACE", SEMICOLON: "SEMICOLON", COMMA: "COMMA",
	COLON: "COLON", DOT: "DOT", ARROW: "ARROW",
	COMMENT: "COMMENT", NEWLINE: "NEWLINE", EOF: "EOF", ILLEGAL: "ILLEGAL",
}

var Keywords = map[string]TokenType{
	"if": IF, "elif": ELIF, "else": ELSE, "for": FOR, "while": WHILE,
	"match": MATCH, "return": RETURN, "await": AWAIT, "break": BREAK, "continue": CONTINUE,
	"bind": BIND, "const": CONST, "craft": CRAFT, "use": USE, "as": AS, "from": FROM,
	"fn": FN, "struct": STRUCT, "try": TRY, "catch": CATCH, "raise": RAISE,
	"type": TYPE, "cast": CAST, "any": ANY, "none": NONE, "trait": TRAIT,
	"async": ASYNC, "emit": EMIT, "listen": LISTEN, "dispatch": DISPATCH, "merge": MERGE,
	"with": WITH, "then": THEN, "defer": DEFER, "pipe": PIPE, "pass": PASS,
	"think": THINK, "ask": ASK, "prompt": PROMPT, "adapt": ADAPT, "call_api": CALL_API,
	"train": TRAIN, "evaluate": EVALUATE, "reason": REASON, "observe": OBSERVE,
	"read": READ, "write": WRITE, "print": PRINT, "log": LOG, "save": SAVE,
	"flow": FLOW, "context": CONTEXT, "memory": MEMORY,
	"Agent": AGENT, "Core": CORE, "model": MODEL, "tools": TOOLS, "role": ROLE,
	"mode": MODE, "own": OWN, "move": MOVE, "drop": DROP, "let": LET,
	"pub": PUB, "priv": PRIV, "global": GLOBAL, "unsafe": UNSAFE, "raw": RAW,
	"future": FUTURE, "macro": MACRO, "delegate": DELEGATE, "route": ROUTE,
	"compose": COMPOSE, "inspect": INSPECT, "do": DO, "please": PLEASE, "maybe": MAYBE,
}

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
}

func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
