// parser.go
// syntax analyzer

package parser

import (
	"fmt"
	"synta-compiler/lexer"
	"synta-compiler/token"
)

// AST Node types
type Node interface {
	node()
	String() string
}

type Program struct {
	Statements []Node
}

type Statement struct {
	Type  string
	Value interface{}
}

type Expression struct {
	Type  string
	Value interface{}
}

type Declaration struct {
	DeclType  string // "bind", "const", "craft", "fn", "struct", "agent", "task"
	Name      string
	Type      string
	Value     Node
	Params    []Parameter
	Body      []Node
	Fields    []Field // for struct and agent
	Decorator string  // for functions (@agent, @task, etc.)
}

type Parameter struct {
	Name string
	Type string
}

type Field struct {
	Name  string
	Type  string
	Value Node
}

type BinaryOp struct {
	Left     Node
	Operator string
	Right    Node
}

type UnaryOp struct {
	Operator string
	Operand  Node
}

type Literal struct {
	Type  string // "int", "float", "string", "bool"
	Value string
}

type Identifier struct {
	Name string
}

type IfStatement struct {
	Condition Node
	ThenBody  []Node
	ElseBody  []Node
}

type WhileStatement struct {
	Condition Node
	Body      []Node
}

type ForStatement struct {
	Variable   string
	Init       Node
	Condition  Node
	Update     Node
	Iterable   Node
	Body       []Node
	Concurrent bool
}

type SwitchStatement struct {
	Expression Node
	Cases      []CaseClause
	Default    []Node
}

type CaseClause struct {
	Value Node
	Body  []Node
}

type ReturnStatement struct {
	Value Node
}

type CallExpression struct {
	Function  Node
	Arguments []Node
}

type ArrayLiteral struct {
	Elements []Node
}

type MapLiteral struct {
	Pairs []KeyValue
}

type KeyValue struct {
	Key   Node
	Value Node
}

type AsyncStatement struct {
	Body []Node
}

type AwaitExpression struct {
	Expression Node
}

type EmitStatement struct {
	EventName string
	Data      Node
}

type ListenStatement struct {
	EventName string
	Handler   Node
}

type ConfigBlock struct {
	Name  string
	Value Node
}

// Parser structure
type Parser struct {
	tokens []token.Token
	pos    int
	errors []ParseError
}

type ParseError struct {
	Line    int
	Column  int
	Message string
	Type    string // "syntax" or "semantic"
}

type ParseResult struct {
	Program  *Program
	Errors   []ParseError
	Warnings []string
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		errors: []ParseError{},
	}
}

// Implement node() interface for all node types
func (n *Program) node()         {}
func (n *Statement) node()       {}
func (n *Expression) node()      {}
func (n *Declaration) node()     {}
func (n *BinaryOp) node()        {}
func (n *UnaryOp) node()         {}
func (n *Literal) node()         {}
func (n *Identifier) node()      {}
func (n *IfStatement) node()     {}
func (n *WhileStatement) node()  {}
func (n *ForStatement) node()    {}
func (n *SwitchStatement) node() {}
func (n *ReturnStatement) node() {}
func (n *CallExpression) node()  {}
func (n *ArrayLiteral) node()    {}
func (n *MapLiteral) node()      {}
func (n *AsyncStatement) node()  {}
func (n *AwaitExpression) node() {}
func (n *EmitStatement) node()   {}
func (n *ListenStatement) node() {}
func (n *ConfigBlock) node()     {}

// String methods for better debugging
func (n *Program) String() string         { return "Program" }
func (n *Statement) String() string       { return fmt.Sprintf("Statement(%s)", n.Type) }
func (n *Expression) String() string      { return fmt.Sprintf("Expression(%s)", n.Type) }
func (n *Declaration) String() string     { return fmt.Sprintf("Declaration(%s: %s)", n.DeclType, n.Name) }
func (n *BinaryOp) String() string        { return fmt.Sprintf("BinaryOp(%s)", n.Operator) }
func (n *UnaryOp) String() string         { return fmt.Sprintf("UnaryOp(%s)", n.Operator) }
func (n *Literal) String() string         { return fmt.Sprintf("Literal(%s: %s)", n.Type, n.Value) }
func (n *Identifier) String() string      { return fmt.Sprintf("Identifier(%s)", n.Name) }
func (n *IfStatement) String() string     { return "IfStatement" }
func (n *WhileStatement) String() string  { return "WhileStatement" }
func (n *ForStatement) String() string    { return "ForStatement" }
func (n *SwitchStatement) String() string { return "SwitchStatement" }
func (n *ReturnStatement) String() string { return "ReturnStatement" }
func (n *CallExpression) String() string  { return "CallExpression" }
func (n *ArrayLiteral) String() string    { return "ArrayLiteral" }
func (n *MapLiteral) String() string      { return "MapLiteral" }
func (n *AsyncStatement) String() string  { return "AsyncStatement" }
func (n *AwaitExpression) String() string { return "AwaitExpression" }
func (n *EmitStatement) String() string   { return fmt.Sprintf("EmitStatement(%s)", n.EventName) }
func (n *ListenStatement) String() string { return fmt.Sprintf("ListenStatement(%s)", n.EventName) }
func (n *ConfigBlock) String() string     { return fmt.Sprintf("ConfigBlock(%s)", n.Name) }

// Helper methods
func (p *Parser) current() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek(offset int) token.Token {
	pos := p.pos + offset
	if pos >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[pos]
}

func (p *Parser) advance() token.Token {
	curr := p.current()
	p.pos++
	return curr
}

func (p *Parser) expect(tokenType token.TokenType) (token.Token, bool) {
	if p.current().Type == tokenType {
		return p.advance(), true
	}
	p.addError(p.current(), fmt.Sprintf("expected %s, got %s", tokenType.String(), p.current().Type.String()))
	return token.Token{}, false
}

// acceptKeywordAsIdent tries to accept a keyword token as an identifier
// Used in contexts where keywords can be used as identifiers (like function arguments)
func (p *Parser) acceptKeywordAsIdent() (token.Token, bool) {
	curr := p.current()
	// Check if current token is a keyword that can be used as identifier
	if isIOKeyword(curr.Type) ||
		isDebugKeyword(curr.Type) ||
		isAIKeyword(curr.Type) ||
		isAgentOperationKeyword(curr.Type) ||
		isConcurrencyKeyword(curr.Type) ||
		isSpecialConstructKeyword(curr.Type) ||
		curr.Type == token.BREAKPOINT ||
		curr.Type == token.CHECKPOINT ||
		curr.Type == token.DEBUG ||
		curr.Type == token.ALLOW ||
		curr.Type == token.PSEUDO ||
		curr.Type == token.TYPE ||
		curr.Type == token.CAST ||
		curr.Type == token.ANY ||
		curr.Type == token.NONE {
		return p.advance(), true
	}
	return token.Token{}, false
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.current().Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) addError(tok token.Token, message string) {
	p.errors = append(p.errors, ParseError{
		Line:    tok.Line,
		Column:  tok.Column,
		Message: message,
		Type:    "syntax",
	})
}

func (p *Parser) skipNewlines() {
	for p.match(token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI) {
		p.advance()
	}
}

