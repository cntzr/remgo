package evaluator

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cntzr/remgo/ast"
	"github.com/cntzr/remgo/object"
	"github.com/cntzr/remgo/timeframe"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.REMStatement:
		if node.Day != nil {
			val := Eval(node.Day, env)
			if isError(val) {
				// TODO muss noch ordentlich verarbeitet werden
				return val
			}
			env.Set("Day", val)
		}
		if node.Month != nil {
			val := Eval(node.Month, env)
			if isError(val) {
				// TODO muss noch ordentlich verarbeitet werden
				return val
			}
			env.Set("Month", val)
		}
		if node.Year != nil {
			val := Eval(node.Year, env)
			if isError(val) {
				// TODO muss noch ordentlich verarbeitet werden
				return val
			}
			env.Set("Year", val)
		}
		if len(node.Weekdays) > 0 {
			for _, d := range node.Weekdays {
				val := Eval(d, env)
				if isError(val) {
					// TODO muss noch ordentlich verarbeitet werden
					return val
				}
				env.Set(val.Inspect(), val)
			}
		}
		if node.At != nil && node.AtTime != nil {
			val := Eval(node.AtTime, env)
			if isError(val) {
				// TODO muss noch ordentlich verarbeitet werden
				return val
			}
			env.Set("From", val)
		}
		if node.Duration != nil && node.DurationTime != nil {
			val := Eval(node.DurationTime, env)
			if isError(val) {
				// TODO muss noch ordentlich verarbeitet werden
				return val
			}
			// calculate the end of date
			if f, ok := env.Get("From"); ok {
				year := time.Now().Year()
				if y, ok := env.Get("Year"); ok {
					year, _ = strconv.Atoi(y.Inspect())
				}
				month := time.Now().Month()
				if m, ok := env.Get("Month"); ok {
					month = mToMonth(m.Inspect())
				}
				day := time.Now().Day()
				if d, ok := env.Get("Day"); ok {
					day, _ = strconv.Atoi(d.Inspect())
				}

				sFrom := strings.Split(f.Inspect(), ":")
				hourFrom, _ := strconv.Atoi(sFrom[0])
				minFrom, _ := strconv.Atoi(sFrom[1])

				sUntil := strings.Split(val.Inspect(), ":")
				hourUntil, _ := strconv.Atoi(sUntil[0])
				minUntil, _ := strconv.Atoi(sUntil[1])

				loc, _ := time.LoadLocation("Europe/Berlin")
				until := time.Date(year, month, day, hourFrom, minFrom, 0, 0, loc)
				until = until.Add(time.Hour * time.Duration(hourUntil))
				until = until.Add(time.Minute * time.Duration(minUntil))
				env.Set("Until", &object.Until{Value: until.Format("15:04")})
			}
		}

		if node.Repeat != nil {
			val := Eval(node.Repeat, env)
			if isError(val) {
				// TODO muss noch ordentlich verarbeitet werden
				return val
			}
			env.Set("Repeat", val)
		}

		// TODO ... handle node.Msg == nil, means MSG is missing
		if node.Msg != nil && node.MsgTxt != nil {
			val := Eval(node.MsgTxt, env)
			if isError(val) {
				// TODO muss noch ordentlich verarbeitet werden
				return val
			}
			env.Set("MsgTxt", val)
		}

		msg := buildMsg(env)
		if msg == "" {
			return nil
		}
		return &object.REM{Value: msg}

		// TODO alle Attribute des REMStatements evaluieren, es fehlen noch ...
		//      [ ] Repeat
		//      [ ] Lead

	case *ast.Weekday:
		return &object.Weekday{Value: node.Value}
	case *ast.Day:
		d, err := strconv.Atoi(node.Value)
		if err != nil {
			return NewError("Conversion from Day to int failed => %s", err.Error())
		}
		if 1 <= d && d <= 31 {
			return &object.Day{Value: node.Value}
		}
		return &object.Error{Message: "Invalid Day, not between 1 and 31"}
	case *ast.Month:
		return &object.Month{Value: node.Value}
	case *ast.Year:
		y, err := strconv.Atoi(node.Value)
		if err != nil {
			return NewError("Conversion from Year to int failed => %s", err.Error())
		}
		if 2000 < y && y < 2055 {
			return &object.Year{Value: node.Value}
		}
		return &object.Error{Message: "Invalid Year, not between 2000 and 2055"}
	case *ast.AtTime:
		if isTime(node.Value) {
			return &object.From{Value: node.Value}
		}
		return &object.Error{Message: "Invalid Time"}
	case *ast.DurationTime:
		return &object.Until{Value: node.Value}
	case *ast.Repeat:
		return &object.Repeat{Value: node.Value[1:]}
	case *ast.MsgTxt:
		return &object.MsgTxt{Value: node.Value}
	}
	return &object.Error{Message: "Can't evaluate node"}
}

