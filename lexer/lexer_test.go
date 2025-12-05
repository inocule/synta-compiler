package lexer

import (
	"synta-lexical/token"
	"testing"
)

func TestTildeDelimiter(t *testing.T) {
	tests := []struct {
		input    string
		expected token.TokenType
	}{
		{"~", token.STATEMENT_END},
		{"x ~", token.STATEMENT_END},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens := l.Tokenize()

		// Find the tilde token (skip NEWLINE and EOF)
		var tilde *token.Token
		for i := range tokens {
			if tokens[i].Type == token.STATEMENT_END {
				tilde = &tokens[i]
				break
			}
		}

		if tilde == nil {
			t.Fatalf("Expected STATEMENT_END token, got none from input: %q", tt.input)
		}
		if tilde.Type != tt.expected {
			t.Errorf("Expected %v, got %v", tt.expected, tilde.Type)
		}
	}
}

func TestMultilineComment(t *testing.T) {
	input := `<! type code here !>`
	l := New(input)
	tokens := l.Tokenize()

	// Should have exactly 2 tokens: COMMENT_MULTI and EOF
	if len(tokens) != 2 {
		t.Fatalf("Expected 2 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != token.COMMENT_MULTI {
		t.Errorf("Expected COMMENT_MULTI, got %v", tokens[0].Type)
	}

	if tokens[0].Lexeme != `<! type code here !>` {
		t.Errorf("Expected lexeme `<! type code here !>`, got %q", tokens[0].Lexeme)
	}

	if tokens[1].Type != token.EOF {
		t.Errorf("Expected EOF, got %v", tokens[1].Type)
	}
}

func TestMultilineCommentWithNewline(t *testing.T) {
	input := `<! this is
a multiline
comment !>`
	l := New(input)
	tokens := l.Tokenize()

	// Should have exactly 2 tokens: COMMENT_MULTI and EOF
	if len(tokens) != 2 {
		t.Fatalf("Expected 2 tokens, got %d: %v", len(tokens), tokens)
	}

	if tokens[0].Type != token.COMMENT_MULTI {
		t.Errorf("Expected COMMENT_MULTI, got %v", tokens[0].Type)
	}

	if tokens[1].Type != token.EOF {
		t.Errorf("Expected EOF, got %v", tokens[1].Type)
	}
}

func TestIncompleteMultilineComment(t *testing.T) {
	input := `<! incomplete comment`
	l := New(input)
	tokens := l.Tokenize()

	// Should have exactly 2 tokens: COMMENT_MULTI and EOF
	if len(tokens) != 2 {
		t.Fatalf("Expected 2 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != token.COMMENT_MULTI {
		t.Errorf("Expected COMMENT_MULTI, got %v", tokens[0].Type)
	}
}
