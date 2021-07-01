package object

import (
	"fmt"
)

type (
	Environment struct {
		store map[string]Object
	}

	ObjectType string

	Object interface {
		Type() ObjectType
		Inspect() string
	}

	// REM collects the whole REMStatement
	REM struct {
		Value string
	}

	Weekday struct {
		Value string
	}

	Day struct {
		Value string
	}

	Month struct {
		Value string
	}

	Year struct {
		Value string
	}

	From struct {
		Value string
	}

	Until struct {
		Value string
	}

	Lead struct {
		Value string
	}

	Repeat struct {
		Value string
	}

	MsgTxt struct {
		Value string
	}

	Error struct {
		Message string
	}
)

const (
	REM_OBJ     = "REM" // means the whole REMStatement
	WEEKDAY_OBJ = "WEEKDAY"
	DAY_OBJ     = "DAY"
	MONTH_OBJ   = "MONTH"
	YEAR_OBJ    = "YEAR"
	FROM_OBJ    = "FROM"
	UNTIL_OBJ   = "UNTIL"
	LEAD_OBJ    = "LEAD"
	REPEAT_OBJ  = "REPEAT"
	MSG_TXT_OBJ = "MSG_TXT"
	ERROR_OBJ   = "ERROR"
)

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (d *REM) Inspect() string { return d.Value }

func (d *REM) Type() ObjectType { return DAY_OBJ }

func (d *Weekday) Inspect() string { return d.Value }

func (d *Weekday) Type() ObjectType { return DAY_OBJ }

func (d *Day) Inspect() string { return d.Value }

func (d *Day) Type() ObjectType { return DAY_OBJ }

func (m *Month) Inspect() string { return m.Value }

func (m *Month) Type() ObjectType { return MONTH_OBJ }

func (y *Year) Inspect() string { return y.Value }

func (y *Year) Type() ObjectType { return YEAR_OBJ }

func (t *From) Inspect() string { return t.Value }

func (t *From) Type() ObjectType { return FROM_OBJ }

func (t *Until) Inspect() string { return t.Value }

func (t *Until) Type() ObjectType { return UNTIL_OBJ }

func (l *Lead) Inspect() string { return l.Value }

func (l *Lead) Type() ObjectType { return LEAD_OBJ }

func (r *Repeat) Inspect() string { return r.Value }

func (r *Repeat) Type() ObjectType { return REPEAT_OBJ }

func (t *MsgTxt) Inspect() string { return t.Value }

func (t *MsgTxt) Type() ObjectType { return MSG_TXT_OBJ }

func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s", e.Message)
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