// Token classification helpers - these determine context-specific token treatment
// Keyword category mappings for parser
var keywordCategories = map[token.TokenType]string{
	// IO Keywords
	token.READ: "IO", token.WRITE: "IO", token.PRINT: "IO", token.LOG: "IO", token.SAVE: "IO",
	// Debug Keywords
	token.DEBUG: "DEBUG", token.CHECKPOINT: "DEBUG", token.TRACE: "DEBUG", token.ASSERT: "DEBUG",
	token.CONFIGURE: "DEBUG", token.GENERATE_REPORT: "DEBUG",
	// AI Keywords
	token.THINK: "AI", token.ASK: "AI", token.PROMPT: "AI", token.ADAPT: "AI",
	token.CALL_API: "AI", token.TRAIN: "AI", token.EVALUATE: "AI", token.REASON: "AI", token.OBSERVE: "AI",
	// Agent/System Keywords
	token.AGENT: "AGENT_SYSTEM", token.CORE: "AGENT_SYSTEM", token.MODEL: "AGENT_SYSTEM",
	token.TOOLS: "AGENT_SYSTEM", token.ROLE: "AGENT_SYSTEM", token.MODE: "AGENT_SYSTEM",
	token.SYS_PROMPT: "AGENT_SYSTEM", token.MAX_CONCURRENT_REQUESTS: "AGENT_SYSTEM", token.RETRY_POLICY: "AGENT_SYSTEM",
	// Concurrency Keywords
	token.ASYNC: "CONCUR", token.AWAIT: "CONCUR", token.EMIT: "CONCUR", token.LISTEN: "CONCUR",
	token.DISPATCH: "CONCUR", token.MERGE: "CONCUR", token.TASK: "CONCUR", token.CONCURRENT: "CONCUR",
	token.STAGE: "CONCUR", token.GATHER: "CONCUR",
	// Special Construct Keywords
	token.WITH: "SPECIAL", token.THEN: "SPECIAL", token.DEFER: "SPECIAL", token.PIPE: "SPECIAL",
	token.PASS: "SPECIAL", token.ALLOW: "SPECIAL", token.THROUGH: "SPECIAL", token.RANGE: "SPECIAL",
	token.STRATEGY: "SPECIAL", token.TIMEOUT: "SPECIAL", token.WINDOW: "SPECIAL",
	token.ALERT_THRESHOLD: "SPECIAL", token.PSEUDO: "SPECIAL",
	// Agent Operation Keywords
	token.DELEGATE: "AGENT_OP", token.ROUTE: "AGENT_OP", token.COMPOSE: "AGENT_OP", token.INSPECT: "AGENT_OP",
	token.CREATE_POOL: "AGENT_OP", token.MAX_WORKERS: "AGENT_OP", token.SUBMIT: "AGENT_OP",
	token.SUBMIT_DELAYED: "AGENT_OP", token.JOIN: "AGENT_OP", token.NOW: "AGENT_OP",
	token.EXECUTION_TIME: "AGENT_OP", token.REPORT: "AGENT_OP",
}

func isKeywordInCategory(tokenType token.TokenType, categories ...string) bool {
	cat, ok := keywordCategories[tokenType]
	if !ok {
		return false
	}
	for _, c := range categories {
		if cat == c {
			return true
		}
	}
	return false
}

func isIOKeyword(tokenType token.TokenType) bool {
	return isKeywordInCategory(tokenType, "IO")
}

func isDebugKeyword(tokenType token.TokenType) bool {
	return isKeywordInCategory(tokenType, "DEBUG")
}

func isAIKeyword(tokenType token.TokenType) bool {
	return isKeywordInCategory(tokenType, "AI")
}

func isAgentSystemKeyword(tokenType token.TokenType) bool {
	return isKeywordInCategory(tokenType, "AGENT_SYSTEM")
}

func isConcurrencyKeyword(tokenType token.TokenType) bool {
	return isKeywordInCategory(tokenType, "CONCUR")
}

func isSpecialConstructKeyword(tokenType token.TokenType) bool {
	return isKeywordInCategory(tokenType, "SPECIAL")
}

func isAgentOperationKeyword(tokenType token.TokenType) bool {
	return isKeywordInCategory(tokenType, "AGENT_OP")
}

// isValidFieldName checks if a token type can be used as a field name in objects/agents
// These are keywords that can appear as field names in specific contexts (like @agent blocks)
// parser.go - FIX isValidFieldName

func isValidFieldName(tokenType token.TokenType) bool {
	return isAgentSystemKeyword(tokenType) ||
		tokenType == token.TYPE ||
		tokenType == token.CONTEXT ||
		tokenType == token.MEMORY ||
		tokenType == token.FLOW ||
		tokenType == token.TIMEOUT ||
		tokenType == token.TASK ||
		tokenType == token.EMIT ||
		tokenType == token.INPUT ||
		tokenType == token.ACTION ||
		tokenType == token.EXECUTION ||
		tokenType == token.RETRY ||
		tokenType == token.ENABLED ||
		tokenType == token.MAX ||
		tokenType == token.DEPENDS_ON ||
		tokenType == token.CONFIG ||
		tokenType == token.OUTPUTS ||
		tokenType == token.BREAKPOINTS ||
		tokenType == token.ON_CONCUR_DEADLOCK ||
		tokenType == token.ON_LOOP ||
		tokenType == token.ON_TIMEOUT ||
		tokenType == token.GLOBAL ||
		// ADD MORE KEYWORDS THAT CAN BE FIELD NAMES:
		tokenType == token.MAX_CONCURRENT_REQUESTS ||
		tokenType == token.RETRY_POLICY ||
		// Allow identifiers that look like "max_req", "prio" etc - these are IDENTIFIER tokens
		tokenType == token.IDENTIFIER
}

