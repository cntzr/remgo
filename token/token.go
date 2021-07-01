package token

import (
	"strconv"
	"strings"
)

type (
	TokenType string

	Token struct {
		Type    TokenType
		Literal string
	}
)

var (
	keywords = map[string]TokenType{
		"REM": REM,
		"RUN": RUN,
	}

	attributes = map[string]TokenType{
		"AT":       AT,
		"DURATION": DURATION,
		"MSG":      MSG,
	}

	months = map[string]TokenType{
		"Jan": MONTH,
		"Feb": MONTH,
		"Mar": MONTH,
		"Apr": MONTH,
		"May": MONTH,
		"Jun": MONTH,
		"Jul": MONTH,
		"Aug": MONTH,
		"Sep": MONTH,
		"Oct": MONTH,
		"Nov": MONTH,
		"Dec": MONTH,
	}

	weekdays = map[string]TokenType{
		"Mon": WEEKDAY,
		"Tue": WEEKDAY,
		"Wed": WEEKDAY,
		"Thu": WEEKDAY,
		"Fri": WEEKDAY,
		"Sat": WEEKDAY,
		"Sun": WEEKDAY,
	}
)

const (
	// partials
	TIME    = "TIME"    // e.g. 14:30
	DAY     = "DAY"     // e.g. 12
	WEEKDAY = "WEEKDAY" // e.g. Mon, Tue or Sat
	MONTH   = "MONTH"   // e.g. Jun
	YEAR    = "YEAR"    // e.g. 2021
	REPEAT  = "REPEAT"  // e.g. *14 for all 2 weeks
	LEAD    = "LEAD"    // e.g. +3 for 3 days before
	MSG_TXT = "MSG"     // e.g. Obst kaufen

	// keywords
	REM = "REM"
	RUN = "RUN"

	// elements
	AT       = "AT"
	DURATION = "DURATION"
	MSG      = "MSG"

	// delimiter
	NEWLINE = "NEWLINE" // end of line

	// all other stuff
	ILLEGAL    = "ILLEGAL"    // illegal token
	WHITESPACE = "WHITESPACE" // inline blank or tab
	EOF        = "EOF"        // end of file
)

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	if tok, ok := attributes[ident]; ok {
		return tok
	}
	if tok, ok := months[ident]; ok {
		return tok
	}
	if tok, ok := weekdays[ident]; ok {
		return tok
	}
	return ILLEGAL
}

func LookupNumerics(ident string) TokenType {
	// first handle times and periods
	s := strings.Split(ident, ":")
	if len(s) == 2 {
		return TIME
	}
	// now check for real numbers
	i, _ := strconv.Atoi(ident)
	if 1 <= i && i <= 31 {
		return DAY
	}
	if 1971 <= i && i <= 2055 {
		return YEAR
	}
	return ILLEGAL
}

func LookupSign(ident string) TokenType {
	if ident[0] == '*' {
		return REPEAT
	}
	if ident[0] == '+' {
		return LEAD
	}
	return ILLEGAL
}
