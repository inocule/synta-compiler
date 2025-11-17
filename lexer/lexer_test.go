// lexer/lexer_test.go
package lexer

import (
	"synta/token"
	"testing"
)

func TestBasicTokens(t *testing.T) {
	input := `bind x := 10;`

	expected := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.BIND, "bind"},
		{token.IDENTIFIER, "x"},
		{token.BIND_ASSIGN, ":="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)
	tokens := l.Tokenize()

	if len(tokens) != len(expected) {
		t.Fatalf("wrong number of tokens. expected=%d, got=%d", len(expected), len(tokens))
	}

	for i, tt := range expected {
		if tokens[i].Type != tt.expectedType {
			t.Errorf("tokens[%d] - wrong type. expected=%q, got=%q",
				i, tt.expectedType, tokens[i].Type)
		}

		if tokens[i].Lexeme != tt.expectedLexeme {
			t.Errorf("tokens[%d] - wrong lexeme. expected=%q, got=%q",
				i, tt.expectedLexeme, tokens[i].Lexeme)
		}
	}
}

func TestFunctionDeclaration(t *testing.T) {
	input := `fn calculate(a, b) {
		return a + b;
	}`

	expected := []token.TokenType{
		token.FN,
		token.IDENTIFIER,
		token.LPAREN,
		token.IDENTIFIER,
		token.COMMA,
		token.IDENTIFIER,
		token.RPAREN,
		token.LBRACE,
		token.NEWLINE,
		token.RETURN,
		token.IDENTIFIER,
		token.PLUS,
		token.IDENTIFIER,
		token.SEMICOLON,
		token.NEWLINE,
		token.RBRACE,
		token.EOF,
	}

	l := New(input)
	tokens := l.Tokenize()

	if len(tokens) != len(expected) {
		t.Fatalf("wrong number of tokens. expected=%d, got=%d", len(expected), len(tokens))
	}

	for i, expectedType := range expected {
		if tokens[i].Type != expectedType {
			t.Errorf("tokens[%d] - wrong type. expected=%q, got=%q",
				i, expectedType, tokens[i].Type)
		}
	}
}

func TestOperators(t *testing.T) {
	input := `+ - * / % == != < > <= >= && || ! ++ -- += -= *= /= %=`

	expected := []token.TokenType{
		token.PLUS,
		token.MINUS,
		token.MULTIPLY,
		token.DIVIDE,
		token.MODULO,
		token.EQ,
		token.NEQ,
		token.LT,
		token.GT,
		token.LTE,
		token.GTE,
		token.AND,
		token.OR,
		token.NOT,
		token.INCREMENT,
		token.DECREMENT,
		token.PLUS_ASSIGN,
		token.MINUS_ASSIGN,
		token.MULT_ASSIGN,
		token.DIV_ASSIGN,
		token.MOD_ASSIGN,
		token.EOF,
	}

	l := New(input)
	tokens := l.Tokenize()

	if len(tokens) != len(expected) {
		t.Fatalf("wrong number of tokens. expected=%d, got=%d", len(expected), len(tokens))
	}

	for i, expectedType := range expected {
		if tokens[i].Type != expectedType {
			t.Errorf("tokens[%d] - wrong type. expected=%q, got=%q",
				i, expectedType, tokens[i].Type)
		}
	}
}

func TestComments(t *testing.T) {
	input := `bind x := 10; // This is a comment
	bind y := 20;
	/~ This is a
	   multi-line comment ~/
	bind z := 30;`

	expected := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.BIND, "bind"},
		{token.IDENTIFIER, "x"},
		{token.BIND_ASSIGN, ":="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\\n"},
		{token.BIND, "bind"},
		{token.IDENTIFIER, "y"},
		{token.BIND_ASSIGN, ":="},
		{token.INTEGER, "20"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\\n"},
		{token.NEWLINE, "\\n"},
		{token.BIND, "bind"},
		{token.IDENTIFIER, "z"},
		{token.BIND_ASSIGN, ":="},
		{token.INTEGER, "30"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)
	tokens := l.Tokenize()

	if len(tokens) != len(expected) {
		t.Fatalf("wrong number of tokens. expected=%d, got=%d", len(expected), len(tokens))
	}

	for i, tt := range expected {
		if tokens[i].Type != tt.expectedType {
			t.Errorf("tokens[%d] - wrong type. expected=%q, got=%q",
				i, tt.expectedType, tokens[i].Type)
		}
	}
}

func TestNumbers(t *testing.T) {
	input := `42 3.14 0.5 999`

	expected := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.INTEGER, "42"},
		{token.FLOAT, "3.14"},
		{token.FLOAT, "0.5"},
		{token.INTEGER, "999"},
		{token.EOF, ""},
	}

	l := New(input)
	tokens := l.Tokenize()

	if len(tokens) != len(expected) {
		t.Fatalf("wrong number of tokens. expected=%d, got=%d", len(expected), len(tokens))
	}

	for i, tt := range expected {
		if tokens[i].Type != tt.expectedType {
			t.Errorf("tokens[%d] - wrong type. expected=%q, got=%q",
				i, tt.expectedType, tokens[i].Type)
		}

		if tokens[i].Lexeme != tt.expectedLexeme {
			t.Errorf("tokens[%d] - wrong lexeme. expected=%q, got=%q",
				i, tt.expectedLexeme, tokens[i].Lexeme)
		}
	}
}

func TestStrings(t *testing.T) {
	input := `"Hello, World!" 'Single quotes'`

	expected := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.STRING, "Hello, World!"},
		{token.STRING, "Single quotes"},
		{token.EOF, ""},
	}

	l := New(input)
	tokens := l.Tokenize()

	if len(tokens) != len(expected) {
		t.Fatalf("wrong number of tokens. expected=%d, got=%d", len(expected), len(tokens))
	}

	for i, tt := range expected {
		if tokens[i].Type != tt.expectedType {
			t.Errorf("tokens[%d] - wrong type. expected=%q, got=%q",
				i, tt.expectedType, tokens[i].Type)
		}

		if tokens[i].Lexeme != tt.expectedLexeme {
			t.Errorf("tokens[%d] - wrong lexeme. expected=%q, got=%q",
				i, tt.expectedLexeme, tokens[i].Lexeme)
		}
	}
}

func TestAIKeywords(t *testing.T) {
	input := `think ask prompt reason observe`

	expected := []token.TokenType{
		token.THINK,
		token.ASK,
		token.PROMPT,
		token.REASON,
		token.OBSERVE,
		token.EOF,
	}

	l := New(input)
	tokens := l.Tokenize()

	if len(tokens) != len(expected) {
		t.Fatalf("wrong number of tokens. expected=%d, got=%d", len(expected), len(tokens))
	}

	for i, expectedType := range expected {
		if tokens[i].Type != expectedType {
			t.Errorf("tokens[%d] - wrong type. expected=%q, got=%q",
				i, expectedType, tokens[i].Type)
		}
	}
}

func TestLineAndColumnTracking(t *testing.T) {
	input := `bind x := 10;
bind y := 20;`

	l := New(input)
	tokens := l.Tokenize()

	// Check first token line/column
	if tokens[0].Line != 1 || tokens[0].Column != 1 {
		t.Errorf("first token position wrong. expected=Line:1,Col:1, got=Line:%d,Col:%d",
			tokens[0].Line, tokens[0].Column)
	}

	// Check token on second line
	secondLineTokenIndex := 6 // "bind" on second line
	if tokens[secondLineTokenIndex].Line != 2 {
		t.Errorf("second line token wrong. expected=Line:2, got=Line:%d",
			tokens[secondLineTokenIndex].Line)
	}
}
