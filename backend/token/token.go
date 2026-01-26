// token.go
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
	WHILE
	FOR
	SWITCH // Added for switch/case support
	MATCH
	CASE
	DEFAULT
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
	INT_TYPE
	FLOAT_TYPE
	CHAR_TYPE
	BOOL_TYPE
	STR_TYPE
	MAP_TYPE
	ARRAY_TYPE
	TRUE_KW
	FALSE_KW

	// Keywords - Concurrency
	ASYNC
	EMIT
	LISTEN
	DISPATCH
	MERGE
	TASK
	CONCURRENT
	STAGE

	// Keywords - Special Constructs
	WITH
	THEN
	DEFER
	PIPE
	PASS
	THROUGH
	RANGE
	ALLOW
	PSEUDO
	STRATEGY
	TIMEOUT
	WINDOW
	ALERT_THRESHOLD

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
	AT_AGENT   // @agent
	AT_TASK    // @task
	AT_STEP    // @step
	AT_INTENT  // @intent
	AT_EXPLAIN // @explain
	DECORATOR  // @other

	// Keywords - I/O
	READ
	WRITE
	PRINT
	LOG
	SAVE
	FLOW
	CONTEXT
	MEMORY

	// Keywords - Debug
	DEBUG
	CHECKPOINT
	TRACE
	ASSERT
	CONFIGURE
	GENERATE_REPORT

	// Reserved Words
	AGENT
	CORE
	MODEL
	TOOLS
	ROLE
	MODE
	SYS_PROMPT
	MAX_CONCURRENT_REQUESTS
	RETRY_POLICY
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
	CREATE_POOL
	MAX_WORKERS
	SUBMIT
	SUBMIT_DELAYED
	JOIN
	NOW
	EXECUTION_TIME
	REPORT

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
	ASSIGN       // =:
	BIND_ASSIGN  // :=
	PLUS_ASSIGN  // +=
	MINUS_ASSIGN // -=
	MULT_ASSIGN  // *=
	DIV_ASSIGN   // /=
	MOD_ASSIGN   // %=
	EQ           // ==
	NEQ          // !=
	LT           // <
	GT           // >
	LTE          // <=
	GTE          // >=
	AND          // &&
	OR           // ||
	NOT          // !
	BITWISE_XOR  // ^
	AMPERSAND    // &
	INCREMENT    // ++
	DECREMENT    // --
	ARROW        // ->
	FAT_ARROW    // =>
	DOLLAR       // $
	PIPE_OP      // |

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
	IF: "IF", ELIF: "ELIF", ELSE: "ELSE", WHILE: "WHILE", FOR: "FOR",
	SWITCH: "SWITCH", MATCH: "MATCH", CASE: "CASE", DEFAULT: "DEFAULT",
	RETURN: "RETURN", AWAIT: "AWAIT",
	BREAK: "BREAK", CONTINUE: "CONTINUE",
	BIND: "BIND", CONST: "CONST", CRAFT: "CRAFT", USE: "USE", AS: "AS", FROM: "FROM",
	FN: "FN", STRUCT: "STRUCT", TRY: "TRY", CATCH: "CATCH", RAISE: "RAISE",
	TYPE: "TYPE", CAST: "CAST", ANY: "ANY", NONE: "NONE", TRAIT: "TRAIT",
	INT_TYPE: "INT", FLOAT_TYPE: "FLOAT", CHAR_TYPE: "CHAR",
	BOOL_TYPE: "BOOL", STR_TYPE: "STR", MAP_TYPE: "MAP", ARRAY_TYPE: "ARRAY",
	TRUE_KW: "TRUE", FALSE_KW: "FALSE",
	ASYNC: "ASYNC", EMIT: "EMIT", LISTEN: "LISTEN", DISPATCH: "DISPATCH", MERGE: "MERGE",
	TASK: "TASK", CONCURRENT: "CONCURRENT", STAGE: "STAGE",
	WITH: "WITH", THEN: "THEN", DEFER: "DEFER", PIPE: "PIPE", PASS: "PASS",
	THROUGH: "THROUGH", RANGE: "RANGE", ALLOW: "ALLOW", PSEUDO: "PSEUDO",
	STRATEGY: "STRATEGY", TIMEOUT: "TIMEOUT", WINDOW: "WINDOW", ALERT_THRESHOLD: "ALERT_THRESHOLD",
	THINK: "THINK", ASK: "ASK", PROMPT: "PROMPT", ADAPT: "ADAPT", CALL_API: "CALL_API",
	TRAIN: "TRAIN", EVALUATE: "EVALUATE", REASON: "REASON", OBSERVE: "OBSERVE",
	AT_AGENT: "AT_AGENT", AT_TASK: "AT_TASK", AT_STEP: "AT_STEP",
	AT_INTENT: "AT_INTENT", AT_EXPLAIN: "AT_EXPLAIN", DECORATOR: "DECORATOR",
	READ: "READ", WRITE: "WRITE", PRINT: "PRINT", LOG: "LOG", SAVE: "SAVE",
	FLOW: "FLOW", CONTEXT: "CONTEXT", MEMORY: "MEMORY",
	DEBUG: "DEBUG", CHECKPOINT: "CHECKPOINT", TRACE: "TRACE", ASSERT: "ASSERT",
	CONFIGURE: "CONFIGURE", GENERATE_REPORT: "GENERATE_REPORT",
	AGENT: "AGENT", CORE: "CORE", MODEL: "MODEL", TOOLS: "TOOLS", ROLE: "ROLE",
	MODE: "MODE", SYS_PROMPT: "SYS_PROMPT", MAX_CONCURRENT_REQUESTS: "MAX_CONCURRENT_REQUESTS",
	RETRY_POLICY: "RETRY_POLICY", OWN: "OWN", MOVE: "MOVE", DROP: "DROP", LET: "LET",
	PUB: "PUB", PRIV: "PRIV", GLOBAL: "GLOBAL", UNSAFE: "UNSAFE", RAW: "RAW",
	FUTURE: "FUTURE", MACRO: "MACRO", DELEGATE: "DELEGATE", ROUTE: "ROUTE",
	COMPOSE: "COMPOSE", INSPECT: "INSPECT", CREATE_POOL: "CREATE_POOL", MAX_WORKERS: "MAX_WORKERS",
	SUBMIT: "SUBMIT", SUBMIT_DELAYED: "SUBMIT_DELAYED", JOIN: "JOIN",
	NOW: "NOW", EXECUTION_TIME: "EXECUTION_TIME", REPORT: "REPORT",
	DO: "DO", PLEASE: "PLEASE", MAYBE: "MAYBE",
	PLUS: "PLUS", MINUS: "MINUS", MULTIPLY: "MULTIPLY", DIVIDE: "DIVIDE", MODULO: "MODULO",
	ASSIGN: "ASSIGN", BIND_ASSIGN: "BIND_ASSIGN", PLUS_ASSIGN: "PLUS_ASSIGN",
	MINUS_ASSIGN: "MINUS_ASSIGN", MULT_ASSIGN: "MULT_ASSIGN", DIV_ASSIGN: "DIV_ASSIGN",
	MOD_ASSIGN: "MOD_ASSIGN", EQ: "EQ", NEQ: "NEQ", LT: "LT", GT: "GT", LTE: "LTE", GTE: "GTE",
	AND: "AND", OR: "OR", NOT: "NOT", BITWISE_XOR: "BITWISE_XOR", AMPERSAND: "AMPERSAND",
	INCREMENT: "INCREMENT", DECREMENT: "DECREMENT", ARROW: "ARROW", FAT_ARROW: "FAT_ARROW",
	DOLLAR: "DOLLAR", PIPE_OP: "PIPE_OP",
	LPAREN: "LPAREN", RPAREN: "RPAREN", LBRACKET: "LBRACKET", RBRACKET: "RBRACKET",
	LBRACE: "LBRACE", RBRACE: "RBRACE", SEMICOLON: "SEMICOLON", COMMA: "COMMA",
	COLON: "COLON", DOT: "DOT", STATEMENT_END: "STMT_END",
	COMMENT_LINE: "COMMENT_LINE", COMMENT_MULTI: "COMMENT_MULTI",
	NEWLINE: "NEWLINE", EOF: "EOF", ILLEGAL: "ILLEGAL",
}

