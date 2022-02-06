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
	plainCheckMark = lipgloss.NewStyle().SetString("✓").
			PaddingRight(1).
			String()

	// plainCmd represents the plain command
	plainCmd = &cobra.Command{
		Use:   "plain",
		Short: "Shows dates for today only with linebreaks",
		Long: `Shows dates for today only with linebreaks. Used for piping and
postprocessing, e.g. via mail. No colors, escape sequences or other stuff.`,
		Run: func(cmd *cobra.Command, args []string) {
			// show only today & tomorrow
			today := time.Now()

			t := timeframe.NewTimeFrame(today, today)

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

			header := prnt.Sprintf("Events for %s", today.Format("02.01.2006"))
			mailHeader := prnt.Sprintf("Subject: %s\n", header)
			mailHeader += "MIME-Version: 1.0\n"
			mailHeader += "Content-Type: text/plain; charset=utf-8\n"
			fmt.Println(mailHeader)
			fmt.Println("\n" + header)
			for _, day := range days {
				events := make([]string, 0)
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
						events = append(events, eventLine)
					} else {
						eventLine := time + e.Description
						events = append(events, fmt.Sprintf("   %s%s", plainCheckMark, eventLine))
					}

				} // END for range Days.Events
				sort.Strings(events)
				for _, e := range events {
					fmt.Println(e)
				}
			} // END for range days
		}, // END of Run
	} // END of plainCmd
)

func init() {
	rootCmd.AddCommand(plainCmd)
}
