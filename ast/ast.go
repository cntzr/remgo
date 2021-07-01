package ast

import (
	"github.com/cntzr/remgo/token"
)

type (
	Node interface {
		TokenLiteral() string
	}

	// Represents one complete reminder as a one-liner
	Statement interface {
		Node
		statementNode()
	}

	// List of all reminders in one file or over all files
	Reminders struct {
		Statements []Statement
	}

	// Marks a reminder ... in the future there may be RUN statements as well
	REMStatement struct {
		Token        token.Token // the REM token
		Weekdays     []*Weekday
		Day          *Day
		Month        *Month
		Year         *Year
		Lead         *Lead
		At           *At
		AtTime       *AtTime
		Duration     *Duration
		DurationTime *DurationTime
		Repeat       *Repeat
		Msg          *Msg
		MsgTxt       *MsgTxt
		Value        string
	}

	Month struct {
		Token token.Token // the Element's token
		Value string
	}

	Weekday struct {
		Token token.Token // the Element's token
		Value string
	}

	Day struct {
		Token token.Token // the Element's token
		Value string
	}

	Year struct {
		Token token.Token // the Element's token
		Value string
	}

	Repeat struct {
		Token token.Token // the Element's token
		Value string
	}

	At struct {
		Token token.Token // the Element's token
		Value string
	}

	AtTime struct {
		Token token.Token // the Element's token
		Value string
	}

	Duration struct {
		Token token.Token // the Element's token
		Value string
	}

	DurationTime struct {
		Token token.Token // the Element's token
		Value string
	}

	Lead struct {
		Token token.Token // the Element's token
		Value string
	}

	Msg struct {
		Token token.Token // the Element's token
		Value string
	}

	MsgTxt struct {
		Token token.Token // the Element's token
		Value string
	}
)

func (r *Reminders) TokenLiteral() string {
	if len(r.Statements) > 0 {
		return r.Statements[0].TokenLiteral()
	}
	return ""
}

func (rs *REMStatement) statementNode() {}

func (rs *REMStatement) TokenLiteral() string { return rs.Token.Literal }

func (e *Weekday) statementNode() {}

func (e *Weekday) TokenLiteral() string { return e.Token.Literal }

func (e *Month) statementNode() {}

func (e *Month) TokenLiteral() string { return e.Token.Literal }

func (e *Day) statementNode() {}

func (e *Day) TokenLiteral() string { return e.Token.Literal }

func (e *Year) statementNode() {}

func (e *Year) TokenLiteral() string { return e.Token.Literal }

func (rs *Repeat) statementNode() {}

func (e *Repeat) TokenLiteral() string { return e.Token.Literal }

func (rs *At) statementNode() {}

func (e *At) TokenLiteral() string { return e.Token.Literal }

func (rs *AtTime) statementNode() {}

func (e *AtTime) TokenLiteral() string { return e.Token.Literal }

func (rs *Duration) statementNode() {}

func (e *Duration) TokenLiteral() string { return e.Token.Literal }

func (rs *DurationTime) statementNode() {}

func (e *DurationTime) TokenLiteral() string { return e.Token.Literal }

func (rs *Lead) statementNode() {}

func (e *Lead) TokenLiteral() string { return e.Token.Literal }

func (rs *Msg) statementNode() {}

func (e *Msg) TokenLiteral() string { return e.Token.Literal }

func (rs *MsgTxt) statementNode() {}

func (e *MsgTxt) TokenLiteral() string { return e.Token.Literal }
