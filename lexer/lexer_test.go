package lexer

import (
	"fmt"
	"os"
	"testing"

	"github.com/cntzr/remgo/token"
)

func TestNextToken(t *testing.T) {
	f, err := os.Open("../test_lexer.rem")
	if err != nil {
		fmt.Printf("Error while open test file: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.REM, "REM"},
		{token.MONTH, "Jun"},
		{token.DAY, "5"},
		{token.YEAR, "2021"},
		{token.LEAD, "+1"},
		{token.AT, "AT"},
		{token.TIME, "08:30"},
		{token.DURATION, "DURATION"},
		{token.TIME, "1:30"},
		{token.REPEAT, "*14"},
		{token.MSG, "MSG"},
		{token.MSG_TXT, "wir testen"},
	}

	l := New(f)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%s, got=%s", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%s, got=%s", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
