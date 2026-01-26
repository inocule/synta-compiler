// token.go
package token

type TokenType int

const (
	// Literals & Identifiers
	IDENTIFIER TokenType = iota
	INTEGER
	FLOAT
	STRING

	// Control Flow
	IF
	ELIF
	ELSE
	WHILE
	FOR
	SWITCH
	MATCH
	CASE
	DEFAULT
	RETURN
	BREAK
	CONTINUE

	// Declarations
	BIND
	CONST
	CRAFT
	USE
	AS
	FROM
	FN
	STRUCT

	// Error Handling
	TRY
	CATCH
	RAISE

	// Type System
	TYPE
	CAST
	ANY
	NONE
	TRAIT
	INT_TYPE
	FLOAT_TYPE
	CHAR_TYPE
	BOOL_TYPE
	STR_TYPE
	MAP_TYPE
	ARRAY_TYPE
	TRUE
	FALSE
	NULL

	// Concurrency & Async
	ASYNC
	AWAIT
	EMIT
	LISTEN
	DISPATCH
	MERGE
	TASK
	CONCURRENT
	STAGE
	GATHER

	// Special Constructs
	WITH
	THEN
	DEFER
	PIPE
	PASS
	ALLOW
	THROUGH
	RANGE
	STRATEGY
	TIMEOUT
	WINDOW
	ALERT_THRESHOLD
	PSEUDO

	// AI Integration
	THINK
	ASK
	PROMPT
	ADAPT
	CALL_API
	TRAIN
	EVALUATE
	REASON
	OBSERVE

	// Decorators
	AT_AGENT   // @agent
	AT_TASK    // @task
	AT_STEP    // @step
	AT_INTENT  // @intent
	AT_EXPLAIN // @explain
	DECORATOR  // @other

	// I/O & Debug
	READ
	WRITE
	PRINT
	LOG
	SAVE
	FLOW
	CONTEXT
	MEMORY
	DEBUG
	CHECKPOINT
	BREAKPOINT
	TRACE
	ASSERT
	CONFIGURE
	GENERATE_REPORT

	// Agent/System Reserved
	AGENT
	CORE
	MODEL
	TOOLS
	ROLE
	MODE
	SYS_PROMPT
	MAX_CONCURRENT_REQUESTS
	RETRY_POLICY

	// Ownership & Memory (reserved for future use)
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

	// Agent Operations (reserved)
	DELEGATE
	ROUTE
	COMPOSE
	INSPECT
	CREATE_POOL
	MAX_WORKERS
	SUBMIT
	SUBMIT_DELAYED
	JOIN
	NOW
	EXECUTION_TIME
	REPORT

	// Noise Words (optional syntax sugar)
	PLEASE
	MAYBE
	DO

	// Loop/Watch Keywords
	LOOP
	GUARD
	WATCH
	ON
	SNAPSHOT
	RESTORE

	// Configuration Keywords
	INPUT
	ACTION
	EXECUTION
	RETRY
	ENABLED
	MAX
	DEPENDS_ON
	CONFIG
	OUTPUTS
	BREAKPOINTS
	ON_CONCUR_DEADLOCK
	ON_LOOP
	ON_TIMEOUT

	// Operators - Arithmetic
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	MODULO   // %

	// Operators - Assignment
	ASSIGN       // =:
	BIND_ASSIGN  // :=
	PLUS_ASSIGN  // +=
	MINUS_ASSIGN // -=
	MULT_ASSIGN  // *=
	DIV_ASSIGN   // /=
	MOD_ASSIGN   // %=

	// Operators - Comparison
	EQ  // ==
	NEQ // !=
	LT  // <
	GT  // >
	LTE // <=
	GTE // >=

	// Operators - Logical
	AND // &&
	OR  // ||
	NOT // !

	// Operators - Bitwise & Other
	BITWISE_XOR   // ^
	AMPERSAND     // &
	INCREMENT     // ++
	DECREMENT     // --
	ARROW         // ->
	FAT_ARROW     // =>
	DOLLAR        // $
	PIPE_OP       // |
	DOUBLE_COLON  // ::
	PIPE_RIGHT    // |>
	PIPE_PARALLEL // |>>

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
	STATEMENT_END

	// Special
	COMMENT_LINE
	COMMENT_MULTI
	NEWLINE
	EOF
	ILLEGAL
)