// Main parsing method
func (p *Parser) Parse() *ParseResult {
	program := &Program{Statements: []Node{}}

	for p.current().Type != token.EOF {
		// Skip newlines and comments
		for p.match(token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI) {
			p.advance()
		}

		if p.current().Type == token.EOF {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		// Skip statement terminators
		for p.match(token.STATEMENT_END, token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI) {
			p.advance()
		}
	}

	return &ParseResult{
		Program:  program,
		Errors:   p.errors,
		Warnings: []string{}, // TODO: implement warnings
	}
}

func (p *Parser) parseStatement() Node {
	// Skip comments and newlines
	for p.match(token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI) {
		p.advance()
	}

	if p.current().Type == token.EOF {
		return nil
	}

	switch p.current().Type {
	case token.DECORATOR:
		return p.parseDecoratorStatement()
	case token.ALLOW:
		return p.parseAllowStatement()
	case token.AT_AGENT:
		return p.parseAgentDeclaration()
	case token.TASK:
		return p.parseTaskDeclaration()
	case token.LOOP:
		return p.parseLoopStatement()
	case token.GUARD:
		return p.parseGuardStatement()
	case token.MATCH:
		return p.parseMatchStatement()
	case token.WATCH:
		return p.parseWatchStatement()
	case token.BIND, token.CONST, token.CRAFT:
		return p.parseDeclaration()
	case token.INT_TYPE, token.FLOAT_TYPE, token.STR_TYPE, token.BOOL_TYPE:
		return p.parseTypedDeclaration()
	case token.FN:
		return p.parseFunctionDeclaration()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.SWITCH:
		return p.parseSwitchStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.STRUCT:
		return p.parseStructDeclaration()
	case token.ASYNC:
		return p.parseAsyncStatement()
	case token.AWAIT:
		return p.parseAwaitExpression()
	case token.EMIT:
		return p.parseEmitStatement()
	case token.LISTEN:
		return p.parseListenStatement()
	case token.TRY:
		return p.parseTryStatement()
	case token.ON:
		return p.parseOnStatement()
	case token.WITH:
		return p.parseWithStatement()
	case token.SNAPSHOT:
		return p.parseSnapshotStatement()
	case token.RESTORE:
		return p.parseRestoreStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDecoratorStatement() Node {
	p.advance() // consume @decorator token

	// Read the decorator name and optional arguments
	decorator := &Statement{
		Type:  "decorator",
		Value: map[string]interface{}{},
	}

	// If followed by an identifier (like 'pseudo'), consume it
	if p.match(token.IDENTIFIER) {
		decorator.Value.(map[string]interface{})["name"] = p.current().Lexeme
		p.advance()
	}

	// If followed by LPAREN, parse decorator arguments
	if p.match(token.LPAREN) {
		p.advance()
		args := []Node{}
		for !p.match(token.RPAREN) && p.current().Type != token.EOF {
			args = append(args, p.parseExpression())
			if !p.match(token.RPAREN) {
				if !p.match(token.COMMA) {
					break
				}
				p.advance()
			}
		}
		p.expect(token.RPAREN)
		decorator.Value.(map[string]interface{})["args"] = args
	}

	return decorator
}

func (p *Parser) parseAllowStatement() Node {
	p.advance() // consume 'allow'

	// allow pseudo(anno, trace, breakpoint, 15);
	// Here 'pseudo' is a keyword, and 'trace', 'breakpoint' are also keywords in this context
	var functionName string

	// Expect specific keywords after 'allow'
	if p.match(token.PSEUDO) {
		functionName = p.current().Lexeme
		p.advance()
	} else if p.match(token.IDENTIFIER) {
		functionName = p.current().Lexeme
		p.advance()
	} else {
		p.addError(p.current(), "expected function name after 'allow'")
		return nil
	}

	var args []Node
	if p.match(token.LPAREN) {
		p.advance()
		for !p.match(token.RPAREN) && p.current().Type != token.EOF {
			// In allow pseudo() context, identifiers like 'anno', 'trace', 'breakpoint' are arguments
			// These should be parsed as identifiers or keywords depending on context
			args = append(args, p.parseExpression())
			if p.match(token.COMMA) {
				p.advance()
				if p.match(token.RPAREN) {
					break
				}
			} else if !p.match(token.RPAREN) {
				break
			}
		}
		p.expect(token.RPAREN)
	}

	return &CallExpression{
		Function:  &Identifier{Name: "allow_" + functionName},
		Arguments: args,
	}
}

func (p *Parser) parseAgentDeclaration() Node {
	p.advance() // consume '@agent'

	nameToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' after agent name")
		return nil
	}
	p.advance()

	fields := []Field{}
	for !p.match(token.RBRACE) && p.current().Type != token.EOF {
		p.skipNewlines()

		// Skip comments inside agent body
		for p.match(token.COMMENT_LINE, token.COMMENT_MULTI) {
			p.advance()
			p.skipNewlines()
		}

		if p.match(token.RBRACE) {
			break
		}

		// In agent context, allow keywords as field names (role, model, mode, etc.)
		// Also allow context-sensitive Agent/System keywords (model, mode, role, etc.)
		fieldToken := p.current()
		var fieldName string

		if fieldToken.Type == token.IDENTIFIER {
			fieldName = fieldToken.Lexeme
			p.advance()
		} else if isValidFieldName(fieldToken.Type) {
			fieldName = fieldToken.Lexeme
			p.advance()
		} else {
			p.addError(fieldToken, fmt.Sprintf("expected field name, got %s", fieldToken.Type.String()))
			p.advance()
			continue
		}

		if !p.match(token.COLON) {
			p.addError(p.current(), "expected ':' after field name")
			p.advance()
			continue
		}
		p.advance()

		value := p.parseExpression()
		if value == nil {
			continue
		}

		fields = append(fields, Field{
			Name:  fieldName,
			Value: value,
		})

		// Skip separators and comments
		for p.match(token.COMMA, token.STATEMENT_END, token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI) {
			p.advance()
		}
		p.skipNewlines()
	}

	if !p.match(token.RBRACE) {
		p.addError(p.current(), "expected '}' after agent body")
		return nil
	}
	p.advance()

	return &Declaration{
		DeclType: "agent",
		Name:     nameToken.Lexeme,
		Fields:   fields,
	}
}

func (p *Parser) parseTaskDeclaration() Node {
	p.advance() // consume 'task'

	nameToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' after task name")
		return nil
	}
	p.advance()

	fields := []Field{}
	for !p.match(token.RBRACE) && p.current().Type != token.EOF {
		p.skipNewlines()

		for p.match(token.COMMENT_LINE, token.COMMENT_MULTI) {
			p.advance()
			p.skipNewlines()
		}

		if p.match(token.RBRACE) {
			break
		}

		// In task context, allow both regular identifiers and context-sensitive keywords as field names
		fieldToken := p.current()
		var fieldName string

		if fieldToken.Type == token.IDENTIFIER {
			fieldName = fieldToken.Lexeme
			p.advance()
		} else if isValidFieldName(fieldToken.Type) {
			fieldName = fieldToken.Lexeme
			p.advance()
		} else {
			p.addError(fieldToken, fmt.Sprintf("expected field name, got %s", fieldToken.Type.String()))
			p.advance()
			continue
		}

		if !p.match(token.COLON) {
			p.addError(p.current(), "expected ':' after field name")
			p.advance()
			continue
		}
		p.advance()

		value := p.parseExpression()
		if value == nil {
			continue
		}

		fields = append(fields, Field{
			Name:  fieldName,
			Value: value,
		})

		for p.match(token.COMMA, token.STATEMENT_END, token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI) {
			p.advance()
		}
		p.skipNewlines()
	}

	if !p.match(token.RBRACE) {
		p.addError(p.current(), "expected '}' after task body")
		return nil
	}
	p.advance()

	return &Declaration{
		DeclType: "task",
		Name:     nameToken.Lexeme,
		Fields:   fields,
	}
}

func (p *Parser) parseDeclaration() Node {
	declToken := p.advance()
	declType := declToken.Lexeme

	// Check for new syntax: bind: { ... } or const: { ... }
	if p.match(token.COLON) && p.peek(1).Type == token.LBRACE {
		p.advance() // consume :
		mapNode := p.parseMapLiteral()

		return &ConfigBlock{
			Name:  declType,
			Value: mapNode,
		}
	}

	// Old syntax: bind name : type = value
	nameToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	varType := ""
	var value Node

	// Optional type annotation
	if p.match(token.COLON) {
		p.advance()
		typeToken := p.current()
		if p.match(token.INT_TYPE, token.FLOAT_TYPE, token.STR_TYPE, token.BOOL_TYPE, token.IDENTIFIER) {
			p.advance()
			varType = typeToken.Lexeme
		}
	}

	// Optional assignment
	if p.match(token.ASSIGN, token.BIND_ASSIGN) {
		p.advance()
		value = p.parseExpression()
	}

	return &Declaration{
		DeclType: declType,
		Name:     nameToken.Lexeme,
		Type:     varType,
		Value:    value,
	}
}

func (p *Parser) parseTypedDeclaration() Node {
	typeToken := p.advance()
	varType := typeToken.Lexeme

	declType := "bind" // default

	nameToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	var value Node

	if p.match(token.ASSIGN, token.BIND_ASSIGN) {
		p.advance()
		value = p.parseExpression()
	}

	return &Declaration{
		DeclType: declType,
		Name:     nameToken.Lexeme,
		Type:     varType,
		Value:    value,
	}
}

func (p *Parser) parseSwitchStatement() Node {
	p.advance() // consume 'switch'

	expr := p.parseExpression()
	if expr == nil {
		p.addError(p.current(), "expected expression after 'switch'")
		return nil
	}

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' after switch expression")
		return nil
	}
	p.advance()

	cases := []CaseClause{}
	var defaultBody []Node

	for !p.match(token.RBRACE) && p.current().Type != token.EOF {
		p.skipNewlines()

		if p.match(token.RBRACE) {
			break
		}

		if p.match(token.CASE) {
			p.advance()

			caseValue := p.parseExpression()
			if caseValue == nil {
				p.addError(p.current(), "expected expression after 'case'")
				continue
			}

			if !p.match(token.LBRACE) {
				p.addError(p.current(), "expected '{' after case value")
				continue
			}
			p.advance()

			caseBody := p.parseBlock()

			cases = append(cases, CaseClause{
				Value: caseValue,
				Body:  caseBody,
			})

		} else if p.match(token.DEFAULT) {
			p.advance()

			if !p.match(token.LBRACE) {
				p.addError(p.current(), "expected '{' after 'default'")
				continue
			}
			p.advance()

			defaultBody = p.parseBlock()

		} else {
			p.addError(p.current(), "expected 'case' or 'default' in switch statement")
			p.advance()
		}
	}

	if !p.match(token.RBRACE) {
		p.addError(p.current(), "expected '}' after switch body")
		return nil
	}
	p.advance()

	return &SwitchStatement{
		Expression: expr,
		Cases:      cases,
		Default:    defaultBody,
	}
}

