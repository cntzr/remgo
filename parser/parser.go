package parser

import (
	"fmt"

	"github.com/cntzr/remind-go/ast"
	"github.com/cntzr/remind-go/lexer"
	"github.com/cntzr/remind-go/token"
)

type (
	Parser struct {
		l         *lexer.Lexer
		errors    []string
		curToken  token.Token
		peekToken token.Token
	}
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// loads curToken and nextToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseReminders() *ast.Reminders {
	reminders := &ast.Reminders{}
	reminders.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			reminders.Statements = append(reminders.Statements, stmt)
		}
		p.nextToken()
	}

	return reminders
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.REM:
		return p.parseREMStatement()
	default:
		return nil
	}
}

func (p *Parser) parseREMStatement() *ast.REMStatement {

	stmt := &ast.REMStatement{
		Token:    p.curToken,
		Weekdays: make([]*ast.Weekday, 0),
	}

	// day, month & year of the event
	switch true {
	case p.expectPeek(token.MONTH):
		stmt.Month = &ast.Month{Token: p.curToken, Value: p.curToken.Literal}
		if p.expectPeek(token.DAY) {
			stmt.Day = &ast.Day{Token: p.curToken, Value: p.curToken.Literal}
		}
		if p.expectPeek(token.YEAR) {
			stmt.Year = &ast.Year{Token: p.curToken, Value: p.curToken.Literal}
		}
	case p.expectPeek(token.DAY):
		stmt.Day = &ast.Day{Token: p.curToken, Value: p.curToken.Literal}
		if p.expectPeek(token.YEAR) {
			stmt.Year = &ast.Year{Token: p.curToken, Value: p.curToken.Literal}
		}
	case p.expectPeek(token.YEAR):
		stmt.Year = &ast.Year{Token: p.curToken, Value: p.curToken.Literal}
	case p.expectPeek(token.WEEKDAY):
		stmt.Weekdays = append(stmt.Weekdays, &ast.Weekday{Token: p.curToken, Value: p.curToken.Literal})
		for p.expectPeek(token.WEEKDAY) {
			stmt.Weekdays = append(stmt.Weekdays, &ast.Weekday{Token: p.curToken, Value: p.curToken.Literal})
		}
	case p.expectPeek(token.LEAD):
		stmt.Lead = &ast.Lead{Token: p.curToken, Value: p.curToken.Literal}
	}

	// At which time (optional)
	switch true {
	case p.expectPeek(token.AT):
		stmt.At = &ast.At{Token: p.curToken, Value: p.curToken.Literal}
		if p.expectPeek(token.TIME) {
			stmt.AtTime = &ast.AtTime{Token: p.curToken, Value: p.curToken.Literal}
		}
	}

	// For how long (optional)
	if p.expectPeek(token.DURATION) {
		stmt.Duration = &ast.Duration{Token: p.curToken, Value: p.curToken.Literal}
		if p.expectPeek(token.TIME) {
			stmt.DurationTime = &ast.DurationTime{Token: p.curToken, Value: p.curToken.Literal}
		}
	}

	// Repeat every n days (optional)
	if p.expectPeek(token.REPEAT) {
		stmt.Repeat = &ast.Repeat{Token: p.curToken, Value: p.curToken.Literal}
	}

	// Description of the event
	if p.expectPeek(token.MSG) {
		stmt.Msg = &ast.Msg{Token: p.curToken, Value: p.curToken.Literal}
		if p.expectPeek(token.MSG_TXT) {
			stmt.MsgTxt = &ast.MsgTxt{Token: p.curToken, Value: p.curToken.Literal}
		}
	}

	// TODO ... Prüfen, ob da noch was übrig ist
	if p.expectPeek(token.NEWLINE) {
		return stmt
	}

	return nil
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