var TokenNames = map[TokenType]string{
	IDENTIFIER: "IDENTIFIER", INTEGER: "INTEGER", FLOAT: "FLOAT", STRING: "STRING",
	IF: "IF", ELIF: "ELIF", ELSE: "ELSE", WHILE: "WHILE", FOR: "FOR", SWITCH: "SWITCH", MATCH: "MATCH",
	CASE: "CASE", DEFAULT: "DEFAULT", RETURN: "RETURN", BREAK: "BREAK", CONTINUE: "CONTINUE",
	BIND: "BIND", CONST: "CONST", CRAFT: "CRAFT", USE: "USE", AS: "AS", FROM: "FROM", FN: "FN", STRUCT: "STRUCT",
	TRY: "TRY", CATCH: "CATCH", RAISE: "RAISE",
	TYPE: "TYPE", CAST: "CAST", ANY: "ANY", NONE: "NONE", TRAIT: "TRAIT",
	INT_TYPE: "INT", FLOAT_TYPE: "FLOAT", CHAR_TYPE: "CHAR", BOOL_TYPE: "BOOL", STR_TYPE: "STR", MAP_TYPE: "MAP", ARRAY_TYPE: "ARRAY",
	TRUE: "TRUE", FALSE: "FALSE", NULL: "NULL",
	ASYNC: "ASYNC", AWAIT: "AWAIT", EMIT: "EMIT", LISTEN: "LISTEN", DISPATCH: "DISPATCH", MERGE: "MERGE", TASK: "TASK",
	CONCURRENT: "CONCURRENT", STAGE: "STAGE", GATHER: "GATHER",
	WITH: "WITH", THEN: "THEN", DEFER: "DEFER", PIPE: "PIPE", PASS: "PASS", ALLOW: "ALLOW", THROUGH: "THROUGH", RANGE: "RANGE",
	STRATEGY: "STRATEGY", TIMEOUT: "TIMEOUT", WINDOW: "WINDOW", ALERT_THRESHOLD: "ALERT_THRESHOLD", PSEUDO: "PSEUDO",
	THINK: "THINK", ASK: "ASK", PROMPT: "PROMPT", ADAPT: "ADAPT", CALL_API: "CALL_API", TRAIN: "TRAIN", EVALUATE: "EVALUATE",
	REASON: "REASON", OBSERVE: "OBSERVE",
	AT_AGENT: "AT_AGENT", AT_TASK: "AT_TASK", AT_STEP: "AT_STEP", AT_INTENT: "AT_INTENT", AT_EXPLAIN: "AT_EXPLAIN", DECORATOR: "DECORATOR",
	READ: "READ", WRITE: "WRITE", PRINT: "PRINT", LOG: "LOG", SAVE: "SAVE", FLOW: "FLOW", CONTEXT: "CONTEXT", MEMORY: "MEMORY",
	DEBUG: "DEBUG", CHECKPOINT: "CHECKPOINT", BREAKPOINT: "BREAKPOINT", TRACE: "TRACE", ASSERT: "ASSERT", CONFIGURE: "CONFIGURE", GENERATE_REPORT: "GENERATE_REPORT",
	AGENT: "AGENT", CORE: "CORE", MODEL: "MODEL", TOOLS: "TOOLS", ROLE: "ROLE", MODE: "MODE", SYS_PROMPT: "SYS_PROMPT",
	MAX_CONCURRENT_REQUESTS: "MAX_CONCURRENT_REQUESTS", RETRY_POLICY: "RETRY_POLICY",
	OWN: "OWN", MOVE: "MOVE", DROP: "DROP", LET: "LET", PUB: "PUB", PRIV: "PRIV", GLOBAL: "GLOBAL",
	UNSAFE: "UNSAFE", RAW: "RAW", FUTURE: "FUTURE", MACRO: "MACRO",
	DELEGATE: "DELEGATE", ROUTE: "ROUTE", COMPOSE: "COMPOSE", INSPECT: "INSPECT", CREATE_POOL: "CREATE_POOL", MAX_WORKERS: "MAX_WORKERS",
	SUBMIT: "SUBMIT", SUBMIT_DELAYED: "SUBMIT_DELAYED", JOIN: "JOIN", NOW: "NOW", EXECUTION_TIME: "EXECUTION_TIME", REPORT: "REPORT",
	DO: "DO", PLEASE: "PLEASE", MAYBE: "MAYBE",
	LOOP: "LOOP", GUARD: "GUARD", WATCH: "WATCH", ON: "ON", SNAPSHOT: "SNAPSHOT", RESTORE: "RESTORE",
	INPUT: "INPUT", ACTION: "ACTION", EXECUTION: "EXECUTION", RETRY: "RETRY", ENABLED: "ENABLED", MAX: "MAX", DEPENDS_ON: "DEPENDS_ON",
	CONFIG: "CONFIG", OUTPUTS: "OUTPUTS", BREAKPOINTS: "BREAKPOINTS", ON_CONCUR_DEADLOCK: "ON_CONCUR_DEADLOCK", ON_LOOP: "ON_LOOP", ON_TIMEOUT: "ON_TIMEOUT",
	PLUS: "PLUS", MINUS: "MINUS", MULTIPLY: "MULTIPLY", DIVIDE: "DIVIDE", MODULO: "MODULO",
	ASSIGN: "ASSIGN", BIND_ASSIGN: "BIND_ASSIGN", PLUS_ASSIGN: "PLUS_ASSIGN", MINUS_ASSIGN: "MINUS_ASSIGN",
	MULT_ASSIGN: "MULT_ASSIGN", DIV_ASSIGN: "DIV_ASSIGN", MOD_ASSIGN: "MOD_ASSIGN",
	EQ: "EQ", NEQ: "NEQ", LT: "LT", GT: "GT", LTE: "LTE", GTE: "GTE", AND: "AND", OR: "OR", NOT: "NOT",
	BITWISE_XOR: "BITWISE_XOR", AMPERSAND: "AMPERSAND", INCREMENT: "INCREMENT", DECREMENT: "DECREMENT",
	ARROW: "ARROW", FAT_ARROW: "FAT_ARROW", DOLLAR: "DOLLAR", PIPE_OP: "PIPE_OP", DOUBLE_COLON: "DOUBLE_COLON", PIPE_RIGHT: "PIPE_RIGHT", PIPE_PARALLEL: "PIPE_PARALLEL",
	LPAREN: "LPAREN", RPAREN: "RPAREN", LBRACKET: "LBRACKET", RBRACKET: "RBRACKET", LBRACE: "LBRACE", RBRACE: "RBRACE",
	SEMICOLON: "SEMICOLON", COMMA: "COMMA", COLON: "COLON", DOT: "DOT", STATEMENT_END: "STMT_END",
	COMMENT_LINE: "COMMENT_LINE", COMMENT_MULTI: "COMMENT_MULTI", NEWLINE: "NEWLINE", EOF: "EOF", ILLEGAL: "ILLEGAL",
}