func (p *Parser) parseFunctionDeclaration() Node {
	p.advance() // consume 'fn'

	// Check for function decorator: fn@agent name(...)
	var decorator string
	if p.match(token.AT_AGENT, token.AT_TASK, token.AT_STEP, token.AT_INTENT, token.AT_EXPLAIN, token.DECORATOR) {
		decorator = p.current().Lexeme
		p.advance()
	}

	nameToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	if !p.match(token.LPAREN) {
		p.addError(p.current(), "expected '(' after function name")
		return nil
	}
	p.advance()

	params := []Parameter{}
	for !p.match(token.RPAREN) && p.current().Type != token.EOF {
		// Allow both identifiers and keywords as parameter names
		var paramName string
		if p.match(token.IDENTIFIER) {
			paramName = p.current().Lexeme
			p.advance()
		} else if isValidFieldName(p.current().Type) || p.match(token.PROMPT, token.THINK, token.ASK) {
			paramName = p.current().Lexeme
			p.advance()
		} else {
			p.addError(p.current(), "expected parameter name")
			break
		}

		paramType := ""
		if p.match(token.COLON) {
			p.advance()
			if p.match(token.INT_TYPE, token.FLOAT_TYPE, token.STR_TYPE, token.BOOL_TYPE, token.IDENTIFIER) {
				typeToken := p.current()
				p.advance()
				paramType = typeToken.Lexeme
			}
		}

		params = append(params, Parameter{Name: paramName, Type: paramType})

		if !p.match(token.RPAREN) {
			p.expect(token.COMMA)
		}
	}

	if !p.match(token.RPAREN) {
		p.addError(p.current(), "expected ')'")
		return nil
	}
	p.advance()

	returnType := ""
	if p.match(token.ARROW, token.FAT_ARROW) {
		p.advance()
		if p.match(token.INT_TYPE, token.FLOAT_TYPE, token.STR_TYPE, token.BOOL_TYPE, token.IDENTIFIER) {
			typeToken := p.current()
			p.advance()
			returnType = typeToken.Lexeme
		}
	}

	// Handle :: block delimiter
	if p.match(token.DOUBLE_COLON) {
		p.advance()
	}

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' for function body")
		return nil
	}
	p.advance()

	body := p.parseBlock()

	return &Declaration{
		DeclType:  "fn",
		Name:      nameToken.Lexeme,
		Type:      returnType,
		Params:    params,
		Body:      body,
		Decorator: decorator,
	}
}

func (p *Parser) parseAsyncStatement() Node {
	p.advance() // consume 'async'

	if p.match(token.FN) {
		return p.parseFunctionDeclaration()
	}

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' or 'fn' after 'async'")
		return nil
	}
	p.advance()

	body := p.parseBlock()

	return &AsyncStatement{Body: body}
}

func (p *Parser) parseAwaitExpression() Node {
	p.advance() // consume 'await'
	// Special: allow await all { ... } and await race { ... }
	if p.match(token.IDENTIFIER) && (p.current().Lexeme == "all" || p.current().Lexeme == "race") {
		p.advance()
		if p.match(token.LBRACE) {
			p.advance()
			exprs := []Node{}
			for !p.match(token.RBRACE) && p.current().Type != token.EOF {
				p.skipNewlines()
				expr := p.parseExpression()
				if expr != nil {
					exprs = append(exprs, expr)
				}
				for p.match(token.COMMA, token.STATEMENT_END, token.NEWLINE) {
					p.advance()
				}
			}
			p.expect(token.RBRACE)
			return &AwaitExpression{Expression: &ArrayLiteral{Elements: exprs}}
		}
	}
	expr := p.parseExpression()
	return &AwaitExpression{Expression: expr}
}

func (p *Parser) parseEmitStatement() Node {
	p.advance() // consume 'emit'

	// Event name can be an identifier or a keyword used as identifier
	var eventName string
	if p.match(token.IDENTIFIER) {
		eventName = p.current().Lexeme
		p.advance()
	} else if tok, ok := p.acceptKeywordAsIdent(); ok {
		eventName = tok.Lexeme
	} else {
		p.addError(p.current(), "expected event name after 'emit'")
		return nil
	}

	var data Node
	if p.match(token.LBRACE) {
		data = p.parseMapLiteral()
	}

	return &EmitStatement{
		EventName: eventName,
		Data:      data,
	}
}

func (p *Parser) parseListenStatement() Node {
	p.advance() // consume 'listen'

	eventToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' after event name")
		return nil
	}
	p.advance()

	handler := p.parseExpression()

	if !p.match(token.RBRACE) {
		p.addError(p.current(), "expected '}'")
		return nil
	}
	p.advance()

	return &ListenStatement{
		EventName: eventToken.Lexeme,
		Handler:   handler,
	}
}

func (p *Parser) parseTryStatement() Node {
	p.advance() // consume 'try'

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' after 'try'")
		return nil
	}
	p.advance()

	tryBody := p.parseBlock()

	var catchBody []Node
	var catchParam string

	if p.match(token.CATCH) {
		p.advance()

		if p.match(token.IDENTIFIER) {
			catchParam = p.advance().Lexeme
		}

		if !p.match(token.LBRACE) {
			p.addError(p.current(), "expected '{' after 'catch'")
			return nil
		}
		p.advance()

		catchBody = p.parseBlock()
	}

	return &Statement{
		Type: "try-catch",
		Value: map[string]interface{}{
			"try":        tryBody,
			"catch":      catchBody,
			"catchParam": catchParam,
		},
	}
}

func (p *Parser) parseIfStatement() Node {
	p.advance() // consume 'if'

	condition := p.parseExpression()
	if condition == nil {
		p.addError(p.current(), "expected condition after 'if'")
		return nil
	}

	// Handle :: inline syntax
	if p.match(token.DOUBLE_COLON) {
		p.advance()
		stmt := p.parseStatement()
		thenBody := []Node{}
		if stmt != nil {
			thenBody = append(thenBody, stmt)
		}

		// Skip newlines and statement terminators before checking for ELSE
		for p.match(token.NEWLINE, token.STATEMENT_END, token.SEMICOLON) {
			p.advance()
		}

		elseBody := []Node{}
		if p.match(token.ELSE) {
			p.advance()
			if p.match(token.DOUBLE_COLON) {
				p.advance()
			}
			stmt := p.parseStatement()
			if stmt != nil {
				elseBody = append(elseBody, stmt)
			}
		}

		return &IfStatement{
			Condition: condition,
			ThenBody:  thenBody,
			ElseBody:  elseBody,
		}
	}

	// Traditional block syntax
	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' or '::' after if condition")
		return nil
	}
	p.advance()

	thenBody := p.parseBlock()

	elseBody := []Node{}
	if p.match(token.ELSE) {
		p.advance()
		if p.match(token.IF) {
			elseBody = []Node{p.parseIfStatement()}
		} else if p.match(token.LBRACE) {
			p.advance()
			elseBody = p.parseBlock()
		}
	}

	return &IfStatement{
		Condition: condition,
		ThenBody:  thenBody,
		ElseBody:  elseBody,
	}
}

