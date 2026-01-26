// parser.go
package parser

import (
	"fmt"
	"strings"
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
	DeclType string // "bind", "const", "craft", "fn", "struct", "agent", "task"
	Name     string
	Type     string
	Value    Node
	Params   []Parameter
	Body     []Node
	Fields   []Field // for struct and agent
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
	Init       Node
	Condition  Node
	Update     Node
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
	for p.match(token.NEWLINE) {
		p.advance()
	}
}

// isValidFieldName checks if a token type can be used as a field name in objects/agents
func isValidFieldName(tokenType token.TokenType) bool {
	// Allow these keywords as field names since they're commonly used in @agent blocks
	validKeywords := map[token.TokenType]bool{
		token.ROLE:            true,
		token.TOOLS:           true,
		token.MODEL:           true,
		token.MODE:            true,
		token.SYS_PROMPT:      true,
		token.PRINT:           true,
		token.TRACE:           true,
		token.PSEUDO:          true,
		token.TIMEOUT:         true,
		token.STAGE:           true,
		token.DEBUG:           true,
		token.WITH:            true,
		token.ALLOW:           true,
		token.CONTEXT:         true,
		token.MEMORY:          true,
		token.FLOW:            true,
		token.STRATEGY:        true,
		token.WINDOW:          true,
		token.ALERT_THRESHOLD: true,
	}
	return validKeywords[tokenType]
}

