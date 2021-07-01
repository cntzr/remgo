package cmd

import "golang.org/x/text/feature/plural"

type entry struct {
	tag, key string
	msg      interface{}
}

var entries = [...]entry{
	{"en", "Error while reading directory %s: %s\n", "Error while reading directory %s: %s\n"},
	{"de", "Error while reading directory %s: %s\n", "Fehler beim Lesen des Verzeichnisses %s: %s\n"},
	{"en", "Error while open rem file %s: %s\n", "Error while open rem file %s: %s\n"},
	{"de", "Error while open rem file %s: %s\n", "Fehler beim Öffnen der Datei %s: %s\n"},
	{"en", "Reading from file: %s\n", "Reading from file: %s\n"},
	{"de", "Reading from file: %s\n", "Lesen aus Datei: %s\n"},
	{"en", "\nEvents from %s - %s", "\nEvents from %s - %s"},
	{"de", "\nEvents from %s - %s", "\nTermine vom %s - %s"},
	{"en", "Using config file: %s\n", "Using config file: %s\n"},
	{"de", "Using config file: %s\n", "Konfiguration eingelesen: %s\n"},
	{"en", "Reading from data directory: %s\n", "Reading from data directory: %s\n"},
	{"de", "Reading from data directory: %s\n", "Verwende Datenverzeichnis: %s\n"},
	{"en", "Language for messages: %s\n", "Language for messages: %s\n"},
	{"de", "Language for messages: %s\n", "Sprache für Nachrichten: %s\n"},
	// samples for singular and plural output
	{"en", "%d task(s) remaining!", plural.Selectf(1, "%d",
		"=1", "One task remaining!",
		"=2", "Two tasks remaining!",
		"other", "[1]d tasks remaining!",
	)},
	{"el", "%d task(s) remaining!", plural.Selectf(1, "%d",
		"=1", "Μία εργασία έμεινε!",
		"=2", "Μια-δυο εργασίες έμειναν!",
		"other", "[1]d εργασίες έμειναν!",
	)},
}