func (p *Parser) parseWhileStatement() Node {
	p.advance() // consume 'while'

	condition := p.parseExpression()
	if condition == nil {
		p.addError(p.current(), "expected condition after 'while'")
		return nil
	}

	// Handle :: inline syntax
	if p.match(token.DOUBLE_COLON) {
		p.advance()
		stmt := p.parseStatement()
		body := []Node{}
		if stmt != nil {
			body = append(body, stmt)
		}
		return &WhileStatement{
			Condition: condition,
			Body:      body,
		}
	}

	// Traditional block syntax
	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' or '::' after while condition")
		return nil
	}
	p.advance()

	body := p.parseBlock()

	return &WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseForStatement() Node {
	p.advance() // consume 'for'

	var init, condition, update Node
	concurrent := false

	// Check for 'in' style loop (for var in collection)
	// The loop variable can be an identifier or a keyword used as identifier
	var loopVar string
	var isInLoop bool

	if p.match(token.IDENTIFIER) {
		loopVar = p.current().Lexeme
		p.advance()
		if p.match(token.IDENTIFIER) && p.current().Lexeme == "in" {
			isInLoop = true
			p.advance() // consume 'in'
		}
	} else if isValidFieldName(p.current().Type) || isIOKeyword(p.current().Type) || isDebugKeyword(p.current().Type) || isAIKeyword(p.current().Type) {
		// Allow keywords as loop variables
		loopVar = p.current().Lexeme
		p.advance()
		if p.match(token.IDENTIFIER) && p.current().Lexeme == "in" {
			isInLoop = true
			p.advance() // consume 'in'
		}
	}

	if isInLoop {
		// for loopVar in collection { ... }
		condition = &BinaryOp{
			Left:     &Identifier{Name: loopVar},
			Operator: "in",
			Right:    p.parseExpression(),
		}
	} else {
		// Traditional for loop: for init; condition; update { ... }
		// If we already parsed a variable, use it as init
		if loopVar != "" {
			init = &Identifier{Name: loopVar}
		} else if !p.match(token.STATEMENT_END) {
			init = p.parseExpression()
		}
		p.expect(token.STATEMENT_END)

		if !p.match(token.STATEMENT_END) {
			condition = p.parseExpression()
		}
		p.expect(token.STATEMENT_END)

		if !p.match(token.LBRACE, token.CONCURRENT) {
			update = p.parseExpression()
		}
	}

	if p.match(token.CONCURRENT) {
		p.advance()
		concurrent = true
	}

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' for for loop body")
		return nil
	}
	p.advance()

	body := p.parseBlock()

	return &ForStatement{
		Init:       init,
		Condition:  condition,
		Update:     update,
		Body:       body,
		Concurrent: concurrent,
	}
}

func (p *Parser) parseReturnStatement() Node {
	p.advance() // consume 'return'

	var value Node
	if !p.match(token.STATEMENT_END, token.NEWLINE, token.RBRACE, token.EOF) {
		value = p.parseExpression()
	}

	return &ReturnStatement{Value: value}
}

func (p *Parser) parseStructDeclaration() Node {
	p.advance() // consume 'struct'

	nameToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' after struct name")
		return nil
	}
	p.advance()

	fields := []Field{}
	for !p.match(token.RBRACE) && p.current().Type != token.EOF {
		// Skip comments and newlines before each field
		for p.match(token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI) {
			p.advance()
		}

		if p.match(token.RBRACE) {
			break
		}

		fieldToken, ok := p.expect(token.IDENTIFIER)
		if !ok {
			break
		}

		fieldType := ""
		if p.match(token.COLON) {
			p.advance()
			if p.match(token.INT_TYPE, token.FLOAT_TYPE, token.STR_TYPE, token.BOOL_TYPE, token.IDENTIFIER) {
				typeToken := p.current()
				p.advance()
				fieldType = typeToken.Lexeme
			}
		}

		fields = append(fields, Field{
			Name: fieldToken.Lexeme,
			Type: fieldType,
		})

		for p.match(token.STATEMENT_END, token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI, token.COMMA) {
			p.advance()
		}
	}

	if !p.match(token.RBRACE) {
		p.addError(p.current(), "expected '}' after struct body")
		return nil
	}
	p.advance()

	return &Declaration{
		DeclType: "struct",
		Name:     nameToken.Lexeme,
		Fields:   fields,
	}
}

func (p *Parser) parseBlock() []Node {
	statements := []Node{}

	for !p.match(token.RBRACE) && p.current().Type != token.EOF {
		p.skipNewlines()

		if p.match(token.RBRACE) {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}

		if p.match(token.STATEMENT_END, token.NEWLINE) {
			p.advance()
		}
	}

	if p.match(token.RBRACE) {
		p.advance()
	}

	return statements
}

func (p *Parser) parseExpressionStatement() Node {
	// Check for configuration block syntax: identifier: { ... } or keyword: { ... }
	// This handles both "config_name: { ... }" and "bind: { ... }", "const: { ... }", etc.
	if p.match(token.IDENTIFIER, token.BIND, token.CONST, token.CRAFT) {
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == token.COLON {
			if p.pos+2 < len(p.tokens) && p.tokens[p.pos+2].Type == token.LBRACE {
				nameToken := p.current()
				p.advance() // consume identifier/keyword
				p.advance() // consume colon

				mapNode := p.parseMapLiteral()

				return &ConfigBlock{
					Name:  nameToken.Lexeme,
					Value: mapNode,
				}
			}
		}
	}

	expr := p.parseExpression()
	return expr
}

func (p *Parser) parseExpression() Node {
	return p.parseAssignment()
}

func (p *Parser) parseTupleLiteral() Node {
	p.advance() // consume '('

	elements := []Node{}
	for !p.match(token.RPAREN) && p.current().Type != token.EOF {
		elements = append(elements, p.parseExpression())
		if !p.match(token.RPAREN) {
			if p.match(token.COMMA) {
				p.advance()
			} else {
				break
			}
		}
	}

	p.expect(token.RPAREN)

	// If single element, it's just a grouped expression
	if len(elements) == 1 {
		return elements[0]
	}

	// Multiple elements = tuple
	return &ArrayLiteral{Elements: elements} // Reuse ArrayLiteral for now
}

func (p *Parser) parseAssignment() Node {
	left := p.parsePipeline()

	if p.match(token.ASSIGN, token.BIND_ASSIGN, token.PLUS_ASSIGN, token.MINUS_ASSIGN,
		token.MULT_ASSIGN, token.DIV_ASSIGN, token.MOD_ASSIGN) {
		opToken := p.advance()
		right := p.parseExpression()
		return &BinaryOp{Left: left, Operator: opToken.Lexeme, Right: right}
	}

	return left
}

func (p *Parser) parsePipeline() Node {
	left := p.parseLogicalOr()

	// Handle pipeline operators: |> (sequential) and |>> (parallel)
	for p.match(token.PIPE_RIGHT, token.PIPE_PARALLEL) {
		opToken := p.advance()
		right := p.parseLogicalOr()
		left = &BinaryOp{Left: left, Operator: opToken.Lexeme, Right: right}
	}

	return left
}

func (p *Parser) parseLogicalOr() Node {
	left := p.parseLogicalAnd()

	for p.match(token.OR) {
		opToken := p.advance()
		right := p.parseLogicalAnd()
		left = &BinaryOp{Left: left, Operator: opToken.Lexeme, Right: right}
	}

	return left
}

func (p *Parser) parseLogicalAnd() Node {
	left := p.parseEquality()

	for p.match(token.AND) {
		opToken := p.advance()
		right := p.parseEquality()
		left = &BinaryOp{Left: left, Operator: opToken.Lexeme, Right: right}
	}

	return left
}

func (p *Parser) parseEquality() Node {
	left := p.parseComparison()

	for p.match(token.EQ, token.NEQ) {
		opToken := p.advance()
		right := p.parseComparison()
		left = &BinaryOp{Left: left, Operator: opToken.Lexeme, Right: right}
	}

	return left
}

func (p *Parser) parseComparison() Node {
	left := p.parseAdditive()

	for p.match(token.LT, token.LTE, token.GT, token.GTE) {
		opToken := p.advance()
		right := p.parseAdditive()
		left = &BinaryOp{Left: left, Operator: opToken.Lexeme, Right: right}
	}

	return left
}

