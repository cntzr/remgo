package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/cntzr/remgo/ast"
	"github.com/cntzr/remgo/lexer"
)

type remStructure struct {
	day          string
	month        string
	year         string
	at           string
	atTime       string
	duration     string
	durationTime string
	msg          string
	msgTxt       string
}

func TestREMStatements(t *testing.T) {
	f, err := os.Open("../test_parser.rem")
	if err != nil {
		fmt.Printf("Error while open test file: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	l := lexer.New(f)
	p := New(l)
	reminders := p.ParseReminders()

	if reminders == nil {
		t.Fatalf("ParseReminders() returned nil")
	}
	if len(reminders.Statements) != 3 {
		t.Fatalf("reminders.Statements does not contain 3 reminders. got=%d", len(reminders.Statements))
	}

	tests := []remStructure{
		{"02", "Jun", "2021", "at", "12:30", "", "", "msg", "Erster Test"},
		{"16", "", "2021", "at", "09:15", "duration", "1:00", "msg", "Zweiter Test"},
		{"27", "May", "", "at", "23:00", "duration", "0:30", "msg", "Dritter Test"},
	}

	for i, tt := range tests {
		stmt := reminders.Statements[i]
		if !testREMStatement(t, stmt, tt) {
			return
		}
	}
}

func testREMStatement(t *testing.T, s ast.Statement, data remStructure) bool {
	if s.TokenLiteral() != "REM" {
		t.Errorf("s.TokenLiteral not 'REM', got=%q", s.TokenLiteral())
		return false
	}
	remStmt, ok := s.(*ast.REMStatement)
	if !ok {
		t.Errorf("s not *ast.REMStatement, got=%T", s)
		return false
	}

	if remStmt.Day != nil && remStmt.Day.Value != data.day {
		t.Errorf("remStmt.Day.Value not '%s'. got=%s", data.day, remStmt.Day.Value)
		return false
	}

	if remStmt.Day != nil && remStmt.Day.TokenLiteral() != data.day {
		t.Errorf("remStmt.Day.TokenLiteral() not '%s'. got=%s",
			data.day, remStmt.Day.TokenLiteral())
		return false
	}
	if remStmt.Month != nil && remStmt.Month.Value != data.month {
		t.Errorf("remStmt.Month.Value not '%s', got=%s", data.month, remStmt.Month.Value)
		return false
	}

	if remStmt.Month != nil && remStmt.Month.TokenLiteral() != data.month {
		t.Errorf("remStmt.Month.TokenLiteral() not '%s', got=%s", data.month, remStmt.Month.TokenLiteral())
		return false
	}

	if remStmt.Year != nil && remStmt.Year.Value != data.year {
		t.Errorf("remStmt.Year.Value not '%s', got=%s", data.year, remStmt.Year.Value)
		return false
	}

	if remStmt.Year != nil && remStmt.Year.TokenLiteral() != data.year {
		t.Errorf("remStmt.Year.TokenLiteral() not '%s', got=%s", data.year, remStmt.Year.TokenLiteral())
		return false
	}

	if remStmt.AtTime != nil && remStmt.AtTime.Value != data.atTime {
		t.Errorf("remStmt.AtTime.Value not '%s', got=%s", data.atTime, remStmt.AtTime.Value)
		return false
	}

	if remStmt.AtTime != nil && remStmt.AtTime.TokenLiteral() != data.atTime {
		t.Errorf("remStmt.AtTime.TokenLiteral() not '%s', got=%s", data.atTime, remStmt.AtTime.TokenLiteral())
		return false
	}

	if remStmt.DurationTime != nil && remStmt.DurationTime.Value != data.durationTime {
		t.Errorf("remStmt.DurationTime.Value not '%s', got=%s", data.durationTime, remStmt.DurationTime.Value)
		return false
	}

	if remStmt.DurationTime != nil && remStmt.DurationTime.TokenLiteral() != data.durationTime {
		t.Errorf("remStmt.DurationTime.TokenLiteral() not '%s', got=%s", data.durationTime, remStmt.DurationTime.TokenLiteral())
		return false
	}

	if remStmt.MsgTxt != nil && remStmt.MsgTxt.Value != data.msgTxt {
		t.Errorf("remStmt.MsgTxt.Value not '%s', got=%s", data.msgTxt, remStmt.MsgTxt.Value)
		return false
	}

	if remStmt.MsgTxt != nil && remStmt.MsgTxt.TokenLiteral() != data.msgTxt {
		t.Errorf("remStmt.MsgTxt.TokenLiteral() not '%s', got=%s", data.msgTxt, remStmt.MsgTxt.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
}