var Keywords = map[string]TokenType{
	// Literals
	"true": TRUE, "false": FALSE, "null": NULL,
	// Control Flow
	"if": IF, "elif": ELIF, "else": ELSE, "while": WHILE, "for": FOR, "switch": SWITCH, "match": MATCH, "case": CASE, "default": DEFAULT,
	"return": RETURN, "break": BREAK, "continue": CONTINUE,
	// Declarations
	"bind": BIND, "const": CONST, "craft": CRAFT, "use": USE, "as": AS, "from": FROM, "fn": FN, "struct": STRUCT,
	// Error Handling
	"try": TRY, "catch": CATCH, "raise": RAISE,
	// Type System
	"type": TYPE, "cast": CAST, "any": ANY, "none": NONE, "trait": TRAIT,
	"int": INT_TYPE, "float": FLOAT_TYPE, "char": CHAR_TYPE, "bool": BOOL_TYPE, "str": STR_TYPE, "Map": MAP_TYPE,
	// Async/Concurrency
	"async": ASYNC, "await": AWAIT, "emit": EMIT, "listen": LISTEN, "dispatch": DISPATCH, "merge": MERGE,
	"task": TASK, "concurrent": CONCURRENT, "stage": STAGE, "gather": GATHER,
	// Intent & Debug
	"allow": ALLOW, "pseudo": PSEUDO,
	// IO & Debug
	"read": READ, "write": WRITE, "print": PRINT, "log": LOG, "save": SAVE, "flow": FLOW, "context": CONTEXT, "memory": MEMORY,
	"debug": DEBUG, "checkpoint": CHECKPOINT, "breakpoint": BREAKPOINT, "trace": TRACE, "assert": ASSERT, "configure": CONFIGURE, "generate_report": GENERATE_REPORT,
	// AI Integration
	"think": THINK, "ask": ASK, "prompt": PROMPT, "adapt": ADAPT, "call_api": CALL_API, "train": TRAIN, "evaluate": EVALUATE, "reason": REASON, "observe": OBSERVE,
	// Special Constructs
	"with": WITH, "then": THEN, "defer": DEFER, "pipe": PIPE, "pass": PASS, "through": THROUGH, "range": RANGE,
	"strategy": STRATEGY, "timeout": TIMEOUT, "window": WINDOW, "alert_threshold": ALERT_THRESHOLD,
	// Ownership & Memory
	"own": OWN, "move": MOVE, "drop": DROP, "let": LET, "pub": PUB, "priv": PRIV, "global": GLOBAL,
	"unsafe": UNSAFE, "raw": RAW, "future": FUTURE, "macro": MACRO,
	// Agent Operations
	"delegate": DELEGATE, "route": ROUTE, "compose": COMPOSE, "inspect": INSPECT, "create_pool": CREATE_POOL, "max_workers": MAX_WORKERS,
	"submit": SUBMIT, "submit_delayed": SUBMIT_DELAYED, "join": JOIN, "now": NOW, "execution_time": EXECUTION_TIME, "Report": REPORT,
	// Noise Words
	"do": DO, "please": PLEASE, "maybe": MAYBE,
	// Loop/Watch Keywords
	"loop": LOOP, "guard": GUARD, "watch": WATCH, "on": ON, "snapshot": SNAPSHOT, "restore": RESTORE,
	// Configuration Keywords
	"input": INPUT, "action": ACTION, "execution": EXECUTION, "retry": RETRY, "enabled": ENABLED, "max": MAX, "depends_on": DEPENDS_ON, "config": CONFIG,
	"outputs": OUTPUTS, "breakpoints": BREAKPOINTS, "on_concurrency_deadlock": ON_CONCUR_DEADLOCK, "on_loop": ON_LOOP, "on_timeout": ON_TIMEOUT,
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

// Context-sensitive keywords: Agent/System Reserved
// These should only be recognized as keywords within agent/task block contexts
var AgentSystemKeywords = map[string]TokenType{
	"Agent":                   AGENT,
	"Core":                    CORE,
	"model":                   MODEL,
	"mode":                    MODE,
	"role":                    ROLE,
	"tools":                   TOOLS,
	"sys_prompt":              SYS_PROMPT,
	"max_concurrent_requests": MAX_CONCURRENT_REQUESTS,
	"retry_policy":            RETRY_POLICY,
}

func IsAgentSystemKeyword(ident string) (TokenType, bool) {
	tok, ok := AgentSystemKeywords[ident]
	return tok, ok
}

func (t TokenType) String() string {
	if name, ok := TokenNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}

// Semantic Groups: Map raw token types to high-level semantic categories
var SemanticGroups = map[TokenType]string{
	// Literals
	IDENTIFIER: "IDENTIFIER", INTEGER: "INTEGER", FLOAT: "FLOAT", STRING: "STRING",
	TRUE: "TRUE", FALSE: "FALSE", NULL: "NULL",
	// Control Flow
	IF: "IF", ELIF: "IF", ELSE: "IF", FOR: "FOR", WHILE: "WHILE", RETURN: "RETURN",
	TRY: "TRY", CATCH: "CATCH", BREAK: "BREAK", CONTINUE: "BREAK",
	// Declarations
	FN: "FN", STRUCT: "STRUCT", BIND: "BIND", CONST: "CONST", CRAFT: "CONST",
	// Async/Concurrency
	ASYNC: "ASYNC", AWAIT: "AWAIT", GATHER: "ASYNC", EMIT: "ASYNC", LISTEN: "ASYNC", TASK: "ASYNC", CONCURRENT: "ASYNC",
	// Intent & Debug
	ALLOW: "INTENT_DEBUG", PSEUDO: "INTENT_DEBUG",
	// IO & Debug
	TRACE: "IO_DEBUG", BREAKPOINT: "IO_DEBUG", DEBUG: "IO_DEBUG", CHECKPOINT: "IO_DEBUG", LOG: "IO_DEBUG",
	PRINT: "IO_DEBUG", ASSERT: "IO_DEBUG", CONFIGURE: "IO_DEBUG", GENERATE_REPORT: "IO_DEBUG",
	// Decorators
	AT_AGENT: "DECORATOR", AT_TASK: "DECORATOR", AT_STEP: "DECORATOR", AT_INTENT: "DECORATOR", AT_EXPLAIN: "DECORATOR", DECORATOR: "DECORATOR",
	// Agent/System
	AGENT: "AGENT_SYSTEM", CORE: "AGENT_SYSTEM", MODEL: "AGENT_SYSTEM", TOOLS: "AGENT_SYSTEM", ROLE: "AGENT_SYSTEM", MODE: "AGENT_SYSTEM", SYS_PROMPT: "AGENT_SYSTEM",
	// Loop/Watch
	LOOP: "LOOP", GUARD: "GUARD", WATCH: "WATCH", ON: "ON", SNAPSHOT: "SNAPSHOT", RESTORE: "RESTORE",
	// Configuration
	INPUT: "CONFIG", ACTION: "CONFIG", EXECUTION: "CONFIG", RETRY: "CONFIG", ENABLED: "CONFIG", MAX: "CONFIG", DEPENDS_ON: "CONFIG",
	CONFIG: "CONFIG", OUTPUTS: "CONFIG", BREAKPOINTS: "CONFIG", ON_CONCUR_DEADLOCK: "CONFIG", ON_LOOP: "CONFIG", ON_TIMEOUT: "CONFIG",
	// Operators
	ASSIGN: "ASSIGN", BIND_ASSIGN: "ASSIGN", PLUS_ASSIGN: "ASSIGN", MINUS_ASSIGN: "ASSIGN", MULT_ASSIGN: "ASSIGN", DIV_ASSIGN: "ASSIGN", MOD_ASSIGN: "ASSIGN",
	FAT_ARROW: "FAT_ARROW", EQ: "EQ", NEQ: "EQ", LT: "LT", GT: "GT", LTE: "LT", GTE: "GT",
	PLUS: "OPERATOR", MINUS: "OPERATOR", MULTIPLY: "OPERATOR", DIVIDE: "OPERATOR", MODULO: "OPERATOR", AND: "OPERATOR", OR: "OPERATOR", NOT: "OPERATOR",
	// Delimiters
	LPAREN: "LPAREN", RPAREN: "RPAREN", LBRACE: "LBRACE", RBRACE: "RBRACE", LBRACKET: "LBRACKET", RBRACKET: "RBRACKET", COMMA: "COMMA", COLON: "COLON", DOT: "DOT",
	// Special
	STATEMENT_END: "STMT_END", COMMENT_LINE: "COMMENT_LINE", COMMENT_MULTI: "COMMENT_MULTI", EOF: "EOF", ILLEGAL: "ILLEGAL",
}

// GetSemanticGroup returns the semantic group label for a token type
func (t TokenType) GetSemanticGroup() string {
	if group, ok := SemanticGroups[t]; ok {
		return group
	}
	return "UNKNOWN"
}