func (p *Parser) parseAdditive() Node {
	left := p.parseMultiplicative()

	for p.match(token.PLUS, token.MINUS) {
		opToken := p.advance()
		right := p.parseMultiplicative()
		left = &BinaryOp{Left: left, Operator: opToken.Lexeme, Right: right}
	}

	return left
}

func (p *Parser) parseMultiplicative() Node {
	left := p.parseUnary()

	for p.match(token.MULTIPLY, token.DIVIDE, token.MODULO) {
		opToken := p.advance()
		right := p.parseUnary()
		left = &BinaryOp{Left: left, Operator: opToken.Lexeme, Right: right}
	}

	return left
}

func (p *Parser) parseUnary() Node {
	if p.match(token.NOT, token.MINUS, token.PLUS) {
		opToken := p.advance()
		operand := p.parseUnary()
		return &UnaryOp{Operator: opToken.Lexeme, Operand: operand}
	}

	if p.match(token.AWAIT) {
		return p.parseAwaitExpression()
	}

	return p.parsePostfix()
}

func (p *Parser) parsePostfix() Node {
	left := p.parsePrimary()

	for {
		if p.match(token.PIPE_RIGHT, token.PIPE_PARALLEL) {
			op := p.advance()
			right := p.parsePrimary()
			left = &BinaryOp{Left: left, Operator: op.Lexeme, Right: right}
		} else if p.match(token.LPAREN) {
			p.advance()
			args := []Node{}
			for !p.match(token.RPAREN) && p.current().Type != token.EOF {
				args = append(args, p.parseExpression())
				if !p.match(token.RPAREN) {
					if !p.match(token.COMMA) {
						break
					}
					p.advance()
				}
			}
			p.expect(token.RPAREN)
			left = &CallExpression{Function: left, Arguments: args}
		} else if p.match(token.LBRACKET) {
			p.advance()
			index := p.parseExpression()
			p.expect(token.RBRACKET)
			left = &BinaryOp{Left: left, Operator: "[]", Right: index}
		} else if p.match(token.DOT) {
			p.advance()
			fieldToken, ok := p.expect(token.IDENTIFIER)
			if !ok {
				break
			}
			left = &BinaryOp{Left: left, Operator: ".", Right: &Identifier{Name: fieldToken.Lexeme}}
		} else if p.match(token.INCREMENT, token.DECREMENT) {
			opToken := p.advance()
			left = &UnaryOp{Operator: opToken.Lexeme + "_post", Operand: left}
		} else if p.match(token.ARROW) {
			p.advance()
			command := p.parseExpression()
			left = &BinaryOp{Left: left, Operator: "->", Right: command}
		} else {
			break
		}
	}

	return left
}

func (p *Parser) parsePrimary() Node {
	switch p.current().Type {
	case token.IDENTIFIER:
		tok := p.advance()
		return &Identifier{Name: tok.Lexeme}

	case token.INTEGER:
		tok := p.advance()
		return &Literal{Type: "int", Value: tok.Lexeme}

	case token.FLOAT:
		tok := p.advance()
		return &Literal{Type: "float", Value: tok.Lexeme}

	case token.STRING:
		tok := p.advance()
		return &Literal{Type: "string", Value: tok.Lexeme}

	case token.TRUE, token.FALSE:
		tok := p.advance()
		return &Literal{Type: "bool", Value: tok.Lexeme}

	case token.NULL:
		tok := p.advance()
		return &Literal{Type: "null", Value: tok.Lexeme}

	case token.LPAREN:
		return p.parseTupleLiteral()

	case token.LBRACKET:
		return p.parseArrayLiteral()

	case token.LBRACE:
		return p.parseMapLiteral()

	case token.ASYNC:
		return p.parseAsyncStatement()

	// Keywords that can be used as identifiers in expression context
	// (e.g., function calls like print(), trace(), now(), breakpoint, etc.)
	default:
		if isIOKeyword(p.current().Type) ||
			isDebugKeyword(p.current().Type) ||
			isAIKeyword(p.current().Type) ||
			isAgentOperationKeyword(p.current().Type) ||
			isConcurrencyKeyword(p.current().Type) ||
			isSpecialConstructKeyword(p.current().Type) ||
			isValidFieldName(p.current().Type) ||
			p.match(token.TRACE, token.NOW, token.GATHER, token.CREATE_POOL,
				token.BREAKPOINT, token.CHECKPOINT, token.DEBUG, token.ALLOW, token.PSEUDO,
				token.TYPE, token.CAST, token.ANY, token.NONE, token.TRAIT) {
			tok := p.advance()
			return &Identifier{Name: tok.Lexeme}
		}

		p.addError(p.current(), fmt.Sprintf("unexpected token: %s", p.current().Type.String()))
		p.advance()
		return nil
	}
}

func (p *Parser) parseArrayLiteral() Node {
	p.advance() // consume '['

	elements := []Node{}
	for !p.match(token.RBRACKET) && p.current().Type != token.EOF {
		elements = append(elements, p.parseExpression())
		if !p.match(token.RBRACKET) {
			if p.match(token.COMMA) {
				p.advance()
			} else {
				break
			}
		}
	}

	p.expect(token.RBRACKET)
	return &ArrayLiteral{Elements: elements}
}

func (p *Parser) parseMapLiteral() Node {
	p.advance() // consume '{'

	pairs := []KeyValue{}
	for !p.match(token.RBRACE) && p.current().Type != token.EOF {
		p.skipNewlines() // This now also skips comments

		if p.match(token.RBRACE) {
			break
		}

		// Parse key - allow identifiers and ALL keywords as field names in maps
		var key Node
		keyToken := p.current()

		// Be MORE permissive - allow almost any token as a map key
		if keyToken.Type == token.IDENTIFIER {
			p.advance()
			key = &Identifier{Name: keyToken.Lexeme}
		} else if isValidFieldName(keyToken.Type) {
			p.advance()
			key = &Identifier{Name: keyToken.Lexeme}
		} else if isIOKeyword(keyToken.Type) ||
			isDebugKeyword(keyToken.Type) ||
			isAIKeyword(keyToken.Type) ||
			isConcurrencyKeyword(keyToken.Type) {
			// Allow IO, debug, AI, concurrency keywords as field names too
			p.advance()
			key = &Identifier{Name: keyToken.Lexeme}
		} else {
			p.addError(keyToken, fmt.Sprintf("expected field name, got %s", keyToken.Type.String()))
			p.advance()
			continue
		}

		// Check for shorthand syntax: { x } means { x: x }
		var value Node
		if p.match(token.COMMA, token.RBRACE, token.NEWLINE, token.STATEMENT_END) {
			value = key
		} else {
			if !p.match(token.COLON) {
				p.addError(p.current(), "expected ':' after field name")
				p.advance()
				continue
			}
			p.advance()

			// Handle nested blocks or statements in object literals
			// Example: concurrency: { max_req: 3, timeout: 60s }
			if p.match(token.LBRACE) {
				value = p.parseMapLiteral()
			} else {
				value = p.parseExpression()
				if value == nil {
					continue
				}
			}
		}

		pairs = append(pairs, KeyValue{Key: key, Value: value})

		// Skip separators and comments
		for p.match(token.COMMA, token.STATEMENT_END, token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI) {
			p.advance()
		}
		p.skipNewlines()
	}

	if !p.match(token.RBRACE) {
		p.addError(p.current(), "expected '}' after map literal")
		return nil
	}
	p.advance()

	return &MapLiteral{Pairs: pairs}
}