var Keywords = map[string]TokenType{
	// Control Flow
	"if": IF, "elif": ELIF, "else": ELSE, "while": WHILE, "for": FOR,
	"switch": SWITCH, "match": MATCH, "case": CASE, "default": DEFAULT,
	"return": RETURN, "await": AWAIT,
	"break": BREAK, "continue": CONTINUE,

	// Declarations
	"bind": BIND, "const": CONST, "craft": CRAFT, "use": USE, "as": AS, "from": FROM,
	"fn": FN, "struct": STRUCT,

	// Error Handling
	"try": TRY, "catch": CATCH, "raise": RAISE,

	// Type System
	"type": TYPE, "cast": CAST, "any": ANY, "none": NONE, "trait": TRAIT,
	"int": INT_TYPE, "float": FLOAT_TYPE, "char": CHAR_TYPE, "bool": BOOL_TYPE,
	"str": STR_TYPE, "Map": MAP_TYPE, "true": TRUE_KW, "false": FALSE_KW,

	// Concurrency
	"async": ASYNC, "emit": EMIT, "listen": LISTEN, "dispatch": DISPATCH, "merge": MERGE,
	"task": TASK, "concurrent": CONCURRENT, "stage": STAGE,

	// Special Constructs
	"with": WITH, "then": THEN, "defer": DEFER, "pipe": PIPE, "pass": PASS,
	"through": THROUGH, "range": RANGE, "allow": ALLOW,
	"strategy": STRATEGY, "timeout": TIMEOUT, "window": WINDOW, "alert_threshold": ALERT_THRESHOLD,

	// AI Integration
	"think": THINK, "ask": ASK, "prompt": PROMPT, "adapt": ADAPT, "call_api": CALL_API,
	"train": TRAIN, "evaluate": EVALUATE, "reason": REASON, "observe": OBSERVE,

	// I/O
	"read": READ, "write": WRITE, "print": PRINT, "log": LOG, "save": SAVE,
	"flow": FLOW, "context": CONTEXT, "memory": MEMORY,

	// Debug
	// Note: "debug", "trace", "checkpoint" are NOT keywords
	// They can be used as identifiers (field names and function args)
	"assert":    ASSERT,
	"configure": CONFIGURE, "generate_report": GENERATE_REPORT,

	// Reserved Words
	// Note: "model", "tools", "role", "mode", "sys_prompt", "pseudo" are NOT keywords
	// They can be used as identifiers (field names in @agent blocks)
	"Agent": AGENT, "Core": CORE,
	"max_concurrent_requests": MAX_CONCURRENT_REQUESTS,
	"retry_policy":            RETRY_POLICY, "own": OWN, "move": MOVE, "drop": DROP, "let": LET,
	"pub": PUB, "priv": PRIV, "global": GLOBAL, "unsafe": UNSAFE, "raw": RAW,
	"future": FUTURE, "macro": MACRO, "delegate": DELEGATE, "route": ROUTE,
	"compose": COMPOSE, "inspect": INSPECT, "create_pool": CREATE_POOL,
	"max_workers": MAX_WORKERS, "submit": SUBMIT, "submit_delayed": SUBMIT_DELAYED,
	"join": JOIN, "now": NOW, "execution_time": EXECUTION_TIME, "Report": REPORT,

	// Noise Words
	"do": DO, "please": PLEASE, "maybe": MAYBE,
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

func (t TokenType) String() string {
	if name, ok := TokenNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}
