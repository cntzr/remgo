package timeframe

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cntzr/remind-go/object"
)

type (
	TimeFrame struct {
		from  time.Time
		until time.Time
		Days  map[string]*Day
	}

	Day struct {
		Date    time.Time
		Weekday string
		Events  []*Event
	}

	Event struct {
		From        string
		Until       string
		Description string
	}
)

func NewTimeFrame(from, until time.Time) TimeFrame {
	d := make(map[string]*Day)
	e := make([]*Event, 0)
	w := from.Weekday().String()[0:3]
	f := Day{Date: from, Weekday: w, Events: e}
	d[from.Format("2006-01-02")] = &f

	nextDay := from.AddDate(0, 0, 1)
	for nextDay.Before(until.AddDate(0, 0, 1)) {
		e := make([]*Event, 0)
		w := nextDay.Weekday().String()[0:3]
		n := Day{Date: nextDay, Weekday: w, Events: e}
		d[nextDay.Format("2006-01-02")] = &n
		nextDay = nextDay.AddDate(0, 0, 1)
	}

	return TimeFrame{from: from, until: until, Days: d}
}

func (t *TimeFrame) FillFrame(env *object.Environment) {
	weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	// process weekdays first
	wasWeekday := false
	for _, day := range weekdays {
		if _, ok := env.Get(day); ok {
			for date, x := range t.Days {
				if x.Weekday == day {
					from := ""
					until := ""
					desc := ""
					if f, ok := env.Get("From"); ok {
						from = f.Inspect()
					}
					if u, ok := env.Get("Until"); ok {
						until = u.Inspect()
					}
					if d, ok := env.Get("MsgTxt"); ok {
						desc = d.Inspect()
					}
					t.addEvent(date, from, until, desc)
				}
			}
			wasWeekday = true
		}
	}
	// there can be only weekdays or real dates ...
	// ... and there can be multiple weekdays so we only return here
	if wasWeekday {
		return
	}
	// process events with at least partial dates
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

	// check for repeated events
	if r, ok := env.Get("Repeat"); ok {
		period, err := strconv.Atoi(r.Inspect())
		if err == nil {
			date := Date(year, int(month), day)
			for date.Before(t.from) { // add repeat period until we are behind start of time frame
				date = date.AddDate(0, 0, period)
			}
			if !date.After(t.until) { // we are in the time frame so update date of event
				year = date.Year()
				month = date.Month()
				day = date.Day()
			}
		}
	}

	if betweenDays(Date(year, int(month), day), t.from, t.until) {
		date := Date(year, int(month), day).Format("2006-01-02")
		from := ""
		until := ""
		desc := ""
		if f, ok := env.Get("From"); ok {
			from = f.Inspect()
		}
		if u, ok := env.Get("Until"); ok {
			until = u.Inspect()
		}
		if d, ok := env.Get("MsgTxt"); ok {
			desc = d.Inspect()
		}

		t.addEvent(date, from, until, desc)
	}
}

func (t *TimeFrame) addEvent(day, from, until, desc string) error {
	if d, ok := t.Days[day]; ok {
		d.Events = append(d.Events, &Event{From: from, Until: until, Description: desc})
		return nil
	} else {
		return fmt.Errorf("%s is not part of the current time frame", day)
	}
}

func Date(year, month, day int) time.Time {
	loc, _ := time.LoadLocation("Europe/Berlin")
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

func SameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func betweenDays(date, bottom, top time.Time) bool {
	return (bottom.Before(date) || SameDay(date, bottom)) && (top.After(date) || SameDay(date, top))
}

func FirstDayOfWeek() time.Time {
	date := Date(time.Now().Year(), int(time.Now().Month()), time.Now().Day())
	if time.Now().Weekday().String() == "Monday" { // today is Monday
		return date
	}
	for date.Weekday().String() != "Monday" {
		date = date.AddDate(0, 0, -1)
	}
	return date
}

func LastDayOfWeek() time.Time {
	date := Date(time.Now().Year(), int(time.Now().Month()), time.Now().Day())
	if time.Now().Weekday().String() == "Sunday" { // today is Sunday
		return date
	}
	for date.Weekday().String() != "Sunday" {
		date = date.AddDate(0, 0, 1)
	}
	return date
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