// FormatGrammarTree outputs a grammar-based parse tree with ASCII formatting
func formatGrammarTree(node Node, isLast bool, prefix string) string {
	if node == nil {
		return ""
	}

	result := ""
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	switch n := node.(type) {
	case *Program:
		result = prefix + connector + "SYNTA_PROGRAM\n"
		if len(n.Statements) > 0 {
			result += prefix + ("├── ") + "STMT_LIST\n"
			for i, stmt := range n.Statements {
				isLastStmt := (i == len(n.Statements)-1)
				result += formatGrammarTree(stmt, isLastStmt, prefix+"│   ")
			}
		}
		return result

	case *Declaration:
		declNode := "DECL_STMT"
		switch n.DeclType {
		case "agent":
			declNode = "AGENT_DECL"
		case "task":
			declNode = "TASK_DECL"
		case "fn":
			declNode = "FN_DECL"
		case "struct":
			declNode = "STRUCT_DECL"
		}
		result = prefix + connector + declNode + "\n"
		result += childPrefix + "├── IDENTIFIER \"" + n.Name + "\"\n"
		if len(n.Fields) > 0 {
			for i, f := range n.Fields {
				isLast := (i == len(n.Fields)-1)
				connStr := "├── "
				if isLast {
					connStr = "└── "
				}
				result += childPrefix + connStr + "PAIR\n"
				nextPrefix := childPrefix
				if isLast {
					nextPrefix += "    "
				} else {
					nextPrefix += "│   "
				}
				result += nextPrefix + "├── IDENTIFIER \"" + f.Name + "\"\n"
				if f.Value != nil {
					result += formatGrammarTree(f.Value, true, nextPrefix)
				}
			}
		}
		if len(n.Body) > 0 {
			result += childPrefix + "└── BLOCK\n"
			for i, stmt := range n.Body {
				isLastStmt := (i == len(n.Body)-1)
				result += formatGrammarTree(stmt, isLastStmt, childPrefix+"    ")
			}
		}
		return result

	case *Literal:
		litType := "LITERAL"
		switch n.Type {
		case "string":
			litType = "STRING_LIT"
		case "int":
			litType = "INT_LIT"
		case "float":
			litType = "FLOAT_LIT"
		case "bool":
			litType = "BOOL_LIT"
		case "null":
			litType = "NULL_LIT"
		}
		return prefix + connector + litType + " \"" + n.Value + "\"\n"

	case *Identifier:
		// Check if the identifier is a keyword and output semantic group if so
		semanticGroup := getSemanticGroupForIdentifier(n.Name)
		return prefix + connector + semanticGroup + " \"" + n.Name + "\"\n"

	case *BinaryOp:
		result = prefix + connector + "BINARY_EXPR\n"
		result += childPrefix + "├── EXPR\n"
		result += formatGrammarTree(n.Left, false, childPrefix+"│   ")
		result += childPrefix + "├── OPERATOR \"" + n.Operator + "\"\n"
		result += childPrefix + "└── EXPR\n"
		result += formatGrammarTree(n.Right, true, childPrefix+"    ")
		return result

	case *CallExpression:
		result = prefix + connector + "CALL_EXPR\n"
		result += childPrefix + "├── MEMBER_EXPR\n"
		result += formatGrammarTree(n.Function, false, childPrefix+"│   ")
		if len(n.Arguments) > 0 {
			result += childPrefix + "└── ARG_LIST\n"
			for i, arg := range n.Arguments {
				isLast := (i == len(n.Arguments)-1)
				result += formatGrammarTree(arg, isLast, childPrefix+"    ")
			}
		}
		return result

	case *ArrayLiteral:
		result = prefix + connector + "ARRAY_LITERAL\n"
		for i, elem := range n.Elements {
			isLast := (i == len(n.Elements)-1)
			result += formatGrammarTree(elem, isLast, childPrefix)
		}
		return result

	case *MapLiteral:
		result = prefix + connector + "MAP_LITERAL\n"
		for i, pair := range n.Pairs {
			isLast := (i == len(n.Pairs)-1)
			connStr := "├── "
			if isLast {
				connStr = "└── "
			}
			result += childPrefix + connStr + "PAIR\n"
			pairPrefix := childPrefix
			if isLast {
				pairPrefix += "    "
			} else {
				pairPrefix += "│   "
			}
			result += formatGrammarTree(pair.Key, false, pairPrefix)
			result += formatGrammarTree(pair.Value, true, pairPrefix)
		}
		return result

	case *IfStatement:
		result = prefix + connector + "IF_STMT\n"
		result += childPrefix + "├── CONDITION\n"
		result += formatGrammarTree(n.Condition, false, childPrefix+"│   ")
		result += childPrefix + "├── BLOCK\n"
		for i, stmt := range n.ThenBody {
			isLast := (i == len(n.ThenBody)-1)
			result += formatGrammarTree(stmt, isLast, childPrefix+"│   ")
		}
		if len(n.ElseBody) > 0 {
			result += childPrefix + "└── ELSE_BLOCK\n"
			for i, stmt := range n.ElseBody {
				isLast := (i == len(n.ElseBody)-1)
				result += formatGrammarTree(stmt, isLast, childPrefix+"    ")
			}
		}
		return result

	case *ReturnStatement:
		result = prefix + connector + "RETURN_STMT\n"
		if n.Value != nil {
			result += formatGrammarTree(n.Value, true, childPrefix)
		}
		return result

	case *AsyncStatement:
		result = prefix + connector + "ASYNC_STMT\n"
		for i, stmt := range n.Body {
			isLast := (i == len(n.Body)-1)
			result += formatGrammarTree(stmt, isLast, childPrefix)
		}
		return result

	case *AwaitExpression:
		result = prefix + connector + "AWAIT_EXPR\n"
		if n.Expression != nil {
			result += formatGrammarTree(n.Expression, true, childPrefix)
		}
		return result

	case *EmitStatement:
		result = prefix + connector + "EMIT_STMT\n"
		result += childPrefix + "├── IDENTIFIER \"" + n.EventName + "\"\n"
		if n.Data != nil {
			result += formatGrammarTree(n.Data, true, childPrefix)
		}
		return result

	case *WhileStatement:
		result = prefix + connector + "WHILE_STMT\n"
		result += childPrefix + "├── CONDITION\n"
		result += formatGrammarTree(n.Condition, false, childPrefix+"│   ")
		result += childPrefix + "└── BLOCK\n"
		for i, stmt := range n.Body {
			isLast := (i == len(n.Body)-1)
			result += formatGrammarTree(stmt, isLast, childPrefix+"    ")
		}
		return result

	case *ForStatement:
		result = prefix + connector + "FOR_STMT\n"
		if n.Init != nil {
			result += childPrefix + "├── INIT\n"
			result += formatGrammarTree(n.Init, false, childPrefix+"│   ")
		}
		if n.Condition != nil {
			result += childPrefix + "├── CONDITION\n"
			result += formatGrammarTree(n.Condition, false, childPrefix+"│   ")
		}
		if n.Update != nil {
			result += childPrefix + "├── UPDATE\n"
			result += formatGrammarTree(n.Update, false, childPrefix+"│   ")
		}
		result += childPrefix + "└── BLOCK\n"
		for i, stmt := range n.Body {
			isLast := (i == len(n.Body)-1)
			result += formatGrammarTree(stmt, isLast, childPrefix+"    ")
		}
		return result

	default:
		return prefix + connector + "STMT\n"
	}
}

// Helper function to get semantic group for an identifier
// If it's a keyword, returns the semantic group; otherwise returns IDENTIFIER
func getSemanticGroupForIdentifier(name string) string {
	// Check if it's a keyword
	if tok, ok := token.Keywords[name]; ok {
		return tok.GetSemanticGroup()
	}
	// Check if it's an agent/system keyword
	if tokType, ok := token.IsAgentSystemKeyword(name); ok {
		return tokType.GetSemanticGroup()
	}
	return "IDENTIFIER"
}

// AST to string representation - Grammar-based
func FormatAST(node Node, indent int) string {
	if prog, ok := node.(*Program); ok {
		return formatGrammarTree(prog, true, "")
	}
	return formatGrammarTree(node, true, "")
}

// Parse from source code directly
func ParseFromSource(source string) *ParseResult {
	lex := lexer.New(source)
	tokens := lex.Tokenize()
	parser := NewParser(tokens)
	return parser.Parse()
}

