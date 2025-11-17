// lexer/lexer.go
package lexer

import (
	"synta/token"
	"unicode"
)

type Lexer struct {
	input  string
	pos    int
	line   int
	column int
	tokens []token.Token
}

func New(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
		tokens: []token.Token{},
	}
}

func (l *Lexer) peek(offset int) byte {
	pos := l.pos + offset
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

func (l *Lexer) advance() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return ch
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) && l.input[l.pos] != '\n' {
		l.advance()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	for l.pos < len(l.input) && (unicode.IsLetter(rune(l.input[l.pos])) ||
		unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '_') {
		l.advance()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	start := l.pos
	tokenType := token.INTEGER

	for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
		l.advance()
	}

	if l.pos < len(l.input) && l.input[l.pos] == '.' &&
		l.pos+1 < len(l.input) && unicode.IsDigit(rune(l.input[l.pos+1])) {
		tokenType = token.FLOAT
		l.advance() // consume '.'
		for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
			l.advance()
		}
	}

	return l.input[start:l.pos], tokenType
}

func (l *Lexer) readString() string {
	quote := l.advance() // consume opening quote
	start := l.pos

	for l.pos < len(l.input) && l.input[l.pos] != quote {
		if l.input[l.pos] == '\\' {
			l.advance() // skip escape char
		}
		l.advance()
	}

	str := l.input[start:l.pos]
	if l.pos < len(l.input) {
		l.advance() // consume closing quote
	}
	return str
}

func (l *Lexer) skipComment() {
	if l.peek(0) == '/' && l.peek(1) == '/' {
		// Single-line comment
		for l.pos < len(l.input) && l.input[l.pos] != '\n' {
			l.advance()
		}
	} else if l.peek(0) == '/' && l.peek(1) == '~' {
		// Multi-line comment
		l.advance() // /
		l.advance() // ~
		for l.pos < len(l.input)-1 {
			if l.input[l.pos] == '~' && l.input[l.pos+1] == '/' {
				l.advance()
				l.advance()
				break
			}
			l.advance()
		}
	}
}

func (l *Lexer) addToken(tokenType token.TokenType, lexeme string, line, column int) {
	l.tokens = append(l.tokens, token.Token{
		Type:   tokenType,
		Lexeme: lexeme,
		Line:   line,
		Column: column,
	})
}

func (l *Lexer) Tokenize() []token.Token {
	for l.pos < len(l.input) {
		l.skipWhitespace()

		if l.pos >= len(l.input) {
			break
		}

		line, column := l.line, l.column
		ch := l.peek(0)

		// Comments
		if ch == '/' && (l.peek(1) == '/' || l.peek(1) == '~') {
			l.skipComment()
			continue
		}

		// Identifiers and keywords
		if unicode.IsLetter(rune(ch)) || ch == '_' {
			ident := l.readIdentifier()
			l.addToken(token.LookupIdent(ident), ident, line, column)
			continue
		}

		// Numbers
		if unicode.IsDigit(rune(ch)) {
			num, tokenType := l.readNumber()
			l.addToken(tokenType, num, line, column)
			continue
		}

		// Strings
		if ch == '"' || ch == '\'' {
			str := l.readString()
			l.addToken(token.STRING, str, line, column)
			continue
		}

		// Operators and delimiters
		switch ch {
		case '+':
			l.advance()
			if l.peek(0) == '+' {
				l.advance()
				l.addToken(token.INCREMENT, "++", line, column)
			} else if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.PLUS_ASSIGN, "+=", line, column)
			} else {
				l.addToken(token.PLUS, "+", line, column)
			}
		case '-':
			l.advance()
			if l.peek(0) == '-' {
				l.advance()
				l.addToken(token.DECREMENT, "--", line, column)
			} else if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.MINUS_ASSIGN, "-=", line, column)
			} else {
				l.addToken(token.MINUS, "-", line, column)
			}
		case '*':
			l.advance()
			if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.MULT_ASSIGN, "*=", line, column)
			} else {
				l.addToken(token.MULTIPLY, "*", line, column)
			}
		case '/':
			l.advance()
			if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.DIV_ASSIGN, "/=", line, column)
			} else {
				l.addToken(token.DIVIDE, "/", line, column)
			}
		case '%':
			l.advance()
			if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.MOD_ASSIGN, "%=", line, column)
			} else {
				l.addToken(token.MODULO, "%", line, column)
			}
		case '=':
			l.advance()
			if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.EQ, "==", line, column)
			} else {
				l.addToken(token.ASSIGN, "=", line, column)
			}
		case ':':
			l.advance()
			if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.BIND_ASSIGN, ":=", line, column)
			} else {
				l.addToken(token.COLON, ":", line, column)
			}
		case '!':
			l.advance()
			if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.NEQ, "!=", line, column)
			} else {
				l.addToken(token.NOT, "!", line, column)
			}
		case '<':
			l.advance()
			if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.LTE, "<=", line, column)
			} else {
				l.addToken(token.LT, "<", line, column)
			}
		case '>':
			l.advance()
			if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.GTE, ">=", line, column)
			} else {
				l.addToken(token.GT, ">", line, column)
			}
		case '&':
			l.advance()
			if l.peek(0) == '&' {
				l.advance()
				l.addToken(token.AND, "&&", line, column)
			} else {
				l.addToken(token.AMPERSAND, "&", line, column)
			}
		case '|':
			l.advance()
			if l.peek(0) == '|' {
				l.advance()
				l.addToken(token.OR, "||", line, column)
			}
		case '^':
			l.advance()
			l.addToken(token.BITWISE_XOR, "^", line, column)
		case '(':
			l.advance()
			l.addToken(token.LPAREN, "(", line, column)
		case ')':
			l.advance()
			l.addToken(token.RPAREN, ")", line, column)
		case '[':
			l.advance()
			l.addToken(token.LBRACKET, "[", line, column)
		case ']':
			l.advance()
			l.addToken(token.RBRACKET, "]", line, column)
		case '{':
			l.advance()
			l.addToken(token.LBRACE, "{", line, column)
		case '}':
			l.advance()
			l.addToken(token.RBRACE, "}", line, column)
		case ';':
			l.advance()
			l.addToken(token.SEMICOLON, ";", line, column)
		case ',':
			l.advance()
			l.addToken(token.COMMA, ",", line, column)
		case '.':
			l.advance()
			l.addToken(token.DOT, ".", line, column)
		case '\n':
			l.advance()
			l.addToken(token.NEWLINE, "\\n", line, column)
		default:
			l.advance()
			l.addToken(token.ILLEGAL, string(ch), line, column)
		}
	}

	l.addToken(token.EOF, "", l.line, l.column)
	return l.tokens
}
