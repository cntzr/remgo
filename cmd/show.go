package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/cntzr/remgo/evaluator"
	"github.com/cntzr/remgo/lexer"
	"github.com/cntzr/remgo/object"
	"github.com/cntzr/remgo/parser"
	"github.com/cntzr/remgo/timeframe"
	"github.com/spf13/cobra"
)

var (
	showToday       = false
	showTomorrow    = false
	showCurrentWeek = false
	showNextWeek    = false

	eventColor  = lipgloss.AdaptiveColor{Light: "#C2C2A3", Dark: "#7A7A52"}
	allDayColor = lipgloss.AdaptiveColor{Light: "#C2D6D6", Dark: "#527A7A"}
	dayColor    = lipgloss.AdaptiveColor{Light: "#EBEBE0", Dark: "#B8B894"}
	markerColor = lipgloss.AdaptiveColor{Light: "#99CC00", Dark: "#739900"}
	warnColor   = lipgloss.AdaptiveColor{Light: "#CC0000", Dark: "#AC3939"}
	headerStyle = lipgloss.NewStyle().Foreground(markerColor).Italic(true)
	dayStyle    = lipgloss.NewStyle().Foreground(dayColor).Italic(true)
	todayStyle  = lipgloss.NewStyle().Foreground(warnColor).Italic(true)
	eventStyle  = lipgloss.NewStyle().Foreground(eventColor)
	allDayStyle = lipgloss.NewStyle().Foreground(allDayColor)

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(markerColor).
			PaddingRight(1).
			String()

	// showCmd represents the show command
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows dates, default is today",
		Long: `Shows dates, default is today. Support for a 
specific day or specific periods maybe coming later.
Output is colored by default.`,
		Run: func(cmd *cobra.Command, args []string) {
			today := time.Now()
			from := today
			until := today
			switch true {
			case showToday:
				// show events for today only
				// time limits are already set
			case showTomorrow:
				// show events for today & tomorrow
				from = today
				until = today.AddDate(0, 0, 1)
			case showCurrentWeek:
				// show events for the current week
				from = timeframe.FirstDayOfWeek(today)
				until = timeframe.LastDayOfWeek(today)
			case showNextWeek:
				// show events for the next week
				today = today.AddDate(0, 0, 7)
				from = timeframe.FirstDayOfWeek(today)
				until = timeframe.LastDayOfWeek(today)
			default:
				// show the whole current week per default
				from = timeframe.FirstDayOfWeek(today)
				until = timeframe.LastDayOfWeek(today)
			} // END of switch true

			t := timeframe.NewTimeFrame(from, until)

			files, err := os.ReadDir(dataDir)
			if err != nil {
				prnt.Printf("Error while reading directory %s: %s\n", dataDir, err.Error())
				os.Exit(1)
			}
			for _, f := range files {
				// Entry is a directory or not a REM file
				if f.IsDir() || filepath.Ext(f.Name()) != ".rem" {
					continue
				}
				if dbgMode {
					// Output in debug mode only
					prnt.Printf("Reading from file: %s\n", f.Name())
				}
				fHandle, err := os.Open(filepath.Join(dataDir, f.Name()))
				if err != nil {
					prnt.Printf("Error while open rem file %s: %s\n", f.Name(), err.Error())
					os.Exit(1)
				}
				defer fHandle.Close()

				l := lexer.New(fHandle)
				p := parser.New(l)
				reminders := p.ParseReminders()

				for _, s := range reminders.Statements {
					env := object.NewEnvironment()
					evaluator.Eval(s, env)
					t.FillFrame(env)
				}
			} // END for range files
			days := make([]string, 0, len(t.Days))
			for key := range t.Days {
				days = append(days, key)
			}
			sort.Strings(days)

			header := prnt.Sprintf("\nEvents from %s - %s", from.Format("02.01.2006"), until.Format("02.01.2006"))
			fmt.Println(headerStyle.Render(header))
			for _, day := range days {
				events := make([]string, 0)
				if timeframe.SameDay(time.Now(), t.Days[day].Date) {
					dayLine := fmt.Sprintf("\n%s", t.Days[day].Date.Format("02.01.2006"))
					fmt.Println(todayStyle.Render(dayLine))
				} else {
					dayLine := fmt.Sprintf("\n%s", t.Days[day].Date.Format("02.01.2006"))
					fmt.Println(dayStyle.Render(dayLine))
				}
				for _, e := range t.Days[day].Events {
					time := ""
					if e.From != "" {
						time = e.From
					}
					if time != "" && e.Until != "" {
						time += " - " + e.Until
					}
					if time != "" {
						time = fmt.Sprintf("%-13s ... ", time)
						eventLine := fmt.Sprintf("   %s%s", time, e.Description)
						events = append(events, eventStyle.Render(eventLine))
					} else {
						eventLine := allDayStyle.Render(time + e.Description)
						events = append(events, fmt.Sprintf("   %s%s", checkMark, eventLine))
					}

				} // END for range Days.Events
				sort.Strings(events)
				for _, e := range events {
					fmt.Println(e)
				}
			} // END for range days
		}, // END of Run
	} // END of ShowCmd
)

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.PersistentFlags().BoolVarP(&showToday, "today", "t", false, "Shows events for today")
	showCmd.PersistentFlags().BoolVarP(&showTomorrow, "tomorrow", "m", false, "Shows events for today & tomorrow")
	showCmd.PersistentFlags().BoolVarP(&showCurrentWeek, "currentweek", "c", false, "Shows events for the current week")
	showCmd.PersistentFlags().BoolVarP(&showNextWeek, "nextweek", "n", false, "Shows events for the upcoming week")
}