// parseLoopStatement handles loop syntax: loop while/foreach/parallel
func (p *Parser) parseLoopStatement() Node {
	p.advance() // consume 'loop'

	var loopType string
	var init Node
	var condition Node
	var update Node
	var variable string
	var iterable Node
	var concurrent bool

	if p.match(token.WHILE) {
		loopType = "while"
		p.advance()

		// Check for "i from 0 to 5" syntax
		if p.match(token.IDENTIFIER) && p.peek(1).Type == token.FROM {
			variable = p.current().Lexeme
			p.advance() // consume variable
			p.advance() // consume 'from'

			init = p.parseExpression()

			if p.match(token.IDENTIFIER) && p.current().Lexeme == "to" {
				p.advance()
				condition = p.parseExpression()
			}
		} else {
			condition = p.parseExpression()
		}
	} else if p.match(token.IDENTIFIER) {
		loopType = p.current().Lexeme
		if loopType == "foreach" {
			p.advance()
			// foreach var in collection
			variable = p.current().Lexeme
			p.advance()

			if !p.match(token.IDENTIFIER) || p.current().Lexeme != "in" {
				p.addError(p.current(), "expected 'in' after loop variable")
			} else {
				p.advance()
			}

			iterable = p.parseExpression()
		} else if loopType == "parallel" {
			p.advance()
			concurrent = true
			// parallel var in collection
			variable = p.current().Lexeme
			p.advance()

			if !p.match(token.IDENTIFIER) || p.current().Lexeme != "in" {
				p.addError(p.current(), "expected 'in' after loop variable")
			} else {
				p.advance()
			}

			iterable = p.parseExpression()
		}
	}

	// Handle :: delimiter if present
	if p.match(token.DOUBLE_COLON) {
		p.advance()
	}

	body := []Node{}
	if p.match(token.LBRACE) {
		p.advance()
		body = p.parseBlock()
	} else {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}

	return &ForStatement{
		Variable:   variable,
		Init:       init,
		Condition:  condition,
		Update:     update,
		Iterable:   iterable,
		Body:       body,
		Concurrent: concurrent,
	}
}

// parseGuardStatement handles guard condition :: statement else :: statement
func (p *Parser) parseGuardStatement() Node {
	p.advance() // consume 'guard'

	condition := p.parseExpression()

	// Handle :: delimiter
	if !p.match(token.DOUBLE_COLON) {
		p.addError(p.current(), "expected '::' after guard condition")
	} else {
		p.advance()
	}

	thenBody := []Node{}
	if p.match(token.LBRACE) {
		p.advance()
		thenBody = p.parseBlock()
	} else {
		thenBody = append(thenBody, p.parseStatement())
	}

	elseBody := []Node{}
	if p.match(token.ELSE) {
		p.advance()
		if p.match(token.DOUBLE_COLON) {
			p.advance()
		}
		if p.match(token.LBRACE) {
			p.advance()
			elseBody = p.parseBlock()
		} else {
			elseBody = append(elseBody, p.parseStatement())
		}
	}

	return &IfStatement{
		Condition: condition,
		ThenBody:  thenBody,
		ElseBody:  elseBody,
	}
}

// parseMatchStatement handles match expression :: { case :: action, ... }
func (p *Parser) parseMatchStatement() Node {
	p.advance() // consume 'match'

	expr := p.parseExpression()

	// Handle :: delimiter
	if !p.match(token.DOUBLE_COLON) {
		p.addError(p.current(), "expected '::' after match expression")
	} else {
		if !p.match(token.LBRACE) {
			p.addError(p.current(), "expected '{' after match expression")
			return nil
		}
		p.advance()

		cases := []CaseClause{}
		for !p.match(token.RBRACE) && p.current().Type != token.EOF {
			p.skipNewlines()
			if p.match(token.RBRACE) {
				break
			}
			// Accept string | _ | identifier | expression :: statement
			val := p.parseExpression()
			if !p.match(token.DOUBLE_COLON) {
				p.addError(p.current(), "expected '::' after match case")
				p.advance()
				continue
			}
			p.advance()
			body := []Node{p.parseStatement()}
			cases = append(cases, CaseClause{Value: val, Body: body})
			for p.match(token.COMMA, token.STATEMENT_END, token.NEWLINE) {
				p.advance()
			}
		}
		p.expect(token.RBRACE)
		return &SwitchStatement{
			Expression: expr,
			Cases:      cases,
			Default:    nil,
		}
	}
	return nil
}

// parseWatchStatement handles watch (reactive) syntax: watch expression :: handler
func (p *Parser) parseWatchStatement() Node {
	p.advance() // consume 'watch'

	watchExpr := p.parseExpression()

	// Handle :: delimiter
	if p.match(token.DOUBLE_COLON) {
		p.advance()
	}

	var handler Node
	if p.match(token.LBRACE) {
		p.advance()
		// Accept block of statements for watch
		stmts := p.parseBlock()
		handler = &ArrayLiteral{Elements: stmts}
	} else if p.match(token.LPAREN) {
		handler = p.parseTupleLiteral()
		if p.match(token.ARROW) {
			p.advance()
			handler = p.parseExpression()
		}
	}

	return &Statement{
		Type: "watch",
		Value: map[string]interface{}{
			"expression": watchExpr,
			"handler":    handler,
		},
	}
}

// parseOnStatement handles event listeners: on expression :: (params) -> body
func (p *Parser) parseOnStatement() Node {
	p.advance() // consume 'on'

	expr := p.parseExpression()

	// Handle :: delimiter
	if p.match(token.DOUBLE_COLON) {
		p.advance()
	}

	var handler Node
	if p.match(token.LPAREN) {
		handler = p.parseTupleLiteral()
		if p.match(token.ARROW) {
			p.advance()
			handler = p.parseStatement()
		}
	} else {
		handler = p.parseStatement()
	}

	return &Statement{
		Type: "on",
		Value: map[string]interface{}{
			"expression": expr,
			"handler":    handler,
		},
	}
}

// parseWithStatement handles context managers: with context { ... } :: { body }
func (p *Parser) parseWithStatement() Node {
	p.advance() // consume 'with'

	if !p.match(token.CONTEXT) {
		p.addError(p.current(), "expected 'context' after 'with'")
	} else {
		p.advance()
	}

	var contextMap Node
	if p.match(token.LBRACE) {
		contextMap = p.parseMapLiteral()
	}

	// Handle :: delimiter
	if p.match(token.DOUBLE_COLON) {
		p.advance()
	}

	var body []Node
	if p.match(token.LBRACE) {
		p.advance()
		body = p.parseBlock()
	}

	return &Statement{
		Type: "with",
		Value: map[string]interface{}{
			"context": contextMap,
			"body":    body,
		},
	}
}

// parseSnapshotStatement handles state persistence: snapshot agent.state -> file
func (p *Parser) parseSnapshotStatement() Node {
	p.advance() // consume 'snapshot'

	source := p.parseExpression()

	// Handle both -> and direct string filepath
	var target Node
	if p.match(token.ARROW) {
		p.advance()
		target = p.parseExpression()
	} else {
		// If no arrow, treat the source as "snapshot_target" command
		// This handles: snapshot AICoder.state (without ->)
		target = &Literal{Type: "string", Value: ""}
	}

	return &Statement{
		Type: "snapshot",
		Value: map[string]interface{}{
			"source": source,
			"target": target,
		},
	}
}

// parseRestoreStatement handles state restoration: restore agent from file
func (p *Parser) parseRestoreStatement() Node {
	p.advance() // consume 'restore'

	target := p.parseExpression()

	if !p.match(token.FROM) {
		p.addError(p.current(), "expected 'from' in restore statement")
	} else {
		p.advance()
	}

	source := p.parseExpression()

	return &Statement{
		Type: "restore",
		Value: map[string]interface{}{
			"source": source,
			"target": target,
		},
	}
}
