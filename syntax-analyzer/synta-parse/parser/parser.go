package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ============================================================================
// AST Nodes
// ============================================================================

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement nodes
type Statement interface {
	Node
	statementNode()
}

// Expression nodes
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	return out.String()
}

// Identifier
type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Lexeme }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral
type IntegerLiteral struct {
	Token Token
	Value string
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Lexeme }
func (il *IntegerLiteral) String() string       { return il.Value }

// FloatLiteral
type FloatLiteral struct {
	Token Token
	Value string
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Lexeme }
func (fl *FloatLiteral) String() string       { return fl.Value }

// StringLiteral
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Lexeme }
func (sl *StringLiteral) String() string       { return fmt.Sprintf("\"%s\"", sl.Value) }

// BooleanLiteral
type BooleanLiteral struct {
	Token Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Lexeme }
func (bl *BooleanLiteral) String() string       { return bl.Token.Lexeme }

// BindStatement: bind x := 10
type BindStatement struct {
	Token Token
	Name  *Identifier
	Value Expression
}

func (bs *BindStatement) statementNode()       {}
func (bs *BindStatement) TokenLiteral() string { return bs.Token.Lexeme }
func (bs *BindStatement) String() string {
	return fmt.Sprintf("bind %s := %s", bs.Name.String(), bs.Value.String())
}

// ConstStatement: const PI := 3.14
type ConstStatement struct {
	Token Token
	Name  *Identifier
	Value Expression
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Lexeme }
func (cs *ConstStatement) String() string {
	return fmt.Sprintf("const %s := %s", cs.Name.String(), cs.Value.String())
}

// AssignStatement: x =: 20
type AssignStatement struct {
	Token Token
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Lexeme }
func (as *AssignStatement) String() string {
	return fmt.Sprintf("%s =: %s", as.Name.String(), as.Value.String())
}

// ReturnStatement: return x
type ReturnStatement struct {
	Token       Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Lexeme }
func (rs *ReturnStatement) String() string {
	if rs.ReturnValue != nil {
		return fmt.Sprintf("return %s", rs.ReturnValue.String())
	}
	return "return"
}

// ExpressionStatement
type ExpressionStatement struct {
	Token      Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Lexeme }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement
type BlockStatement struct {
	Token      Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Lexeme }
func (bs *BlockStatement) String() string {
	var out strings.Builder
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString("  ")
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	out.WriteString("}")
	return out.String()
}

// IfStatement
type IfStatement struct {
	Token       Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative Statement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Lexeme }
func (is *IfStatement) String() string {
	var out strings.Builder
	out.WriteString("if ")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())
	if is.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(is.Alternative.String())
	}
	return out.String()
}

// WhileStatement
type WhileStatement struct {
	Token     Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Lexeme }
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("while %s %s", ws.Condition.String(), ws.Body.String())
}

// ForStatement
type ForStatement struct {
	Token    Token
	Variable *Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Lexeme }
func (fs *ForStatement) String() string {
	return fmt.Sprintf("for %s in %s %s", fs.Variable.String(), fs.Iterable.String(), fs.Body.String())
}

// FunctionStatement
type FunctionStatement struct {
	Token      Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Lexeme }
func (fs *FunctionStatement) String() string {
	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("fn %s(%s) %s", fs.Name.String(), strings.Join(params, ", "), fs.Body.String())
}

// PrintStatement
type PrintStatement struct {
	Token      Token
	Expression Expression
}

func (ps *PrintStatement) statementNode()       {}
func (ps *PrintStatement) TokenLiteral() string { return ps.Token.Lexeme }
func (ps *PrintStatement) String() string {
	return fmt.Sprintf("print %s", ps.Expression.String())
}

// PrefixExpression
type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Lexeme }
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

// InfixExpression
type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Lexeme }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

// CallExpression
type CallExpression struct {
	Token     Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Lexeme }
func (ce *CallExpression) String() string {
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return fmt.Sprintf("%s(%s)", ce.Function.String(), strings.Join(args, ", "))
}

// ArrayLiteral
type ArrayLiteral struct {
	Token    Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Lexeme }
