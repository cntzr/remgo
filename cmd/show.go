package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/cntzr/remgo/evaluator"
	"github.com/cntzr/remgo/lexer"
	"github.com/cntzr/remgo/object"
	"github.com/cntzr/remgo/parser"
	"github.com/cntzr/remgo/timeframe"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	showToday    = false
	showNextWeek = false

	// showCmd represents the show command
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows dates, default period is current week",
		Long: `Shows dates, default period is current week. Support for a 
specific day, next week or specific periods is coming soon.
Output is colored by default.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Color styles for pterm
			headerStyle := pterm.NewStyle(pterm.FgGreen, pterm.Italic)
			dayStyle := pterm.NewStyle(pterm.FgLightYellow, pterm.Italic)
			todayStyle := pterm.NewStyle(pterm.FgLightRed, pterm.Italic)
			eventStyle := pterm.NewStyle(pterm.FgGray)

			date := time.Now()
			// show the whole current week per default
			from := timeframe.FirstDayOfWeek(date)
			until := timeframe.LastDayOfWeek(date)
			if showToday {
				// show only today & tomorrow
				from = time.Now()
				until = from.AddDate(0, 0, 1)
			} else if showNextWeek {
				// show whole next week
				date = time.Now().AddDate(0, 0, 7)
				from = timeframe.FirstDayOfWeek(date)
				until = timeframe.LastDayOfWeek(date)
			}

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
			headerStyle.Println(header)
			for _, day := range days {
				events := make([]string, 0)
				if timeframe.SameDay(time.Now(), t.Days[day].Date) {
					todayStyle.Printf("\n%s\n", t.Days[day].Date.Format("02.01.2006"))
				} else {
					dayStyle.Printf("\n%s\n", t.Days[day].Date.Format("02.01.2006"))
				}
				for _, e := range t.Days[day].Events {
					time := ""
					if e.From != "" {
						time = e.From
					}
					if time != "" && e.Until != "" {
						time = time + " - " + e.Until
					}
					if time != "" {
						time = fmt.Sprintf("%-13s ... ", time)
					}
					events = append(events, fmt.Sprintf("   %s%s", time, e.Description))
				} // END for range Days.Events
				sort.Strings(events)
				for _, e := range events {
					eventStyle.Println(e)
				}
			} // END for range days
		}, // END of Run
	} // END of ShowCmd
)

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.PersistentFlags().BoolVarP(&showToday, "today", "t", false, "Shows events for today & tomorrow only")
	showCmd.PersistentFlags().BoolVarP(&showNextWeek, "nextweek", "n", false, "Shows events for the upcoming week")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
