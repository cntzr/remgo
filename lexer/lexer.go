package lexer

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"unicode"

	"github.com/cntzr/remind-go/token"
)

type (
	Lexer struct {
		reader    *bufio.Reader
		ch        rune // current char under examination
		lastToken token.Token
	}
)

// New ... creates a new Lexer
func New(f *os.File) *Lexer {

	l := &Lexer{reader: bufio.NewReader(f)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	ch, _, err := l.reader.ReadRune()
	if err == io.EOF {
		l.ch = rune(0)
	} else {
		l.ch = ch
	}
}

func (l *Lexer) unreadChar() {
	l.reader.UnreadRune()
}

// NextToken ... extracts the next Token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhitespace()

	switch true {
	case l.isEOF():
		tok.Literal = ""
		tok.Type = token.EOF
	case unicode.IsLetter(l.ch):
		l.unreadChar() // roll back so we can collect the letter
		if l.lastToken.Type == token.MSG {
			tok.Literal = l.readTheRest()
			tok.Type = token.MSG_TXT
		} else {
			tok.Literal = l.readWord()
			tok.Type = token.LookupIdent(tok.Literal)
		}
	case l.isDigit():
		l.unreadChar() // roll back so we can collect the digit
		if l.lastToken.Type == token.MSG {
			tok.Literal = l.readTheRest()
			tok.Type = token.MSG_TXT
		} else {
			tok.Literal = l.readNumber()
			tok.Type = token.LookupNumerics(tok.Literal)
		}
	case l.isSign():
		l.unreadChar() // roll back so we can collect the sign
		tok.Literal = l.readSign()
		tok.Type = token.LookupSign(tok.Literal)
	case l.isNewline():
		l.readChar()
		tok.Type = token.NEWLINE
		tok.Literal = ""
	default:
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.lastToken = tok
	return tok
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
func (l *Lexer) readWhitespace() string {
	var buf bytes.Buffer

	l.readChar()
	buf.WriteRune(l.ch)
	for {
		l.readChar()
		if l.ch == rune(0) {
			break
		}
		if !l.isWhitespace() {
			l.unreadChar()
			break
		}
		buf.WriteRune(l.ch)
	}

	return buf.String()
}

func (l *Lexer) readWord() string {
	var buf bytes.Buffer

	l.readChar()
	buf.WriteRune(l.ch)
	for {
		l.readChar()
		if l.ch == rune(0) { // EOF
			break
		}
		if !unicode.IsLetter(l.ch) {
			l.unreadChar()
			break
		}
		buf.WriteRune(l.ch)
	}

	return buf.String()
}

func (l *Lexer) readTheRest() string {
	var buf bytes.Buffer

	l.readChar()
	buf.WriteRune(l.ch)
	for {
		l.readChar()
		if l.ch == rune(0) { // EOF
			break
		}
		if l.ch == '\n' { // EOL
			break
		}
		buf.WriteRune(l.ch)
	}

	return buf.String()
}

func (l *Lexer) readNumber() string {
	var buf bytes.Buffer

	l.readChar()
	buf.WriteRune(l.ch)
	for {
		l.readChar()
		if l.ch == rune(0) { // EOF
			break
		}
		if !l.isDigit() {
			l.unreadChar()
			break
		}
		buf.WriteRune(l.ch)
	}

	return buf.String()
}

func (l *Lexer) readSign() string {
	var buf bytes.Buffer

	l.readChar()
	buf.WriteRune(l.ch)
	for {
		l.readChar()
		if l.ch == rune(0) { // EOF
			break
		}
		if !l.isDigit() { // signs are followed only by digits
			l.unreadChar()
			break
		}
		buf.WriteRune(l.ch)
	}

	return buf.String()
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}
}

/*********************************************
helpers for first identification
*********************************************/
func (l *Lexer) isWhitespace() bool {
	return l.ch == ' ' || l.ch == '\t'
}

func (l *Lexer) isDigit() bool {
	return '0' <= l.ch && l.ch <= '9' || l.ch == ':'
}

func (l *Lexer) isSign() bool {
	return l.ch == '+' || l.ch == '*'
}

func (l *Lexer) isNewline() bool {
	return l.ch == '\n'
}

func (l *Lexer) isEOF() bool {
	return l.ch == rune(0)
}