func (al *ArrayLiteral) String() string {
	elements := []string{}
	for _, e := range al.Elements {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// IndexExpression
type IndexExpression struct {
	Token Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Lexeme }
func (ie *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left.String(), ie.Index.String())
}

// GetNodePosition returns the token line and column for a node if available
func GetNodePosition(n Node) (line int, column int) {
	switch v := n.(type) {
	case *Identifier:
		return v.Token.Line, v.Token.Column
	case *IntegerLiteral:
		return v.Token.Line, v.Token.Column
	case *FloatLiteral:
		return v.Token.Line, v.Token.Column
	case *StringLiteral:
		return v.Token.Line, v.Token.Column
	case *BooleanLiteral:
		return v.Token.Line, v.Token.Column
	case *BindStatement:
		return v.Token.Line, v.Token.Column
	case *ConstStatement:
		return v.Token.Line, v.Token.Column
	case *AssignStatement:
		return v.Token.Line, v.Token.Column
	case *ReturnStatement:
		return v.Token.Line, v.Token.Column
	case *ExpressionStatement:
		return v.Token.Line, v.Token.Column
	case *BlockStatement:
		return v.Token.Line, v.Token.Column
	case *IfStatement:
		return v.Token.Line, v.Token.Column
	case *WhileStatement:
		return v.Token.Line, v.Token.Column
	case *ForStatement:
		return v.Token.Line, v.Token.Column
	case *FunctionStatement:
		return v.Token.Line, v.Token.Column
	case *PrintStatement:
		return v.Token.Line, v.Token.Column
	case *PrefixExpression:
		return v.Token.Line, v.Token.Column
	case *InfixExpression:
		return v.Token.Line, v.Token.Column
	case *CallExpression:
		return v.Token.Line, v.Token.Column
	case *ArrayLiteral:
		return v.Token.Line, v.Token.Column
	case *IndexExpression:
		return v.Token.Line, v.Token.Column
	default:
		return 0, 0
	}
}

// ============================================================================
// Parser
// ============================================================================

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[TokenType]int{
	EQ:       EQUALS,
	NEQ:      EQUALS,
	LT:       LESSGREATER,
	GT:       LESSGREATER,
	LTE:      LESSGREATER,
	GTE:      LESSGREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	DIVIDE:   PRODUCT,
	MULTIPLY: PRODUCT,
	MODULO:   PRODUCT,
	LPAREN:   CALL,
	LBRACKET: INDEX,
	AND:      LOWEST + 1,
	OR:       LOWEST + 1,
}

type ParseError struct {
	Tok Token
	Msg string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("Line %d:%d: %s (%s)", e.Tok.Line, e.Tok.Column, e.Msg, e.Tok.Lexeme)
}

type Parser struct {
	tokens   []Token
	pos      int
	curToken Token
	errors   []error
	debugLog []string

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

func New(tokens []Token) *Parser {
	p := &Parser{
		tokens:         tokens,
		pos:            0,
		errors:         []error{},
		debugLog:       []string{},
		prefixParseFns: make(map[TokenType]prefixParseFn),
		infixParseFns:  make(map[TokenType]infixParseFn),
	}

	// Register prefix parse functions
	p.registerPrefix(IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(FLOAT, p.parseFloatLiteral)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(NOT, p.parsePrefixExpression)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)
	p.registerPrefix(LBRACKET, p.parseArrayLiteral)

	// Register infix parse functions
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(DIVIDE, p.parseInfixExpression)
	p.registerInfix(MULTIPLY, p.parseInfixExpression)
	p.registerInfix(MODULO, p.parseInfixExpression)
	p.registerInfix(EQ, p.parseInfixExpression)
	p.registerInfix(NEQ, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(LTE, p.parseInfixExpression)
	p.registerInfix(GTE, p.parseInfixExpression)
	p.registerInfix(AND, p.parseInfixExpression)
	p.registerInfix(OR, p.parseInfixExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)
	p.registerInfix(LBRACKET, p.parseIndexExpression)

	if len(tokens) > 0 {
		p.curToken = tokens[0]
	}

	return p
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Parse() (*Program, []error, []string) {
	p.log("Starting parse")
	program := &Program{
		Statements: []Statement{},
	}

	for p.curToken.Type != EOF {
		// Skip newlines and comments
		if p.curToken.Type == NEWLINE || p.curToken.Type == COMMENT_LINE || p.curToken.Type == COMMENT_MULTI {
			p.advance()
			continue
		}

		// Skip noise words
		if p.curToken.Type == DO || p.curToken.Type == PLEASE || p.curToken.Type == MAYBE {
			p.advance()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
			p.log(fmt.Sprintf("Parsed statement: %T", stmt))
		}

		p.advance()
	}

	p.log(fmt.Sprintf("Parsing complete. %d statements parsed", len(program.Statements)))
	return program, p.errors, p.debugLog
}

func (p *Parser) parseStatement() Statement {
	p.log(fmt.Sprintf("Parsing statement at token: %s", p.curToken.Type.String()))

	switch p.curToken.Type {
	case BIND, LET:
		return p.parseBindStatement()
	case CONST:
		return p.parseConstStatement()
	case RETURN:
		return p.parseReturnStatement()
	case IF:
		return p.parseIfStatement()
	case WHILE:
		return p.parseWhileStatement()
	case FOR:
		return p.parseForStatement()
	case FN:
		return p.parseFunctionStatement()
	case PRINT:
		return p.parsePrintStatement()
	case IDENTIFIER:
		// Check if this is an assignment
		if p.peekToken().Type == ASSIGN {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBindStatement() Statement {
	stmt := &BindStatement{Token: p.curToken}

	p.advance()
	if p.curToken.Type != IDENTIFIER {
		p.error(p.curToken, "expected identifier after 'bind'")
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Lexeme}

	p.advance()
	if p.curToken.Type != BIND_ASSIGN {
		p.error(p.curToken, "expected ':=' after identifier")
		return nil
	}

	p.advance()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseConstStatement() Statement {
	stmt := &ConstStatement{Token: p.curToken}

	p.advance()
	if p.curToken.Type != IDENTIFIER {
		p.error(p.curToken, "expected identifier after 'const'")
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Lexeme}

	p.advance()
	if p.curToken.Type != BIND_ASSIGN {
		p.error(p.curToken, "expected ':=' after identifier")
		return nil
	}

	p.advance()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseAssignStatement() Statement {
	stmt := &AssignStatement{}
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Lexeme}

	p.advance() // Move to =:
	stmt.Token = p.curToken

	p.advance() // Move to value
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseReturnStatement() Statement {
	stmt := &ReturnStatement{Token: p.curToken}

	p.advance()

	// Return can be empty
	if p.curToken.Type == NEWLINE || p.curToken.Type == SEMICOLON || p.curToken.Type == EOF {
		return stmt
	}

	stmt.ReturnValue = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseIfStatement() Statement {
	stmt := &IfStatement{Token: p.curToken}

	p.advance()
	stmt.Condition = p.parseExpression(LOWEST)

	p.advance()
	if p.curToken.Type != LBRACE {
		p.error(p.curToken, "expected '{' after if condition")
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	// Check for elif or else
	if p.peekToken().Type == ELIF {
		p.advance()
		stmt.Alternative = p.parseIfStatement() // Recursive for elif
	} else if p.peekToken().Type == ELSE {
		p.advance()
		p.advance()
		if p.curToken.Type == LBRACE {
			stmt.Alternative = p.parseBlockStatement()
		}
	}

	return stmt
}

func (p *Parser) parseWhileStatement() Statement {
	stmt := &WhileStatement{Token: p.curToken}

	p.advance()
	stmt.Condition = p.parseExpression(LOWEST)

	p.advance()
	if p.curToken.Type != LBRACE {
		p.error(p.curToken, "expected '{' after while condition")
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() Statement {
	stmt := &ForStatement{Token: p.curToken}

	p.advance()
	if p.curToken.Type != IDENTIFIER {
		p.error(p.curToken, "expected identifier after 'for'")
		return nil
	}

	stmt.Variable = &Identifier{Token: p.curToken, Value: p.curToken.Lexeme}

	p.advance()
	// Expect 'in' keyword
	if p.curToken.Type != IDENTIFIER || p.curToken.Lexeme != "in" {
		p.error(p.curToken, "expected 'in' after for variable")
		return nil
	}

	p.advance()
	stmt.Iterable = p.parseExpression(LOWEST)

	p.advance()
	if p.curToken.Type != LBRACE {
		p.error(p.curToken, "expected '{' after for iterable")
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionStatement() Statement {
	stmt := &FunctionStatement{Token: p.curToken}

	p.advance()
	if p.curToken.Type != IDENTIFIER {
		p.error(p.curToken, "expected function name")
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Lexeme}

	p.advance()
	if p.curToken.Type != LPAREN {
		p.error(p.curToken, "expected '(' after function name")
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	p.advance()
	if p.curToken.Type != LBRACE {
		p.error(p.curToken, "expected '{' after function parameters")
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}

	p.advance()
	if p.curToken.Type == RPAREN {
		return identifiers
	}

	identifiers = append(identifiers, &Identifier{
		Token: p.curToken,
		Value: p.curToken.Lexeme,
	})

	for p.peekToken().Type == COMMA {
		p.advance()
		p.advance()
		identifiers = append(identifiers, &Identifier{
			Token: p.curToken,
			Value: p.curToken.Lexeme,
		})
	}

	p.advance()
	if p.curToken.Type != RPAREN {
		p.error(p.curToken, "expected ')' after function parameters")
		return nil
	}

	return identifiers
}

func (p *Parser) parsePrintStatement() Statement {
	stmt := &PrintStatement{Token: p.curToken}

	p.advance()
	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.advance()

	for p.curToken.Type != RBRACE && p.curToken.Type != EOF {
		// Skip newlines and comments
		if p.curToken.Type == NEWLINE || p.curToken.Type == COMMENT_LINE || p.curToken.Type == COMMENT_MULTI {
			p.advance()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.advance()
	}

	return block
}

func (p *Parser) parseExpressionStatement() Statement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.error(p.curToken, fmt.Sprintf("no prefix parse function for %s", p.curToken.Type.String()))
		return nil
	}

	leftExp := prefix()

	for p.peekToken().Type != SEMICOLON && p.peekToken().Type != NEWLINE && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken().Type]
		if infix == nil {
			return leftExp
		}

		p.advance()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Lexeme}
}

func (p *Parser) parseIntegerLiteral() Expression {
	return &IntegerLiteral{Token: p.curToken, Value: p.curToken.Lexeme}
}

func (p *Parser) parseFloatLiteral() Expression {
	return &FloatLiteral{Token: p.curToken, Value: p.curToken.Lexeme}
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Lexeme}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Lexeme,
	}

	p.advance()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Lexeme,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.advance()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.advance()
	exp := p.parseExpression(LOWEST)
	p.advance()
	if p.curToken.Type != RPAREN {
		p.error(p.curToken, "expected ')' after grouped expression")
		return nil
	}
	return exp
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseArrayLiteral() Expression {
	array := &ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(RBRACKET)
	return array
}

func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}

	p.advance()
	exp.Index = p.parseExpression(LOWEST)

	p.advance()
	if p.curToken.Type != RBRACKET {
		p.error(p.curToken, "expected ']'")
		return nil
	}

	return exp
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}

	p.advance()
	if p.curToken.Type == end {
		return list
	}

	list = append(list, p.parseExpression(LOWEST))

	for p.peekToken().Type == COMMA {
		p.advance()
		p.advance()
		list = append(list, p.parseExpression(LOWEST))
	}

	p.advance()
	if p.curToken.Type != end {
		p.error(p.curToken, fmt.Sprintf("expected %s", end.String()))
		return nil
	}

	return list
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens)-1 {
		p.pos++
		p.curToken = p.tokens[p.pos]
	}
}

func (p *Parser) peekToken() Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return Token{Type: EOF}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken().Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) error(tok Token, message string) {
	err := ParseError{Tok: tok, Msg: message}
	p.errors = append(p.errors, err)
	p.log(fmt.Sprintf("ERROR: %s", err.Error()))
}

func (p *Parser) log(message string) {
	p.debugLog = append(p.debugLog, message)
}

// ============================================================================
// Utilities
// ============================================================================

// LoadTokens reads a JSON token file and unmarshals into []Token
func LoadTokens(path string) ([]Token, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading token file: %w", err)
	}
	var tokens []Token
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, fmt.Errorf("error parsing token file: %w", err)
	}
	return tokens, nil
}

// WriteDebugLog writes debug lines to a file
func WriteDebugLog(path string, debug []string) error {
	content := ""
	for _, l := range debug {
		content += l + "\n"
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// WriteErrors writes parse errors to a file
func WriteErrors(path string, errors []error) error {
	if len(errors) == 0 {
		return nil
	}
	content := ""
	for _, e := range errors {
		content += e.Error() + "\n"
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// ============================================================================
// Tree Generation
// ============================================================================

// GeneratePrettyTree creates a pretty-printed tree representation
func GeneratePrettyTree(program *Program) string {
	var sb strings.Builder
	sb.WriteString("Program\n")
	for i, stmt := range program.Statements {
		sb.WriteString(fmt.Sprintf("├── Statement %d: %T\n", i+1, stmt))
		sb.WriteString(prettyPrintNode(stmt, "│   ", i == len(program.Statements)-1))
	}
	return sb.String()
}

func prettyPrintNode(node Node, prefix string, isLast bool) string {
	var sb strings.Builder

	switch n := node.(type) {
	case *BindStatement:
		sb.WriteString(fmt.Sprintf("%s├── Name: %s\n", prefix, n.Name.Value))
		sb.WriteString(fmt.Sprintf("%s└── Value: %s\n", prefix, n.Value.String()))

	case *ConstStatement:
		sb.WriteString(fmt.Sprintf("%s├── Name: %s\n", prefix, n.Name.Value))
		sb.WriteString(fmt.Sprintf("%s└── Value: %s\n", prefix, n.Value.String()))

	case *AssignStatement:
		sb.WriteString(fmt.Sprintf("%s├── Name: %s\n", prefix, n.Name.Value))
		sb.WriteString(fmt.Sprintf("%s└── Value: %s\n", prefix, n.Value.String()))

	case *ReturnStatement:
		if n.ReturnValue != nil {
			sb.WriteString(fmt.Sprintf("%s└── Value: %s\n", prefix, n.ReturnValue.String()))
		}

	case *IfStatement:
		sb.WriteString(fmt.Sprintf("%s├── Condition: %s\n", prefix, n.Condition.String()))
		sb.WriteString(fmt.Sprintf("%s├── Consequence: %d statements\n", prefix, len(n.Consequence.Statements)))
		if n.Alternative != nil {
			sb.WriteString(fmt.Sprintf("%s└── Alternative: %T\n", prefix, n.Alternative))
		}

	case *WhileStatement:
		sb.WriteString(fmt.Sprintf("%s├── Condition: %s\n", prefix, n.Condition.String()))
		sb.WriteString(fmt.Sprintf("%s└── Body: %d statements\n", prefix, len(n.Body.Statements)))

	case *ForStatement:
		sb.WriteString(fmt.Sprintf("%s├── Variable: %s\n", prefix, n.Variable.Value))
		sb.WriteString(fmt.Sprintf("%s├── Iterable: %s\n", prefix, n.Iterable.String()))
		sb.WriteString(fmt.Sprintf("%s└── Body: %d statements\n", prefix, len(n.Body.Statements)))

	case *FunctionStatement:
		params := []string{}
		for _, p := range n.Parameters {
			params = append(params, p.Value)
		}
		sb.WriteString(fmt.Sprintf("%s├── Name: %s\n", prefix, n.Name.Value))
		sb.WriteString(fmt.Sprintf("%s├── Parameters: [%s]\n", prefix, strings.Join(params, ", ")))
		sb.WriteString(fmt.Sprintf("%s└── Body: %d statements\n", prefix, len(n.Body.Statements)))

	case *PrintStatement:
		sb.WriteString(fmt.Sprintf("%s└── Expression: %s\n", prefix, n.Expression.String()))

	case *ExpressionStatement:
		sb.WriteString(fmt.Sprintf("%s└── Expression: %s\n", prefix, n.Expression.String()))
	}

	return sb.String()
}

// GenerateCompactTree creates a compact one-line-per-statement tree
func GenerateCompactTree(program *Program) string {
	var sb strings.Builder
	for i, stmt := range program.Statements {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, stmt.String()))
	}
	return sb.String()
}

// GenerateDetailedTree creates a detailed tree with type information
func GenerateDetailedTree(program *Program) string {
	var sb strings.Builder
	sb.WriteString("=== DETAILED PARSE TREE ===\n\n")
	sb.WriteString(fmt.Sprintf("Total Statements: %d\n\n", len(program.Statements)))

	for i, stmt := range program.Statements {
		sb.WriteString(fmt.Sprintf("Statement %d:\n", i+1))
		sb.WriteString(fmt.Sprintf("  Type: %T\n", stmt))
		sb.WriteString(fmt.Sprintf("  Token: %s\n", stmt.TokenLiteral()))
		sb.WriteString(fmt.Sprintf("  String: %s\n", stmt.String()))

		// Add type-specific details
		switch n := stmt.(type) {
		case *BindStatement:
			sb.WriteString(fmt.Sprintf("  Name: %s\n", n.Name.Value))
			sb.WriteString(fmt.Sprintf("  Value Type: %T\n", n.Value))

		case *FunctionStatement:
			sb.WriteString(fmt.Sprintf("  Function Name: %s\n", n.Name.Value))
			sb.WriteString(fmt.Sprintf("  Parameter Count: %d\n", len(n.Parameters)))
			sb.WriteString(fmt.Sprintf("  Body Statements: %d\n", len(n.Body.Statements)))

		case *IfStatement:
			sb.WriteString(fmt.Sprintf("  Condition Type: %T\n", n.Condition))
			sb.WriteString(fmt.Sprintf("  Consequence Statements: %d\n", len(n.Consequence.Statements)))
			sb.WriteString(fmt.Sprintf("  Has Alternative: %v\n", n.Alternative != nil))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}