// Main parsing method
func (p *Parser) Parse() *ParseResult {
	program := &Program{Statements: []Node{}}

	for p.current().Type != token.EOF {
		// Skip newlines and comments
		if p.match(token.NEWLINE, token.COMMENT_LINE, token.COMMENT_MULTI) {
			p.advance()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		// Skip statement terminators
		if p.match(token.STATEMENT_END, token.NEWLINE) {
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
	p.skipNewlines()

	if p.current().Type == token.EOF {
		return nil
	}

	switch p.current().Type {
	case token.ALLOW:
		return p.parseAllowStatement()
	case token.AT_AGENT:
		return p.parseAgentDeclaration()
	case token.TASK:
		return p.parseTaskDeclaration()
	case token.BIND, token.CONST, token.CRAFT:
		return p.parseDeclaration()
	case token.INT_TYPE, token.FLOAT_TYPE, token.STR_TYPE, token.BOOL_TYPE:
		// Handle type-first declarations like "int x =: 10;"
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
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseAllowStatement() Node {
	p.advance() // consume 'allow'

	// allow pseudo(anno, trace, breakpoint, 15);
	functionName := ""
	if p.match(token.IDENTIFIER) {
		functionName = p.advance().Lexeme
	}

	var args []Node
	if p.match(token.LPAREN) {
		p.advance()
		for !p.match(token.RPAREN) && p.current().Type != token.EOF {
			args = append(args, p.parseExpression())
			if !p.match(token.RPAREN) {
				p.expect(token.COMMA)
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

		if p.match(token.RBRACE) {
			break
		}

		// Allow any identifier or keyword as field name
		fieldToken := p.current()
		var fieldName string
		if fieldToken.Type == token.IDENTIFIER || isValidFieldName(fieldToken.Type) {
			fieldName = fieldToken.Lexeme
			p.advance()
		} else {
			p.addError(fieldToken, "expected field name")
			break
		}

		if !p.match(token.COLON) {
			p.addError(p.current(), "expected ':' after field name")
			break
		}
		p.advance()

		value := p.parseExpression()

		fields = append(fields, Field{
			Name:  fieldName,
			Value: value,
		})

		if p.match(token.COMMA, token.STATEMENT_END, token.NEWLINE) {
			p.advance()
		}
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

		if p.match(token.RBRACE) {
			break
		}

		fieldToken, ok := p.expect(token.IDENTIFIER)
		if !ok {
			break
		}

		if !p.match(token.COLON) {
			p.addError(p.current(), "expected ':' after field name")
			break
		}
		p.advance()

		value := p.parseExpression()

		fields = append(fields, Field{
			Name:  fieldToken.Lexeme,
			Value: value,
		})

		if p.match(token.COMMA, token.STATEMENT_END, token.NEWLINE) {
			p.advance()
		}
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

	// Expect identifier
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
	// Handle "int x =: 10;" or "bind int x =: 10;" syntax
	typeToken := p.advance()
	varType := typeToken.Lexeme

	// Check if there's a bind/const/craft keyword before the type
	declType := "bind" // default

	nameToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	var value Node

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

func (p *Parser) parseSwitchStatement() Node {
	p.advance() // consume 'switch'

	// Parse the expression to switch on
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

			// Parse case value
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

			// Parse case body
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
			p.advance() // skip unexpected token
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

	nameToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	// Parse parameters
	if !p.match(token.LPAREN) {
		p.addError(p.current(), "expected '(' after function name")
		return nil
	}
	p.advance()

	params := []Parameter{}
	for !p.match(token.RPAREN) && p.current().Type != token.EOF {
		paramToken, ok := p.expect(token.IDENTIFIER)
		if !ok {
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

		params = append(params, Parameter{Name: paramToken.Lexeme, Type: paramType})

		if !p.match(token.RPAREN) {
			p.expect(token.COMMA)
		}
	}

	if !p.match(token.RPAREN) {
		p.addError(p.current(), "expected ')'")
		return nil
	}
	p.advance()

	// Parse optional return type
	returnType := ""
	if p.match(token.FAT_ARROW) {
		p.advance()
		if p.match(token.INT_TYPE, token.FLOAT_TYPE, token.STR_TYPE, token.BOOL_TYPE, token.IDENTIFIER) {
			typeToken := p.current()
			p.advance()
			returnType = typeToken.Lexeme
		}
	}

	// Parse function body
	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' for function body")
		return nil
	}
	p.advance()

	body := p.parseBlock()

	return &Declaration{
		DeclType: "fn",
		Name:     nameToken.Lexeme,
		Type:     returnType,
		Params:   params,
		Body:     body,
	}
}

func (p *Parser) parseAsyncStatement() Node {
	p.advance() // consume 'async'

	// async fn or async { }
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

	expr := p.parseExpression()

	return &AwaitExpression{Expression: expr}
}

func (p *Parser) parseEmitStatement() Node {
	p.advance() // consume 'emit'

	eventToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}

	var data Node
	if p.match(token.LBRACE) {
		data = p.parseMapLiteral()
	}

	return &EmitStatement{
		EventName: eventToken.Lexeme,
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

	// Parse handler (could be identifier => async { ... })
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

		// catch error { ... }
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

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' after if condition")
		return nil
	}
	p.advance()

	thenBody := p.parseBlock()

	elseBody := []Node{}
	if p.match(token.ELSE) {
		p.advance()
		if p.match(token.IF) {
			// else if
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

	if !p.match(token.LBRACE) {
		p.addError(p.current(), "expected '{' after while condition")
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

	// for i =: 0; i < 10; i++ { }
	// for item in collection { }
	// for i =: 0; i < 4; i++ concurrent { }

	// Check for 'in' style loop
	if p.match(token.IDENTIFIER) && p.peek(1).Type == token.IDENTIFIER {
		// for i in range { }
		loopVar := p.advance().Lexeme
		p.expect(token.IDENTIFIER) // skip 'in'

		condition = &BinaryOp{
			Left:     &Identifier{Name: loopVar},
			Operator: "in",
			Right:    p.parseExpression(),
		}
	} else {
		// Traditional for loop
		if !p.match(token.STATEMENT_END) {
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

	// Check for concurrent keyword
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
		p.skipNewlines()

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

		if p.match(token.STATEMENT_END, token.NEWLINE, token.COMMA) {
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
	// Check for configuration block syntax: identifier: { ... }
	if p.current().Type == token.IDENTIFIER && p.current().Type == token.IDENTIFIER {
		// Look ahead to check for colon and brace
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == token.COLON {
			if p.pos+2 < len(p.tokens) && p.tokens[p.pos+2].Type == token.LBRACE {
				// This is a configuration block
				nameToken := p.current()
				p.advance() // consume identifier
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

func (p *Parser) parseAssignment() Node {
	left := p.parseLogicalOr()

	if p.match(token.ASSIGN, token.BIND_ASSIGN, token.PLUS_ASSIGN, token.MINUS_ASSIGN,
		token.MULT_ASSIGN, token.DIV_ASSIGN, token.MOD_ASSIGN) {
		opToken := p.advance()
		right := p.parseExpression()
		return &BinaryOp{Left: left, Operator: opToken.Lexeme, Right: right}
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
			// Agent call: DataProcessor -> "command"
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

	case token.TRUE_KW, token.FALSE_KW:
		tok := p.advance()
		return &Literal{Type: "bool", Value: tok.Lexeme}

	case token.LPAREN:
		p.advance()
		expr := p.parseExpression()
		p.expect(token.RPAREN)
		return expr

	case token.LBRACKET:
		return p.parseArrayLiteral()

	case token.LBRACE:
		return p.parseMapLiteral()

	case token.ASYNC:
		return p.parseAsyncStatement()

	default:
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
		p.skipNewlines()
		if p.match(token.RBRACE) {
			break
		}

		key := p.parseExpression()
		p.expect(token.COLON)
		value := p.parseExpression()
		pairs = append(pairs, KeyValue{Key: key, Value: value})

		if p.match(token.COMMA, token.STATEMENT_END, token.NEWLINE) {
			p.advance()
		}

		if p.match(token.RBRACE) {
			break
		}
	}

	p.expect(token.RBRACE)
	return &MapLiteral{Pairs: pairs}
}

// AST to string representation
func FormatAST(node Node, indent int) string {
	prefix := strings.Repeat("  ", indent)

	switch n := node.(type) {
	case *Program:
		result := prefix + "Program\n"
		for _, stmt := range n.Statements {
			result += FormatAST(stmt, indent+1)
		}
		return result

	case *Declaration:
		result := prefix + fmt.Sprintf("Declaration (%s)\n", n.DeclType)
		result += prefix + fmt.Sprintf("  name: %s\n", n.Name)
		if n.Type != "" {
			result += prefix + fmt.Sprintf("  type: %s\n", n.Type)
		}
		if n.Value != nil {
			result += prefix + "  value:\n" + FormatAST(n.Value, indent+2)
		}
		if len(n.Params) > 0 {
			result += prefix + "  params:\n"
			for _, p := range n.Params {
				result += prefix + fmt.Sprintf("    %s: %s\n", p.Name, p.Type)
			}
		}
		if len(n.Fields) > 0 {
			result += prefix + "  fields:\n"
			for _, f := range n.Fields {
				result += prefix + fmt.Sprintf("    %s: %s\n", f.Name, f.Type)
				if f.Value != nil {
					result += FormatAST(f.Value, indent+3)
				}
			}
		}
		if len(n.Body) > 0 {
			result += prefix + "  body:\n"
			for _, stmt := range n.Body {
				result += FormatAST(stmt, indent+2)
			}
		}
		return result

	case *BinaryOp:
		result := prefix + fmt.Sprintf("BinaryOp (%s)\n", n.Operator)
		if n.Left != nil {
			result += prefix + "  left:\n" + FormatAST(n.Left, indent+2)
		}
		if n.Right != nil {
			result += prefix + "  right:\n" + FormatAST(n.Right, indent+2)
		}
		return result

	case *UnaryOp:
		result := prefix + fmt.Sprintf("UnaryOp (%s)\n", n.Operator)
		if n.Operand != nil {
			result += prefix + "  operand:\n" + FormatAST(n.Operand, indent+2)
		}
		return result

	case *Literal:
		return prefix + fmt.Sprintf("Literal (%s): %s\n", n.Type, n.Value)

	case *Identifier:
		return prefix + fmt.Sprintf("Identifier: %s\n", n.Name)

	case *IfStatement:
		result := prefix + "IfStatement\n"
		if n.Condition != nil {
			result += prefix + "  condition:\n" + FormatAST(n.Condition, indent+2)
		}
		result += prefix + "  then:\n"
		for _, stmt := range n.ThenBody {
			result += FormatAST(stmt, indent+2)
		}
		if len(n.ElseBody) > 0 {
			result += prefix + "  else:\n"
			for _, stmt := range n.ElseBody {
				result += FormatAST(stmt, indent+2)
			}
		}
		return result

	case *WhileStatement:
		result := prefix + "WhileStatement\n"
		if n.Condition != nil {
			result += prefix + "  condition:\n" + FormatAST(n.Condition, indent+2)
		}
		result += prefix + "  body:\n"
		for _, stmt := range n.Body {
			result += FormatAST(stmt, indent+2)
		}
		return result

	case *ForStatement:
		result := prefix + "ForStatement\n"
		if n.Concurrent {
			result = prefix + "ForStatement (concurrent)\n"
		}
		if n.Init != nil {
			result += prefix + "  init:\n" + FormatAST(n.Init, indent+2)
		}
		if n.Condition != nil {
			result += prefix + "  condition:\n" + FormatAST(n.Condition, indent+2)
		}
		if n.Update != nil {
			result += prefix + "  update:\n" + FormatAST(n.Update, indent+2)
		}
		result += prefix + "  body:\n"
		for _, stmt := range n.Body {
			result += FormatAST(stmt, indent+2)
		}
		return result

	case *SwitchStatement:
		result := prefix + "SwitchStatement\n"
		if n.Expression != nil {
			result += prefix + "  expression:\n" + FormatAST(n.Expression, indent+2)
		}
		result += prefix + "  cases:\n"
		for _, c := range n.Cases {
			result += prefix + "    case:\n"
			if c.Value != nil {
				result += prefix + "      value:\n" + FormatAST(c.Value, indent+4)
			}
			result += prefix + "      body:\n"
			for _, stmt := range c.Body {
				result += FormatAST(stmt, indent+4)
			}
		}
		if len(n.Default) > 0 {
			result += prefix + "  default:\n"
			for _, stmt := range n.Default {
				result += FormatAST(stmt, indent+2)
			}
		}
		return result

	case *ReturnStatement:
		result := prefix + "ReturnStatement\n"
		if n.Value != nil {
			result += FormatAST(n.Value, indent+1)
		}
		return result

	case *CallExpression:
		result := prefix + "CallExpression\n"
		if n.Function != nil {
			result += prefix + "  function:\n" + FormatAST(n.Function, indent+2)
		}
		if len(n.Arguments) > 0 {
			result += prefix + "  arguments:\n"
			for _, arg := range n.Arguments {
				result += FormatAST(arg, indent+2)
			}
		}
		return result

	case *ArrayLiteral:
		result := prefix + "ArrayLiteral\n"
		for _, elem := range n.Elements {
			result += FormatAST(elem, indent+1)
		}
		return result

	case *MapLiteral:
		result := prefix + "MapLiteral\n"
		for _, pair := range n.Pairs {
			result += prefix + "  pair:\n"
			if pair.Key != nil {
				result += prefix + "    key:\n" + FormatAST(pair.Key, indent+3)
			}
			if pair.Value != nil {
				result += prefix + "    value:\n" + FormatAST(pair.Value, indent+3)
			}
		}
		return result

	case *AsyncStatement:
		result := prefix + "AsyncStatement\n"
		for _, stmt := range n.Body {
			result += FormatAST(stmt, indent+1)
		}
		return result

	case *AwaitExpression:
		result := prefix + "AwaitExpression\n"
		if n.Expression != nil {
			result += FormatAST(n.Expression, indent+1)
		}
		return result

	case *EmitStatement:
		result := prefix + fmt.Sprintf("EmitStatement (%s)\n", n.EventName)
		if n.Data != nil {
			result += FormatAST(n.Data, indent+1)
		}
		return result

	case *ListenStatement:
		result := prefix + fmt.Sprintf("ListenStatement (%s)\n", n.EventName)
		if n.Handler != nil {
			result += FormatAST(n.Handler, indent+1)
		}
		return result

	case *ConfigBlock:
		result := prefix + fmt.Sprintf("ConfigBlock (%s)\n", n.Name)
		if n.Value != nil {
			result += FormatAST(n.Value, indent+1)
		}
		return result

	case *Statement:
		result := prefix + fmt.Sprintf("Statement (%s)\n", n.Type)
		return result

	default:
		return prefix + "Unknown node\n"
	}
}

// Parse from source code directly
func ParseFromSource(source string) *ParseResult {
	// Create lexer
	lex := lexer.New(source)
	tokens := lex.Tokenize()

	// Create parser
	parser := NewParser(tokens)
	return parser.Parse()
}
