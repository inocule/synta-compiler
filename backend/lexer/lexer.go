// lexer.go
// lexical analyzer

package lexer

import (
	"synta-compiler/token"
	"unicode"
)

type Lexer struct {
	input  string
	pos    int
	line   int
	column int
	tokens []token.Token
}

// Single-character operator mapping
var singleCharOps = map[byte]token.TokenType{
	'(': token.LPAREN, ')': token.RPAREN, '[': token.LBRACKET, ']': token.RBRACKET,
	'{': token.LBRACE, '}': token.RBRACE, ',': token.COMMA, '^': token.BITWISE_XOR,
	';': token.STATEMENT_END, '$': token.DOLLAR,
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

	// Read integer part
	for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
		l.advance()
	}

	// Check for decimal point
	if l.pos < len(l.input) && l.input[l.pos] == '.' &&
		l.pos+1 < len(l.input) && unicode.IsDigit(rune(l.input[l.pos+1])) {
		tokenType = token.FLOAT
		l.advance() // consume '.'

		// Read fractional part
		for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
			l.advance()
		}
	}

	// Handle scientific notation or time suffixes (s, m, h)
	if l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == 's' || ch == 'm' || ch == 'h' {
			l.advance() // consume time suffix
		}
	}

	return l.input[start:l.pos], tokenType
}

func (l *Lexer) readString() string {
	quote := l.advance() // consume opening quote

	// Check for triple-quoted string
	isTripleQuoted := false
	if l.pos+1 < len(l.input) && l.input[l.pos] == quote && l.input[l.pos+1] == quote {
		isTripleQuoted = true
		l.advance() // consume second quote
		l.advance() // consume third quote
	}

	start := l.pos

	if isTripleQuoted {
		// For triple-quoted strings, look for closing """
		for l.pos < len(l.input) {
			if l.pos+2 < len(l.input) && l.input[l.pos] == quote && l.input[l.pos+1] == quote && l.input[l.pos+2] == quote {
				str := l.input[start:l.pos]
				l.advance() // consume first closing quote
				l.advance() // consume second closing quote
				l.advance() // consume third closing quote
				return str
			}
			l.advance()
		}
	} else {
		// For regular strings, look for single closing quote
		for l.pos < len(l.input) && l.input[l.pos] != quote {
			if l.input[l.pos] == '\\' && l.pos+1 < len(l.input) {
				l.advance() // consume backslash
				l.advance() // consume escaped character
			} else {
				l.advance()
			}
		}

		str := l.input[start:l.pos]
		if l.pos < len(l.input) {
			l.advance() // consume closing quote
		}
		return str
	}

	// Unclosed triple-quoted string
	return l.input[start:l.pos]
}

// Synta single-line comment: !>
func (l *Lexer) readLineComment() string {
	// consume the comment marker '!>' then capture the comment text
	l.advance() // !
	l.advance() // >

	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.advance()
	}

	return l.input[start:l.pos]
}

// Synta multi-line comment: <! !>
func (l *Lexer) readMultiComment() string {
	start := l.pos
	l.advance() // <
	l.advance() // !

	for l.pos < len(l.input) {
		if l.pos < len(l.input)-1 && l.input[l.pos] == '!' && l.input[l.pos+1] == '>' {
			l.advance() // !
			l.advance() // >
			return l.input[start:l.pos]
		}
		l.advance()
	}

	// If we reach here, comment wasn't closed properly
	return l.input[start:l.pos]
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

		// Handle Synta comments first (must be checked before < and ! operators)
		if ch == '<' && l.peek(1) == '!' {
			text := l.readMultiComment()
			l.addToken(token.COMMENT_MULTI, text, line, column)
			continue
		}

		if ch == '!' && l.peek(1) == '>' {
			text := l.readLineComment()
			l.addToken(token.COMMENT_LINE, text, line, column)
			continue
		}

		// Handle @ decorators
		if ch == '@' {
			l.advance()
			if unicode.IsLetter(rune(l.peek(0))) {
				ident := l.readIdentifier()
				switch ident {
				case "agent":
					l.addToken(token.AT_AGENT, "@agent", line, column)
				case "task":
					l.addToken(token.AT_TASK, "@task", line, column)
				case "step":
					l.addToken(token.AT_STEP, "@step", line, column)
				case "intent":
					l.addToken(token.AT_INTENT, "@intent", line, column)
				case "explain":
					l.addToken(token.AT_EXPLAIN, "@explain", line, column)
				case "allow":
					l.addToken(token.DECORATOR, "@allow", line, column)
				default:
					l.addToken(token.DECORATOR, "@"+ident, line, column)
				}
			} else {
				l.addToken(token.ILLEGAL, "@", line, column)
			}
			continue
		}

		// Handle identifiers and keywords
		if unicode.IsLetter(rune(ch)) || ch == '_' {
			ident := l.readIdentifier()
			l.addToken(token.LookupIdent(ident), ident, line, column)
			continue
		}

		// Handle numbers
		if unicode.IsDigit(rune(ch)) {
			num, tokenType := l.readNumber()
			l.addToken(tokenType, num, line, column)
			continue
		}

		// Handle strings
		if ch == '"' || ch == '\'' {
			str := l.readString()
			l.addToken(token.STRING, str, line, column)
			continue
		}

		// Handle operators and delimiters
		if tok, ok := singleCharOps[ch]; ok {
			l.advance()
			l.addToken(tok, string(ch), line, column)
			continue
		}

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
			} else if l.peek(0) == '>' {
				l.advance()
				l.addToken(token.ARROW, "->", line, column)
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
			} else if l.peek(0) == ':' {
				l.advance()
				l.addToken(token.ASSIGN, "=:", line, column)
			} else if l.peek(0) == '>' {
				l.advance()
				l.addToken(token.FAT_ARROW, "=>", line, column)
			}

		case ':':
			l.advance()
			if l.peek(0) == '=' {
				l.advance()
				l.addToken(token.BIND_ASSIGN, ":=", line, column)
			} else if l.peek(0) == ':' {
				l.advance()
				l.addToken(token.DOUBLE_COLON, "::", line, column)
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
			} else if l.peek(0) == '>' {
				l.advance()
				if l.peek(0) == '>' {
					l.advance()
					l.addToken(token.PIPE_PARALLEL, "|>>", line, column)
				} else {
					l.addToken(token.PIPE_RIGHT, "|>", line, column)
				}
			} else {
				l.addToken(token.PIPE_OP, "|", line, column)
			}

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