func NewError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isTime(t string) bool {
	s := strings.Split(t, ":")
	h, _ := strconv.Atoi(s[0])
	m, _ := strconv.Atoi(s[1])

	if 0 <= h && h <= 23 && 0 <= m && m <= 59 {
		return true
	}
	return false
}

// TODO ... method is deprecated cause it's reimplemented by timeframe.FillFrame()
func buildMsg(env *object.Environment) string {
	d, _ := env.Get("Day")
	m, _ := env.Get("Month")
	y, _ := env.Get("Year")

	day := time.Now().Day()
	month := time.Now().Month()
	year := time.Now().Year()

	if d != nil {
		day, _ = strconv.Atoi(d.Inspect())
	}
	if m != nil {
		month = mToMonth(m.Inspect())
	}
	if y != nil {
		year, _ = strconv.Atoi(y.Inspect())
	}

	tRem := timeframe.Date(year, int(month), day)
	today := timeframe.Date(time.Now().Year(), int(time.Now().Month()), time.Now().Day())
	deltaDays := today.Sub(tRem).Hours() / 24
	// deltaDays > 0 ... event in the past
	// deltaDays < -1 ... event in the future
	if -1 > deltaDays || deltaDays > 0 { // not today or tomorrow
		return ""
	}

	var t string
	f, _ := env.Get("From")
	u, _ := env.Get("Until")

	if f != nil && u != nil {
		sFrom := strings.Split(f.Inspect(), ":")
		hourFrom, _ := strconv.Atoi(sFrom[0])
		minFrom, _ := strconv.Atoi(sFrom[1])

		sUntil := strings.Split(u.Inspect(), ":")
		hourUntil, _ := strconv.Atoi(sUntil[0])
		minUntil, _ := strconv.Atoi(sUntil[1])

		loc, _ := time.LoadLocation("Europe/Berlin")
		until := time.Date(year, month, day, hourFrom, minFrom, 0, 0, loc)
		until = until.Add(time.Hour * time.Duration(hourUntil))
		until = until.Add(time.Minute * time.Duration(minUntil))
		t = fmt.Sprintf(", %s - %s", f.Inspect(), until.Format("15:04"))
	}
	if f != nil && u == nil {
		t = fmt.Sprintf(", %s", f.Inspect())
	}

	msg, _ := env.Get("MsgTxt")

	return fmt.Sprintf("%d-%02d-%02d%s, %s", year, month, day, t, msg.Inspect())
}

func mToMonth(m string) time.Month {
	var i int

	switch m {
	case "Jan":
		i = 1
	case "Feb":
		i = 2
	case "Mar":
		i = 3
	case "Apr":
		i = 4
	case "May":
		i = 5
	case "Jun":
		i = 6
	case "Jul":
		i = 7
	case "Aug":
		i = 8
	case "Sep":
		i = 9
	case "Oct":
		i = 10
	case "Nov":
		i = 11
	case "Dec":
		i = 12
	default:
		i = 0
	}

	return time.Month(i)
}
